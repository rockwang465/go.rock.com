package v1

import (
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
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
