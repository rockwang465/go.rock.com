package v1

import (
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
)

type CreateProjectEnvReq struct {
	Name        string `json:"name" binding:"required" example:"cluster name"` // 要创建的项目环境的名字(project_env表的name)
	EnvId       int64  `json:"env_id" binding:"required,min=1" example:"1"`    // env表中对应(基于cluster_id进行查询的)namespace对应的id字段值。如cluster_id为164的default名称空间对应的id为448
	Description string `json:"description" binding:"omitempty,max=100" example:"description the cluster project env"`
}

type ProjectEnvBriefResp struct {
	Id          int64            `json:"id" example:"1"`
	Name        string           `json:"name" example:"test-project"`
	Description string           `json:"description" example:"description the project"`
	EnvId       int64            `json:"env_id" example:"1"`
	ProjectId   int64            `json:"project_id" example:"1"`
	CreatedAt   models.LocalTime `json:"created_at" example:"2021-01-28 20:20:20"`
	UpdateAt    models.LocalTime `json:"updated_at" example:"2021-01-28 20:20:20"`
}

type PaginationProjectEnvResp struct {
	PageNum  int64                  `json:"page_num" binding:"required" example:"1"`
	PageSize int64                  `json:"page_size" binding:"required" example:"10"`
	Total    int64                  `json:"total" binding:"required" example:"100"`
	Pages    int64                  `json:"pages" binding:"required" example:"1"`
	Items    []*ProjectEnvBriefResp `json:"items" binding:"required"`
}

type ProjectIdAndEnvIdReq struct {
	Id           int64 `json:"id" uri:"id" binding:"required,min=1" example:"1"`       // project id
	ProjectEnvId int64 `json:"pe_id" uri:"pe_id" binding:"required,min=1" example:"1"` // project_env id
}

type UpdateProjectEnvReq struct {
	Description string `json:"description" binding:"omitempty,max=100" example:"description the project env"`
}

// @Summary Create project env by project_id env_id name
// @Description api for create project env by project_id env_id name
// @Tags PROJECT
// @Accept json
// @Produce json
// @Param id path integer true "Project ID"
// @Param input_body body v1.CreateProjectEnvReq true "JSON type input body"
// @Success 200 {object} v1.ProjectEnvBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/projects/{id}/project-envs [post]
func (c *Controller) CreateProjectEnv(ctx *gin.Context) {
	// 功能: 给一个project增加项目环境(在project_env表中新增数据)
	// 项目环境含义: 在一个project(项目)上添加一个集群环境的namespace，并起个名字，添加描述
	// 作用: 允许发布这个project(项目)下所有app(应用)到这个集群的这个namespace(名称空间)
	// 先点击[查看] 查看一个项目,所以这里获取到了project_id。打开新的页面，操作如下:
	// A.点击 [新建项目环境]按钮,触发 http://10.151.3.xx:8888/v1/clusters?page_num=1&page_size=1000&queryfiled=, (ctlv1.GetClusters router)获取cluster集群环境的ip
	// B.右侧弹窗: 点击 [选择集群]下拉按钮,选择A中生成渲染出的IP地址(集群ip)
	// C.右侧弹窗: 点击 [选择集群环境]下拉按钮,拿到B中选择的ip,触发 http://10.151.3.xx:8888/v1/clusters/164/envs?page_num=1&page_page=10, ( ctlv1.GetClusterEnvs router)将所有名称空间渲染出来,让用户选择名称空间。注意，该名称空间是对应env表中的一个id的。(通过cluster_id在env表中的查询,拿到该cluster_id拥有的名称空间，并渲染到此处，且保存env_id用于下面的传参)
	// D.右侧弹窗: [环境名称] 中输入环境的名称
	// E.右侧弹窗: [简要描述] 中输入环境的描述
	// F.右侧弹窗: 点击 [新建] 按钮，进行项目环境创建。传参为: C中的env的id, D中的环境名称, E中的描述
	var idReq IdReq // project id
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	var req CreateProjectEnvReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}

	projectEnv, err := api.CreateProjectEnvByProjectId(idReq.Id, req.EnvId, req.Name, req.Description)
	if err != nil {
		panic(err)
	}

	resp := ProjectEnvBriefResp{}

	if err := utils.MarshalResponse(projectEnv, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Create project_env id:%v by project_id:%v env_id:%v name:%v", resp.Id, idReq.Id, req.EnvId, req.Name)
	ctx.JSON(http.StatusCreated, resp)
}

