package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/clients/k8s"
	"go.rock.com/rock-platform/rock/server/clients/license"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/utils"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	LicenseMode = "voucher" // dongle or voucher
)

// license-ca status struct from cactl.go
type CAStatusResp struct {
	Server      uint                   `json:"server" example:"0"`                                     // 0 is master, 1 is slave
	Mode        string                 `json:"mode" example:"voucher"`                                 // license-ca authorization mode, default is voucher mode
	Disable     bool                   `json:"disable" example:"false"`                                // is ca disabled, by soft start/stop ca, if disabled, ca can't supply nomal service
	IsActive    bool                   `json:"is_active" example:"true"`                               // master or standby ca
	ActiveLimit int32                  `json:"active_limit" example:"100"`                             // cluster total active limit
	AloneTime   int32                  `json:"alone_time" example:"0"`                                 // ca alone time, uint seconds, 0 means forever
	DongleTime  int64                  `json:"dongle_time" example:"1616762924"`                       // dongle timestamp
	Status      string                 `json:"status" example:"alone"`                                 // ca status, "alone" or "alive" or "dead", means whether ca is in alive
	AuthID      string                 `json:"auth_id" example:"495788f9-9797-4bf8-a3e1-d65d09b107cd"` // cluster license sn
	Product     string                 `json:"product" example:"IVA-VIPER"`                            // product name
	DongleID    string                 `json:"dongle_id" example:"494330853"`                          // dongle id
	ExpiredAt   string                 `json:"expired_at" example:"99991231"`                          // expire time
	Company     string                 `json:"company" example:"sensetime_SC"`                         // company name
	FeatureIds  []uint64               `json:"feature_ids" example:"22000"`                            // feature ids
	Quotas      map[string]quotaLimit  `json:"quotas"`                                                 // cluster quotas, used and total
	Consts      map[string]interface{} `json:"consts"`                                                 // cluster consts, value type will be int32 or string
	Devices     []caDeviceInfo         `json:"devices"`                                                // the quotas that devices have taken
}

type quotaLimit struct {
	Used  int32 `json:"used" example:"1"`  // used quotas
	Total int32 `json:"total" example:"2"` // total quotas
}

type caDeviceInfo struct {
	UdID       string           `json:"udid,omitempty" example:"engine-face-extract-service-kd4k9-a954a1f74cd23d97248249d04de10221-fba9aae9f524e083"`
	QuotaUsage map[string]int32 `json:"quota_usage,omitempty"`
}

type K8sMasterInfo struct {
	Master1IP   string `json:"master1_ip" example:"10.151.5.136"`
	Master2IP   string `json:"master2_ip" example:"10.151.5.137"`
	MasterTotal uint   `json:"master_total" example:"3"`
}

type LicenseModeReq struct {
	LicenseMode string `json:"license_mode" form:"license_mode" binding:"omitempty,min=1" example:"voucher or dongle"` // default voucher
}

type LicenseServerTypeReq struct {
	// 由于 serverType 只有0 和1 两个值。但定义required，则0位uint的零值，gin validate以为你没有输入。
	// 所以要定义最大为1，最小为0，不要加required字段。这里可以用 min=0,max=1 或者 oneof=0 1两种写法。
	ServerType uint `json:"server_type" form:"server_type" binding:"required,oneof=0 1" example:"0"` // 0 is master, 1 is slave
	LicenseModeReq
}

type ActiveResp struct {
	StatusCode    string `json:"status_code,omitempty" example:"000000"`
	StatusMessage string `json:"status_message,omitempty" example:"SUCCESS"`
}

// Client Licenses Response
type ClientLicResp struct {
	Licenses []string `json:"licenses,omitempty"`
}

