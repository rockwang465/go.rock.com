package v1

import (
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/clients/k8s"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"time"
)

type InstanceBriefResp struct {
	Id           int64            `json:"id" example:"1"`
	ClusterName  string           `json:"cluster_name" example:"test-cluster"`
	EnvNamespace string           `json:"env_namespace" example:"default"`
	ProjectName  string           `json:"project_name" example:"test-project"`
	Name         string           `json:"name" example:"senseguard-oauth2-default"`
	ChartName    string           `json:"chart_name" example:"senseguard-oauth2"` // helm deploy in cluster name, example: senseguard-oauth2
	ChartVersion string           `json:"chart_version" example:"1.0.0-dev-fe380d"`
	CreatedAt    models.LocalTime `json:"created_at" example:"2021-03-11 16:47:37"`
	UpdateAt     models.LocalTime `json:"updated_at" example:"2021-03-11 16:47:37"`
	Version      int64            `json:"version" example:"1"`
}

type PaginationInstanceResp struct {
	PageNum  int64                `json:"page_num" binding:"required" example:"1"`
	PageSize int64                `json:"page_size" binding:"required" example:"10"`
	Total    int64                `json:"total" binding:"required" example:"100"`
	Pages    int64                `json:"pages" binding:"required" example:"1"`
	Items    []*InstanceBriefResp `json:"items" binding:"required"`
}

type InstanceDetailResp struct {
	InstanceBriefResp
	LastDeployment int64 `json:"last_deployment" example:"1"` // deployment_id
	AppId          int64 `json:"app_id" example:"1"`
	EnvId          int64 `json:"env_id" example:"1"`
}

type InstancePaginationReq struct {
	GetPaginationReq
	Cluster string `json:"cluster" form:"cluster" binding:"omitempty" example:"test-cluster"` // cluster name
	Project string `json:"project" form:"project" binding:"omitempty" example:"test-project"` // project name
	AppId   int64  `json:"app_id"  form:"app_id"  binding:"omitempty,min=1" example:"1"`
}

type InstanceConfResp struct {
	Name    string `json:"name" example:"bootstrap.yaml"`
	Content string `json:"content" example:"server:\n  port: 8080\n  servlet:\n    context-path: /sys\nspring:\n  application:\n    name: aurora-system-service\n  ..."`
}

type InstancePodResp struct {
	Namespace string         `json:"namespace" binding:"required" example:"default"`
	Pods      []*InstancePod `json:"pods" binding:"required"`
}

type InstancePod struct {
	Name          string           `json:"name" binding:"required" example:"senseguard-oauth2-7b78686878-vcx79"`
	Containers    []*ContainerResp `json:"containers" binding:"required"`
	RestartPolicy string           `json:"restart_policy" binding:"required" example:"Always"`
	DnsPolicy     string           `json:"dns_policy" binding:"required" example:"ClusterFirst"`
	NodeName      string           `json:"node_name" binding:"required" example:"k8s-master1"`
	HostNetwork   bool             `json:"host_network" binding:"required" example:"false"`
	Status        string           `json:"status" binding:"required" example:"Running"`
	PodIp         string           `json:"pod_ip" binding:"required" example:"10.244.0.70"`
	StartTime     time.Time        `json:"start_time" binding:"required" example:"2021-03-11T14:49:55+08:00"`
}

type ContainerResp struct {
	Name  string `json:"name" binding:"required" example:"senseguard-oauth2"`
	Image string `json:"image" binding:"required" example:"10.151.3.75/sensenebula-guard-std/senseguard-oauth2:1.0.0-dev-fe380d"`
}

type InstanceLogReq struct {
	Pod       string `json:"pod" form:"pod" binding:"required" example:"senseguard-oauth2-7b78686878-vcx79"` // pod name
	Container string `json:"container" form:"container" binding:"required" example:"senseguard-oauth2"`      // container name
}

type InstanceLogResp struct {
	PodName       string `json:"pod_name" example:"senseguard-oauth2-7b78686878-vcx79"`
	ContainerName string `json:"container_name" example:"senseguard-oauth2"`
	Content       string `json:"content" binding:"required" example:"log content here"`
}

