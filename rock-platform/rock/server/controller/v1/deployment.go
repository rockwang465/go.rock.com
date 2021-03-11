package v1

import (
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/clients/helm"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
)

type CreateDeploymentReq struct {
	Description  string `json:"description" binding:"omitempty,max=250" example:"description for deployment"`
	ChartName    string `json:"chart_name" binding:"required" example:"senseguard-oauth2-default"`
	ChartVersion string `json:"chart_version" binding:"required" example:"1.0.0-dev-fe380d"`
	ProjectEnvId int64  `json:"project_env_id" binding:"required,min=1" example:"1"` // project id & namespace id
	AppId        int64  `json:"app_id" binding:"required,min=1" example:"1"`
}

type DeploymentBriefResp struct {
	Id          int64            `json:"id" example:"1"`
	Name        string           `json:"name" example:"senseguard-oauth2-default"`
	Description string           `json:"description" example:"description for deployment"`
	CreatedAt   models.LocalTime `json:"created_at" example:"2021-03-09 15:18:13"`
	UpdatedAt   models.LocalTime `json:"updated_at" example:"2021-03-09 15:18:13"`
	Version     int64            `json:"version" example:"1"`
}

type DeploymentDetailResp struct {
	DeploymentBriefResp
	ChartName    string `json:"chart_name" example:"senseguard-oauth2"`
	ChartVersion string `json:"chart_version" example:"1.0.0-dev-000c37"`
	AppId        int64  `json:"app_id" example:"1"`
	EnvId        int64  `json:"env_id" example:"1"`
}

// @Summary Create deployment
// @Description Api to create deployment
// @Tags DEPLOYMENT
// @Accept json
// @Produce json
// @Param input_body body v1.CreateDeploymentReq true "JSON type input body"
// @Success 201 {object} v1.DeploymentDetailResp "StatusCreated"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/deployments [post]
func (c *Controller) CreateDeployment(ctx *gin.Context) {
	// 1.拉取最新代码编译发版逻辑(CreateAppBuild): -> 被动触发
	//      前端进行单个服务发版，将相关参数传给CreateApp Api，
	//      CreateApp Api 通过 drone-go 模块将相关参数发给 drone-server，
	//      drone-server 将任务下发给 drone-agent，
	//      drone-agent 拉取该应用的源码，根据 .drone.yaml(pipeline) 定义进行任务执行。
	//      当执行 .drone.yaml 最后一步(deploy_to_env)部署应用到指定环境时，会运行infra-drone-plugins中的python脚本，
	//      通过admin jwt token(galaxias_api_token)调用运维平台的 CreateDeployment(当前函数) Api 进行应用部署到指定环境。

	// 2.指定存在的应用发版逻辑(CreateDeployment 当前函数): -> 主动触发
	// 应用管理 -> 应用 -> (翻页/根据project名称&app名称查询)选择某个应用 -> 查看(查看一个应用)
	// -> 应用版本 -> 部署(每个版本的应用后面都有3个选项:查看/部署/删除) ->
	// 右侧弹窗:
	//    版本: 1.0.0-dev-fe380d // 此处无法修改，因为是从这个版本号进来的弹窗
	//    选择项目环境(必选): 10.151.3.99-default(10.151.3.99-default)  // 下拉选择一个环境
	//    简要描述(可选): A Helm chart for Kubernetes  // 简单描述此次部署的原因

	// CreateDeployment api 逻辑描述:
	// a.传入参数: ChartName ChartVersion ProjectEnvId AppId Description
	// b.使用app_id project_env_id拿到从app_conf表中拿到该app的config(values.yaml)
	// c.使用helm模块，检查发版环境tiller服务端是否健康
	// d.使用helm模块，删除发版环境的当前应用版本
	// e.将当前要部署的信息存到deployment表中(appId, envId, chartName, chartVersion, description, namespace)
	// f.使用helm二进制命令，将应用的chart tgz包部署到环境上:
	//    相关参数获取: chartmuseum地址、发布环境地址、tiller地址、app配置转成文件、chart包名、名称空间
	// g.将部署完成后的实例(已经将服务安装到指定环境中，此服务为一个实例)信息保存到instance表中，用于追溯指定环境安装过哪些服务。
	//    相关参数信息包含: clusterName, namespace, projectName, deployment.Name, chartName, chartVersion, deployment.Id, AppId, EnvId
	//    详细信息见: database/models/instance.go

	// 另外，在此处执行之前，该运维平台所在的机器上必须先安装好helm命令，并helm init --client-only --stable-repo-url=http://10.151.3.75:8080初始化成功

	var req CreateDeploymentReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}

	// Get app config by app_id and project_env_id
	appConf, err := api.GetAppConfByAppAndProjectEnvId(req.AppId, req.ProjectEnvId)
	if err != nil {
		panic(err)
	}

	// get env id
	projectEnv, err := api.GetProjectEnvById(req.ProjectEnvId)
	if err != nil {
		panic(err)
	}
	envId := projectEnv.EnvId

	// get cluster id
	env, err := api.GetEnvById(envId)
	if err != nil {
		panic(err)
	}
	clusterId := env.ClusterId

	// check helm tiller is ok by cluster id
	if err := helm.PingTillerServer(clusterId); err != nil {
		panic(err)
	}

	// check service is not failed on remote host, if so, remove it
	releaseName := utils.GenerateChartName(req.ChartVersion, env.Namespace) // kafka-component
	err = helm.DeleteReleaseIfFailedOrDeleted(clusterId, releaseName)
	if err != nil {
		panic(err)
	}

	// Check manual helm release is already installed, if so, remove it
	err = helm.DeleteManualInstallReleaseIfExist(clusterId, releaseName)

	// insert the current deployment info into the database
	deployment, err := api.CreateDeployment(req.AppId, envId, req.ChartName, req.ChartVersion, req.Description, env.Namespace)
	if err != nil {
		panic(err)
	}

	cluster, err := api.GetClusterById(clusterId)
	if err != nil {
		panic(err)
	}

	project, err := api.GetProjectById(projectEnv.ProjectId)
	if err != nil {
		panic(err)
	}

	// Prepare installation parameters
	repoUrl := utils.GetRepoUrl()                                          // get chartmuseum repo address
	chartTgzName := utils.GenChartTgzName(req.ChartName, req.ChartVersion) // get chart tgz package name
	k8sConfig := cluster.Config                                            // get k8s cluster admin.conf
	ns := env.Namespace
	appConfig := appConf.Config

	// deploy the helm chart service by helm client
	if err := utils.InstallOrUpgradeChart(repoUrl, chartTgzName, k8sConfig, ns, deployment.Name, appConfig); err != nil {
		panic(err)
	}

	// create or update instance
	_, err = api.CreateOrUpdateInstance(req.ChartName, req.ChartVersion, cluster.Name, ns, project.Name, deployment)
	if err != nil {
		panic(err)
	}

	resp := DeploymentDetailResp{}
	if err := utils.MarshalResponse(deployment, &resp); err != nil {
		panic(err)
	}
	c.Logger.Infof("Create deployment by id(%v) name(%v)", resp.Id, resp.Name)
	ctx.JSON(http.StatusCreated, resp)
}