// @Summary Get license status by cluster id and license mode
// @Description api for get license status by cluster id and license mode
// @Tags CLUSTER
// @Accept json
// @Produce json
// @Param id path integer true "Cluster ID"
// @Param license_mode query string false "license mode"
// @Success 200 {array} v1.CAStatusResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/clusters/{id}/license-status [get]
func (c *Controller) GetLicenseStatus(ctx *gin.Context) {
	var idReq IdReq // cluster id
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	var req LicenseModeReq
	if err := ctx.ShouldBind(&req); err != nil {
		panic(err)
	}

	// ensure license authorization mode
	mode := getLicenseMode(req.LicenseMode)

	cluster, err := api.GetClusterById(idReq.Id)
	if err != nil {
		panic(err)
	}

	// get the k8s master nodes info
	k8sClusterInfo, err := getClusterIp(cluster.Config)
	if err != nil {
		panic(err)
	}

	// get license-ca status
	var caMasterStatus *license.CAStatus
	var caSlaveStatus *license.CAStatus
	var caStatus []*license.CAStatus
	if k8sClusterInfo.MasterTotal >= 3 { // when cluster mode
		masterCAUrl := utils.GetLicenseCaUrl(k8sClusterInfo.Master1IP)
		slaveCAUrl := utils.GetLicenseCaUrl(k8sClusterInfo.Master2IP)
		caCtl, err := license.NewServiceCtl(masterCAUrl, slaveCAUrl) // get license-ca client
		if err != nil {
			panic(err)
		}
		caMasterStatus, err = caCtl.GetCAStatus(0, mode) // get license-ca master status
		if err != nil {
			panic(err)
		}

		caSlaveStatus, err = caCtl.GetCAStatus(1, mode) // get license-ca slave status
		if err != nil {
			panic(err)
		}

		caStatus = []*license.CAStatus{caMasterStatus, caSlaveStatus}
	} else { // when single node mode
		masterCAUrl := utils.GetLicenseCaUrl(k8sClusterInfo.Master1IP)
		caCtl, err := license.NewServiceCtl(masterCAUrl, masterCAUrl) // get license-ca client
		if err != nil {
			panic(err)
		}
		caMasterStatus, err = caCtl.GetCAStatus(0, mode)
		if err != nil {
			panic(err)
		}

		caStatus = []*license.CAStatus{caMasterStatus}
	}

	resp := make([]CAStatusResp, 2)
	err = utils.MarshalResponse(caStatus, &resp)
	if err != nil {
		panic(err)
	}

	c.Logger.Infof("Get license ca status by cluster id %v, authorization mode is %v", cluster.Id, mode)
	ctx.JSON(http.StatusOK, resp)
}

func getLicenseMode(licenseMode string) string {
	if licenseMode != "" {
		return licenseMode
	} else {
		return LicenseMode
	}
}

func getClusterIp(k8sConf string) (*K8sMasterInfo, error) {
	// check the total number of K8S cluster master
	nodeList, err := k8s.GetClusterNodes(k8sConf)
	if err != nil {
		return nil, err
	}
	nodes, err := formatNodesResp(nodeList.Items)
	if err != nil {
		return nil, err
	}

	// gets the total number of k8s cluster master nodes
	masterTotal, err := getClusterMaster(nodes)
	if err != nil {
		return nil, err
	}

	// get the master node's IP by k8s cluster
	masterClusterInfo := new(K8sMasterInfo)
	if masterTotal >= 3 { // when cluster mode
		masterClusterInfo.Master1IP = (*nodes)[0].InternalIP
		masterClusterInfo.Master2IP = (*nodes)[1].InternalIP
		masterClusterInfo.MasterTotal = masterTotal
		return masterClusterInfo, nil
	} else if masterTotal == 1 || masterTotal == 2 { // when single node mode
		masterClusterInfo.Master1IP = (*nodes)[0].InternalIP
		masterClusterInfo.MasterTotal = masterTotal
		return masterClusterInfo, nil
	} else {
		err := utils.NewRockError(404, 40400015, fmt.Sprintf("k8s cluster node not found"))
		return nil, err
	}
}

// Gets the total number of K8S cluster master nodes
func getClusterMaster(nodes *[]ClusterNodeResp) (uint, error) {
	var masterTotal uint = 0
	for _, node := range *nodes {
		for _, label := range node.Labels {
			if label.Key == "node-role.kubernetes.io/master" && label.Value == "" {
				masterTotal += 1
			}
		}
	}
	return masterTotal, nil
}

