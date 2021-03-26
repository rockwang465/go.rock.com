package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/clients/k8s"
	"go.rock.com/rock-platform/rock/server/clients/license"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
)

const (
	LicenseMode = "voucher" // dongle or voucher
)

type LicenseStatusReq struct {
	LicenseMode string `json:"license_mode" form:"license_mode" binding:"omitempty,min=1" example:"voucher or dongle"` // default voucher
}

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

type NodeBriefInfo struct {
	IP string `json:"internal_ip" example:"10.151.5.136"`
	//Name     string `json:"name" example:"k8s-master1"`
	//Hostname string `json:"hostname" example:"k8s-master1"`
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
// @Router /v1/clusters/{id}/envs [get]
func (c *Controller) GetLicenseStatus(ctx *gin.Context) {
	var idReq IdReq // cluster id
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	var req LicenseStatusReq
	if err := ctx.ShouldBind(&req); err != nil {
		panic(err)
	}

	// ensure license authorization mode
	var mode string
	if req.LicenseMode != "" {
		mode = req.LicenseMode
	} else {
		mode = LicenseMode
	}

	cluster, err := api.GetClusterById(idReq.Id)
	if err != nil {
		panic(err)
	}

	// check the total number of K8S cluster master
	nodeList, err := k8s.GetClusterNodes(cluster.Config)
	if err != nil {
		panic(err)
	}
	nodes, err := formatNodesResp(nodeList.Items)
	if err != nil {
		panic(err)
	}

	// gets the total number of k8s cluster master nodes
	masterTotal, err := getClusterMaster(nodes)
	if err != nil {
		panic(err)
	}

	// get the master node's IP by k8s cluster
	var master1 NodeBriefInfo
	var master2 NodeBriefInfo
	if masterTotal >= 3 { // when cluster mode
		master1 = NodeBriefInfo{
			IP: (*nodes)[0].InternalIP,
		}
		master2 = NodeBriefInfo{
			IP: (*nodes)[1].InternalIP,
		}
	} else if masterTotal == 1 || masterTotal == 2 { // when single node mode
		master1 = NodeBriefInfo{
			IP: (*nodes)[0].InternalIP,
		}
	} else {
		err := utils.NewRockError(404, 40400015, fmt.Sprintf("k8s cluster node not found"))
		panic(err)
	}

	// get license-ca status
	var caMasterStatus *license.CAStatus
	var caSlaveStatus *license.CAStatus
	var caStatus []*license.CAStatus
	if masterTotal >= 3 { // when cluster mode
		masterCAUrl := utils.GetLicenseCaUrl(master1.IP)
		slaveCAUrl := utils.GetLicenseCaUrl(master2.IP)
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
		masterCAUrl := utils.GetLicenseCaUrl(master1.IP)
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

	c.Logger.Infof("Get cluster id %v license ca status success", cluster.Id)
	ctx.JSON(http.StatusOK, resp)
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