// @Summary Get all project env
// @Description api for get all project env
// @Tags PROJECT
// @Accept json
// @Produce json
// @Param id path integer true "Project ID"
// @Success 200 {object} v1.ProjectEnvBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/projects/{id}/project-envs [get]
func (c *Controller) GetProjectEnvs(ctx *gin.Context) {
	var idReq IdReq // project id
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	var paginationReq GetPaginationReq
	if err := ctx.ShouldBind(&paginationReq); err != nil {
		panic(err)
	}

	projectPg, err := api.GetProjectEnvs(idReq.Id, paginationReq.PageNum, paginationReq.PageSize, paginationReq.QueryField)
	if err != nil {
		panic(err)
	}

	resp := PaginationProjectEnvResp{}
	if err := utils.MarshalResponse(projectPg, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get all project_envs, this pagination project number is: %v", len(resp.Items))
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Delete specific project env by id
// @Description api for delete specific project env by id
// @Tags PROJECT
// @Accept json
// @Produce json
// @Param id path integer true "Project ID"
// @Param pe_id path integer true "Project Env ID"
// @Success 204 {object} string "StatusNoContent"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/projects/{id}/project-envs/{pe_id} [delete]
func (c *Controller) DeleteProjectEnv(ctx *gin.Context) {
	var uriReq ProjectIdAndEnvIdReq
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	err := api.DeleteProjectEnvById(uriReq.ProjectEnvId)
	if err != nil {
		panic(err)
	}

	c.Logger.Infof("Delete project(%v)'s env(%v) successfully", uriReq.Id, uriReq.ProjectEnvId)
	ctx.JSON(http.StatusNoContent, "")
}

// @Summary Get project env by id
// @Description api for get project env by id
// @Tags PROJECT
// @Accept json
// @Produce json
// @Param id path integer true "Project ID"
// @Param pe_id path integer true "Project Env ID"
// @Success 200 {object} v1.ProjectEnvBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/projects/{id}/project-envs/{pe_id} [get]
func (c *Controller) GetProjectEnv(ctx *gin.Context) {
	var uriReq ProjectIdAndEnvIdReq
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	projectEnv, err := api.GetProjectEnvById(uriReq.ProjectEnvId)
	if err != nil {
		panic(err)
	}

	resp := ProjectEnvBriefResp{}
	if err := utils.MarshalResponse(projectEnv, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get project(%v) by env id(%v)", uriReq.Id, uriReq.ProjectEnvId)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Update project env description by id and body
// @Description api for update project env description
// @Tags PROJECT
// @Accept json
// @Produce json
// @Param id path integer true "Project ID"
// @Param pe_id path integer true "Project Env ID"
// @Param update_body body v1.UpdateProjectEnvReq true "JSON type for update project env description"
// @Success 200 {object} v1.ProjectEnvBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/projects/{id}/project-envs/{pe_id} [put]
func (c *Controller) UpdateProjectEnv(ctx *gin.Context) {
	var uriReq ProjectIdAndEnvIdReq
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	var descReq UpdateProjectEnvReq
	if err := ctx.ShouldBindJSON(&descReq); err != nil {
		panic(err)
	}
	projectEnv, err := api.UpdateProjectEnv(uriReq.ProjectEnvId, descReq.Description)
	if err != nil {
		panic(err)
	}

	resp := ProjectEnvBriefResp{}
	if err := utils.MarshalResponse(projectEnv, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Update project env description by id:%v", uriReq.Id)
	ctx.JSON(http.StatusOK, resp)
}