// @Summary Download license hardware c2v file
// @Description api for download license hardware c2v file
// @Tags CLUSTER
// @Accept json
// @Produce json
// @Param id path integer true "Cluster ID"
// @Param server_type query uint true "license master or slave type"
// @Param license_mode query string false "license mode"
// @Success 200 {string} string "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/clusters/{id}/license-c2v [get]
func (c *Controller) GetC2vFile(ctx *gin.Context) {
	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	var serverTypeReq LicenseServerTypeReq
	if err := ctx.ShouldBind(&serverTypeReq); err != nil {
		panic(err)
	}

	// ensure license authorization mode
	mode := getLicenseMode(serverTypeReq.LicenseMode)

	cluster, err := api.GetClusterById(idReq.Id)
	if err != nil {
		panic(err)
	}

	// get the k8s master nodes info
	k8sClusterInfo, err := getClusterIp(cluster.Config)
	if err != nil {
		panic(err)
	}

	var caCtl *license.CACtl
	if k8sClusterInfo.MasterTotal >= 3 { // when cluster mode
		masterCAUrl := utils.GetLicenseCaUrl(k8sClusterInfo.Master1IP)
		slaveCAUrl := utils.GetLicenseCaUrl(k8sClusterInfo.Master2IP)
		caCtl, err = license.NewServiceCtl(masterCAUrl, slaveCAUrl) // get license-ca client
		if err != nil {
			panic(err)
		}
	} else {
		masterCAUrl := utils.GetLicenseCaUrl(k8sClusterInfo.Master1IP)
		caCtl, err = license.NewServiceCtl(masterCAUrl, masterCAUrl) // get license-ca client
		if err != nil {
			panic(err)
		}
	}

	// ServerType: 0 is master, 1 is slave
	// Type: 0 is c2v + fingerprint
	hardwareInfoResp, err := caCtl.HardwareInfo(license.ServerType(serverTypeReq.ServerType), 0)
	if err != nil {
		panic(err)
	}

	// save c2v file
	c2vTmpFile, err := ioutil.TempFile("/tmp", "c2v-file-*")
	if err != nil {
		panic(err)
	}
	defer os.Remove(c2vTmpFile.Name())

	_, err = c2vTmpFile.WriteString(hardwareInfoResp.C2V)
	if err != nil {
		panic(err)
	}
	defer c2vTmpFile.Close()

	fileName := "default.c2v"
	status, err := caCtl.GetCAStatus(license.ServerType(serverTypeReq.ServerType), mode) // get license-ca master/slave status
	if err != nil {
		c.Logger.Warnf("Get dongle id failed, set dongle id to default and skip it: %v", err)
	} else {
		fileName = fmt.Sprintf("%s.c2v", status.DongleID)
	}

	c.Logger.Infof("Download c2v hardware info by cluster id %d, authorization mode is %v", idReq.Id, mode)

	// 为了前端通过调用当前接口就能直接下载文件，这里必须配置如下格式(filename + application/octet-stream):
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.File(c2vTmpFile.Name()) // 读取文件内容并返回
}