// @Summary Get app instance's list by app id
// @Description Api for get app app instance's list by app id
// @Tags APP
// @Accept json
// @Produce json
// @Param id path integer true "App ID"
// @Param page_num query integer true "Request page number" default(1)
// @Param page_size query integer true "Instance number page size " default(10)
// @Success 200 {object} v1.PaginationInstanceResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/apps/{id}/instances [get]
func (c *Controller) GetAppInstances(ctx *gin.Context) {
	// 通过app_id获取该应用的部署实例(应用管理-应用-查看)
	// 查看该应用部署到哪些集群上去了
	var paginationReq GetPaginationReq
	if err := ctx.ShouldBind(&paginationReq); err != nil {
		panic(err)
	}

	var uriReq IdReq // app_id
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	instancePg, err := api.GetAppInstances(uriReq.Id, paginationReq.PageNum, paginationReq.PageSize)
	if err != nil {
		panic(err)
	}

	resp := PaginationInstanceResp{}
	if err := utils.MarshalResponse(instancePg, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get app instances by app_id(%v), this pagination instance number is: %v", uriReq.Id, len(resp.Items))
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get all instances
// @Description Api for get all instances
// @Tags INSTANCE
// @Accept json
// @Produce json
// @Param page_num query integer true "Request page number" default(1)
// @Param page_size query integer true "Instance number page size " default(10)
// @Param query_field query string false "Fuzzy Query(field: name)"
// @Param cluster query string false "Cluster name "
// @Param project query string false "Project name "
// @Param app_id query integer false "App Id"
// @Success 200 {object} v1.PaginationInstanceResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/instances [get]
func (c *Controller) GetInstances(ctx *gin.Context) {
	var req InstancePaginationReq
	if err := ctx.ShouldBind(&req); err != nil {
		panic(err)
	}

	instancesPg, err := api.GetInstances(req.PageNum, req.PageSize, req.QueryField, req.Cluster, req.Project, req.AppId)
	if err != nil {
		panic(err)
	}

	resp := PaginationInstanceResp{}
	if err := utils.MarshalResponse(instancesPg, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get all instances, this pagination instance number is %v", len(resp.Items))
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get instances by id
// @Description Api for get instances by id
// @Tags INSTANCE
// @Accept json
// @Produce json
// @Param id path integer true "Instance ID"
// @Success 200 {object} v1.InstanceDetailResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/instances/{id} [get]
func (c *Controller) GetInstance(ctx *gin.Context) {
	var uriReq IdReq
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	instance, err := api.GetInstanceById(uriReq.Id)
	if err != nil {
		panic(err)
	}

	resp := InstanceDetailResp{}
	if err := utils.MarshalResponse(instance, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get instance by id %v", uriReq.Id)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get instance's config by instance id
// @Description Api for get instance's config by instance id
// @Tags INSTANCE
// @Accept json
// @Produce json
// @Param id path integer true "Instance ID"
// @Success 200 {array} v1.InstanceConfResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/instances/{id}/configs [get]
func (c *Controller) GetInstanceConfig(ctx *gin.Context) {
	// get app config file in config map by instance id
	var uriReq IdReq // instance_id
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	instance, err := api.GetInstanceById(uriReq.Id)
	if err != nil {
		panic(err)
	}

	env, err := api.GetEnvById(instance.EnvId)
	if err != nil {
		panic(err)
	}

	cluster, err := api.GetClusterById(env.ClusterId)
	if err != nil {
		panic(err)
	}

	// get instance configMap
	configMapList, err := k8s.GetInstanceConfig(cluster.Config, env.Namespace, instance.Name)
	if err != nil {
		panic(err)
	}

	// format the configmap to config
	resp := formatInstanceConfigs(configMapList)

	c.Logger.Infof("Get instance's config map with instance id %v", uriReq.Id)
	ctx.JSON(http.StatusOK, resp)
}

// format the configmap to config
func formatInstanceConfigs(configMapList *corev1.ConfigMapList) []*InstanceConfResp {
	iConfResps := []*InstanceConfResp{}
	for _, configMap := range configMapList.Items {
		for name, content := range configMap.Data {
			resp := &InstanceConfResp{
				Name:    name,
				Content: content,
				// [
				//    {
				//        "name": "bootstrap.yml",
				//        "content": "server:\n  port: 8080\n  servlet:\n    context-path: /sys\nspring:\n  application:\n    name: aurora-system-service\n  main:\n    allow-bean-definition-overriding: true\n  messages:\n    basename: i18n/messages\n    encoding: UTF-8\n  cloud:\n    nacos:\n      config:\n        server-addr: ${NACOS_HOST}:${NACOS_PORT}\n        namespace: ${NACOS_NAMESPACE}\n        file-extension: yml\n        shared-configs:\n          - data-id: nacos_discovery.yml\n            group: COMMON_GROUP\n            refresh: true\n          - data-id: feign.yml\n            group: COMMON_GROUP\n            refresh: true\n          - data-id: log.yml\n            group: COMMON_GROUP\n            refresh: true\n          - data-id: datasource.yml\n            group: COMMON_GROUP\n            refresh: true\n          - data-id: mybatis.yml\n            group: COMMON_GROUP\n            refresh: true\n          - data-id: redis.yml\n            group: COMMON_GROUP\n            refresh: true\n          - data-id: aurora-system-service.yml\n            group: DEFAULT_GROUP\n            refresh: true\n\nNACOS_HOST: nacos.component.svc.cluster.local\nNACOS_PORT: 8848\nNACOS_NAMESPACE: prod\n\nfeign:\n  client:\n    config: #设备心跳和删除过期人像定时任务特定超时时间\n      aurora-device-manager-service:\n        connectTimeout: 100000\n        readTimeout: 100000\n        loggerLevel: basic\n      aurora-portrait-manager-service:\n        connectTimeout: 100000\n        readTimeout: 100000\n        loggerLevel: basic\n\nminio:\n  url: http://minio-default.component:9000\n  accesskey: minio\n  secretKey: minio123\n  bucketName: aurora-system-service-export\n"
				//    }
				//]
			}
			iConfResps = append(iConfResps, resp)
		}
	}
	return iConfResps
}

// @Summary Get instance's relevant deployments by instance id
// @Description Api to get instance's relevant deployments by instance id
// @Tags INSTANCE
// @Accept json
// @Produce json
// @Param id path integer true "Instance ID"
// @Param page_num query integer true "Request page number" default(1)
// @Param page_size query integer true "Deployment number page size " default(10)
// @Success 200 {array} v1.PaginationDeploymentResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/instances/{id}/deployments [get]
func (c *Controller) GetInstanceDeployment(ctx *gin.Context) {
	var paginationReq GetPaginationReq
	if err := ctx.ShouldBind(&paginationReq); err != nil {
		panic(err)
	}

	var uriReq IdReq // instance_id
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	instance, err := api.GetInstanceById(uriReq.Id)
	if err != nil {
		panic(err)
	}

	deploymentPg, err := api.GetDeploymentsByName(instance.Name, paginationReq.PageNum, paginationReq.PageSize)
	if err != nil {
		panic(err)
	}

	resp := PaginationDeploymentResp{}
	if err := utils.MarshalResponse(deploymentPg, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get name %v deployments by instance id %v", instance.Name, uriReq.Id)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get instance's log by instance id
// @Description Api to get instance's log by instance id
// @Tags INSTANCE
// @Accept json
// @Produce json
// @Param id path integer true "Instance ID"
// @Param pod query string true "Instance's pod name"
// @Param container query string true "Pod's container name"
// @Success 200 {array} v1.InstanceLogResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/instances/{id}/logs [get]
func (c *Controller) GetInstanceLog(ctx *gin.Context) {
	// 需要先调 GetInstancePods api，获取pod container名称
	// 然后将pod container名称传入当前 GetInstanceLog api中
	// get pod log by pod name and container name
	var uriReq IdReq // instance_id
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	var logReq InstanceLogReq
	if err := ctx.ShouldBind(&logReq); err != nil {
		panic(err)
	}

	instance, err := api.GetInstanceById(uriReq.Id)
	if err != nil {
		panic(err)
	}

	env, err := api.GetEnvById(instance.EnvId)
	if err != nil {
		panic(err)
	}
	cluster, err := api.GetClusterById(env.ClusterId)
	if err != nil {
		panic(err)
	}

	// get the instance's pod log
	podLog, err := k8s.GetInstanceLog(cluster.Config, instance.EnvNamespace, logReq.Pod, logReq.Container, false)
	if err != nil {
		panic(err)
	}

	resp := InstanceLogResp{
		PodName:       logReq.Pod,
		ContainerName: logReq.Container,
		Content:       podLog,
	}

	c.Logger.Infof("Get specific pod name %v, container name %v instance log by instance id %v", logReq.Pod, logReq.Container, uriReq.Id)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get instance's pods name by instance id
// @Description Api to get instance's pods name and containers name by instance id
// @Tags INSTANCE
// @Accept json
// @Produce json
// @Param id path integer true "Instance ID"
// @Success 200 {array} v1.InstancePodResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/instances/{id}/pods [get]
func (c *Controller) GetInstancePods(ctx *gin.Context) {
	// get pod name and container name by instance id
	var uriReq IdReq // instance_id
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	instance, err := api.GetInstanceById(uriReq.Id)
	if err != nil {
		panic(err)
	}

	env, err := api.GetEnvById(instance.EnvId)
	if err != nil {
		panic(err)
	}

	cluster, err := api.GetClusterById(env.ClusterId)
	if err != nil {
		panic(err)
	}

	podList, err := k8s.GetInstancePods(cluster.Config, env.Namespace, instance.Name)
	if err != nil {
		panic(err)
	}
	resp := formatInstancePods(podList)

	c.Logger.Infof("Get instance pods name by instance id %v", uriReq.Id)
	ctx.JSON(http.StatusOK, resp)
	// resp:
	// [
	//    {
	//        "namespace": "default",
	//        "pods": [
	//            {
	//                "name": "senseguard-device-management-75f7565f58-m6wp4",
	//                "containers": [
	//                    {
	//                        "name": "senseguard-device-management",
	//                        "image": "10.151.3.75/sensenebula-guard-std/senseguard-device-management:2.3.0-2.3.0-dev-643384"
	//                    }
	//                ],
	//                "restart_policy": "Always",
	//                "dns_policy": "ClusterFirst",
	//                "node_name": "k8s-master1",
	//                "host_network": false,
	//                "status": "Running",
	//                "pod_ip": "10.244.0.70",
	//                "start_time": "2021-03-11T14:49:55+08:00"
	//            }
	//        ]
	//    }
	//]
}

// format the podList to InstancePodResp
func formatInstancePods(podList *corev1.PodList) []*InstancePodResp {
	nsMapper := make(map[string][]*InstancePod)
	for _, pod := range podList.Items {
		ns := pod.Namespace
		status := getPodsContainersStatus(pod.Status)
		containers := getPodsContainers(pod.Spec.Containers)
		iPod := &InstancePod{
			Name:          pod.Name,
			RestartPolicy: string(pod.Spec.RestartPolicy),
			DnsPolicy:     string(pod.Spec.DNSPolicy),
			NodeName:      pod.Spec.NodeName,
			HostNetwork:   pod.Spec.HostNetwork,
			Status:        status,
			PodIp:         pod.Status.HostIP,
			StartTime:     pod.Status.StartTime.Time,
			Containers:    containers,
		}
		nsMapper[ns] = append(nsMapper[ns], iPod)
	}

	iPods := []*InstancePodResp{}
	for ns, pods := range nsMapper {
		iPods = append(iPods, &InstancePodResp{
			Namespace: ns,
			Pods:      pods,
		})
	}
	return iPods
}

// get containers info
func getPodsContainers(containers []corev1.Container) []*ContainerResp {
	containerResp := []*ContainerResp{}
	for _, c := range containers {
		container := &ContainerResp{
			Name:  c.Name,
			Image: c.Image,
		}
		containerResp = append(containerResp, container)
	}
	return containerResp
}

// get pod status
func getPodsContainersStatus(status corev1.PodStatus) string {
	for _, c := range status.ContainerStatuses {
		if c.Ready == false {
			if c.State.Waiting != nil {
				return c.State.Waiting.Reason
			}
		}
	}
	return "Running"
}