// @Summary Download license hardware fingerprint file
// @Description api for download license hardware fingerprint file
// @Tags CLUSTER
// @Accept json
// @Produce json
// @Param id path integer true "Cluster ID"
// @Param server_type query uint true "license master or slave type"
// @Param license_mode query string false "license mode"
// @Success 200 {string} string "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/clusters/{id}/license-fingerprint [get]
func (c *Controller) GetFingerprintFile(ctx *gin.Context) {
	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	var serverTypeReq LicenseServerTypeReq
	if err := ctx.ShouldBind(&serverTypeReq); err != nil {
		panic(err)
	}

	// ensure license authorization mode
	mode := getLicenseMode(serverTypeReq.LicenseMode)

	cluster, err := api.GetClusterById(idReq.Id)
	if err != nil {
		panic(err)
	}

	// get the k8s master nodes info
	k8sClusterInfo, err := getClusterIp(cluster.Config)
	if err != nil {
		panic(err)
	}

	var caCtl *license.CACtl
	if k8sClusterInfo.MasterTotal >= 3 { // when cluster mode
		masterCAUrl := utils.GetLicenseCaUrl(k8sClusterInfo.Master1IP)
		slaveCAUrl := utils.GetLicenseCaUrl(k8sClusterInfo.Master2IP)
		caCtl, err = license.NewServiceCtl(masterCAUrl, slaveCAUrl) // get license-ca client
		if err != nil {
			panic(err)
		}
	} else {
		masterCAUrl := utils.GetLicenseCaUrl(k8sClusterInfo.Master1IP)
		caCtl, err = license.NewServiceCtl(masterCAUrl, masterCAUrl) // get license-ca client
		if err != nil {
			panic(err)
		}
	}

	// ServerType: 0 is master, 1 is slave
	// Type: 0 is c2v + fingerprint
	hardwareInfoResp, err := caCtl.HardwareInfo(license.ServerType(serverTypeReq.ServerType), 0)
	if err != nil {
		panic(err)
	}

	// save fingerprint file
	fpTmpFile, err := ioutil.TempFile("/tmp", "fingerprint-file-*")
	if err != nil {
		panic(err)
	}
	defer os.Remove(fpTmpFile.Name())

	_, err = fpTmpFile.WriteString(hardwareInfoResp.FingerPrint)
	if err != nil {
		panic(err)
	}
	defer fpTmpFile.Close()

	fileName := "default.fingerprint"
	status, err := caCtl.GetCAStatus(license.ServerType(serverTypeReq.ServerType), mode) // get license-ca master/slave status
	if err != nil {
		c.Logger.Warnf("Get dongle id failed, set dongle id to default and skip it: %v", err)
	} else {
		fileName = fmt.Sprintf("%s.fingerprint", status.DongleID)
	}

	c.Logger.Infof("Download fingerprint hardware info by cluster id %d, authorization mode is %v", idReq.Id, mode)

	// 为了前端通过调用当前接口就能直接下载文件，这里必须配置如下格式(filename + application/octet-stream):
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.File(fpTmpFile.Name()) // 读取文件内容并返回
}

// @Summary Online active the license-ca
// @Description api for online active the license-ca
// @Tags CLUSTER
// @Accept json
// @Produce json
// @Param id path integer true "Cluster ID"
// @Param server_type query uint true "license master or slave type"
// @Param license_mode query string false "license mode"
// @Success 200 {object} v1.ActiveResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/clusters/{id}/license-online [post]
func (c *Controller) ActiveOnline(ctx *gin.Context) {
	// 如果在线激活报错，请去官网查看报错原因。https://dongle.sensetime.com/errors
	// 因为license-ca源码是闭源的，只能基础研发才有权限查看。
	// 示例报错码:
	//    "status_code": "000002",
	//    "status_message": "get online data from sdk error of -2014281727"
	// 查询结果:
	// ID	错误码		厂商错误码	厂商错误码描述
	// 154	-2014281727	0x8001		链接激活服务器失败

	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	var serverTypeReq LicenseServerTypeReq
	if err := ctx.ShouldBind(&serverTypeReq); err != nil {
		panic(err)
	}

	// ensure license authorization mode
	mode := getLicenseMode(serverTypeReq.LicenseMode)

	cluster, err := api.GetClusterById(idReq.Id)
	if err != nil {
		panic(err)
	}

	// get the k8s master nodes info
	k8sClusterInfo, err := getClusterIp(cluster.Config)
	if err != nil {
		panic(err)
	}

	var caCtl *license.CACtl
	if k8sClusterInfo.MasterTotal >= 3 { // when cluster mode
		masterCAUrl := utils.GetLicenseCaUrl(k8sClusterInfo.Master1IP)
		slaveCAUrl := utils.GetLicenseCaUrl(k8sClusterInfo.Master2IP)
		caCtl, err = license.NewServiceCtl(masterCAUrl, slaveCAUrl) // get license-ca client
		if err != nil {
			panic(err)
		}
	} else {
		masterCAUrl := utils.GetLicenseCaUrl(k8sClusterInfo.Master1IP)
		caCtl, err = license.NewServiceCtl(masterCAUrl, masterCAUrl) // get license-ca client
		if err != nil {
			panic(err)
		}
	}

	// online active
	activeResp, err := caCtl.OnlineActivate(license.ServerType(serverTypeReq.ServerType), "activate")
	if err != nil {
		panic(err)
	}

	if activeResp.StatusCode != "000000" {
		message := activeResp.StatusMessage + ", You can use the error code to query the reason in here: https://dongle.sensetime.com/errors"
		err := utils.NewRockError(500, 50000007, message)
		panic(err)
	}

	// check license-ca status is active
	status, err := caCtl.GetCAStatus(license.ServerType(serverTypeReq.ServerType), mode)
	if err != nil {
		panic(err)
	}
	if status.IsActive == false {
		err := utils.NewRockError(500, 50000008, "License is not active even after executing license online activation")
		panic(err)
	}

	resp := ActiveResp{}
	if err := utils.MarshalResponse(activeResp, &resp); err != nil {
		panic(err)
	}

	// get client licenses data from license-ca
	clientLicResp, err := caCtl.GetClientLics(license.ServerType(serverTypeReq.ServerType))
	if err != nil || len(clientLicResp.Licenses) == 0 {
		err := fmt.Sprintf("get client license data failed: %v", err)
		panic(err)
	}

	// update license-config configmap
	err = k8s.UpdateConfigmapWithLicense(cluster.Config, mode, clientLicResp.Licenses[0])
	if err != nil {
		err := fmt.Sprintf("update configmap failed: %v", err)
		panic(err)
	}

	// restart engine namespace pods
	err = k8s.RestartPodsWithLicense(cluster.Config, 60)
	if err != nil {
		panic(err)
	}

	c.Logger.Infof("Online active license success, active by cluster_id %d, authorization mode is %v", idReq.Id, mode)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Offline active the license-ca
// @Description api for offline active the license-ca
// @Tags CLUSTER
// @Accept json
// @Produce json
// @Param id path integer true "Cluster ID"
// @Param server_type query uint true "license master or slave type"
// @Param license_mode query string false "license mode"
// @Success 200 {object} v1.ActiveResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/clusters/{id}/license-offline [post]
func (c *Controller) ActiveOffline(ctx *gin.Context) {
	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	var serverTypeReq LicenseServerTypeReq
	if err := ctx.ShouldBind(&serverTypeReq); err != nil {
		panic(err)
	}

	// get the authorized v2c file
	file, err := ctx.FormFile("v2c")
	if err != nil {
		panic(err)
	}

	// create a temporary file
	tmpFile, err := ioutil.TempFile("", "v2c-file-*")
	if err != nil {
		panic(err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// save v2c file to temporary file
	if err := ctx.SaveUploadedFile(file, tmpFile.Name()); err != nil {
		panic(err)
	}

	// read temporary file content
	v2cData, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		panic(err)
	}

	// ensure license authorization mode
	mode := getLicenseMode(serverTypeReq.LicenseMode)

	cluster, err := api.GetClusterById(idReq.Id)
	if err != nil {
		panic(err)
	}

	// get the k8s master nodes info
	k8sClusterInfo, err := getClusterIp(cluster.Config)
	if err != nil {
		panic(err)
	}

	var caCtl *license.CACtl
	if k8sClusterInfo.MasterTotal >= 3 { // when cluster mode
		masterCAUrl := utils.GetLicenseCaUrl(k8sClusterInfo.Master1IP)
		slaveCAUrl := utils.GetLicenseCaUrl(k8sClusterInfo.Master2IP)
		caCtl, err = license.NewServiceCtl(masterCAUrl, slaveCAUrl) // get license-ca client
		if err != nil {
			panic(err)
		}
	} else {
		masterCAUrl := utils.GetLicenseCaUrl(k8sClusterInfo.Master1IP)
		caCtl, err = license.NewServiceCtl(masterCAUrl, masterCAUrl) // get license-ca client
		if err != nil {
			panic(err)
		}
	}

	// offline active
	activeResp, err := caCtl.OfflineActivate(license.ServerType(serverTypeReq.ServerType), string(v2cData))
	if err != nil {
		panic(err)
	}
	if activeResp.StatusCode != "000000" {
		message := activeResp.StatusMessage + ", You can use the error code to query the reason in here: https://dongle.sensetime.com/errors"
		err := utils.NewRockError(500, 50000009, message)
		panic(err)
	}

	// check license-ca status is active
	status, err := caCtl.GetCAStatus(license.ServerType(serverTypeReq.ServerType), mode)
	if err != nil {
		panic(err)
	}
	if status.IsActive == false {
		err := utils.NewRockError(500, 50000010, "License is not active even after executing license offline activation")
		panic(err)
	}

	resp := ActiveResp{}
	if err := utils.MarshalResponse(activeResp, &resp); err != nil {
		panic(err)
	}

	// get client licenses data from license-ca
	clientLicResp, err := caCtl.GetClientLics(license.ServerType(serverTypeReq.ServerType))
	if err != nil || len(clientLicResp.Licenses) == 0 {
		err := fmt.Sprintf("get client license data failed: %v", err)
		panic(err)
	}

	// update license-config configmap
	err = k8s.UpdateConfigmapWithLicense(cluster.Config, mode, clientLicResp.Licenses[0])
	if err != nil {
		err := fmt.Sprintf("update configmap failed: %v", err)
		panic(err)
	}

	// restart engine namespace pods
	err = k8s.RestartPodsWithLicense(cluster.Config, 60)
	if err != nil {
		panic(err)
	}

	c.Logger.Infof("Offline active license success, active by cluster_id %d, authorization mode is %v", idReq.Id, mode)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get client licenses
// @Description api for get client licenses
// @Tags CLUSTER
// @Accept json
// @Produce json
// @Param id path integer true "Cluster ID"
// @Success 200 {object} v1.ClientLicResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/clusters/{id}/license-clics [get]
func (c *Controller) GetClientLicenses(ctx *gin.Context) {
	// 获取 client_license.lic 证书内容
	// 命令获取方式: license_client license | sed '1,6d'
	// 集群中，不管serverType是0还是1，licence-ca的license.lic证书肯定是一样的，因为lic证书是放在(license-config)configmap中的，所以一样。
	var idReq IdReq // cluster id
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	// 集群中，不管serverType是0还是1，licence-ca的license.lic证书肯定是一样的，因为lic证书是放在(license-config)configmap中的，所以一样。
	// 所以不需要单独获取 license serverType 了
	//var serverTypeReq LicenseServerTypeReq
	//if err := ctx.ShouldBind(&serverTypeReq); err != nil {
	//	panic(err)
	//}

	cluster, err := api.GetClusterById(idReq.Id)
	if err != nil {
		panic(err)
	}

	// get the k8s master nodes info
	k8sClusterInfo, err := getClusterIp(cluster.Config)
	if err != nil {
		panic(err)
	}

	var caCtl *license.CACtl
	if k8sClusterInfo.MasterTotal >= 3 { // when cluster mode
		masterCAUrl := utils.GetLicenseCaUrl(k8sClusterInfo.Master1IP)
		slaveCAUrl := utils.GetLicenseCaUrl(k8sClusterInfo.Master2IP)
		caCtl, err = license.NewServiceCtl(masterCAUrl, slaveCAUrl) // get license-ca client
		if err != nil {
			panic(err)
		}
	} else {
		masterCAUrl := utils.GetLicenseCaUrl(k8sClusterInfo.Master1IP)
		caCtl, err = license.NewServiceCtl(masterCAUrl, masterCAUrl) // get license-ca client
		if err != nil {
			panic(err)
		}
	}

	clientLicResp, err := caCtl.GetClientLics(0)
	if err != nil {
		panic(err)
	}

	resp := ClientLicResp{}
	if err := utils.MarshalResponse(clientLicResp, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get client licenses by cluster_id %v", idReq.Id)
	ctx.JSON(http.StatusOK, resp)
}
