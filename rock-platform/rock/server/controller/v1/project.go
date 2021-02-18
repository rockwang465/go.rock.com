package v1

import (
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
)

type CreateProjectReq struct {
	Name        string `json:"name" binding:"required" example:"test-project"`
	Description string `json:"description" binding:"omitempty,max=100" example:"description the project"`
}

type ProjectBriefResp struct {
	Id          int64            `json:"id" example:"1"`
	Name        string           `json:"name" example:"test-project"`
	Description string           `json:"description" example:"description the project"`
	CreatedAt   models.LocalTime `json:"created_at" example:"2021-01-28 20:20:20"`
	UpdateAt    models.LocalTime `json:"updated_at" example:"2021-01-28 20:20:20"`
}

type PaginationProjectResp struct {
	PageNum  int64               `json:"page_num" binding:"required" example:"1"`
	PageSize int64               `json:"page_size" binding:"required" example:"10"`
	Total    int64               `json:"total" binding:"required" example:"100"`
	Pages    int64               `json:"pages" binding:"required" example:"1"`
	Items    []*ProjectBriefResp `json:"items" binding:"required"`
}

type UpdateProjectReq struct {
	Description string `json:"description" binding:"omitempty,max=100" example:"description the project"`
}

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

// @Summary Create project
// @Description Api to create project
// @Tags PROJECT
// @Accept json
// @Produce json
// @Param input_body body v1.CreateProjectReq true "JSON type input body"
// @Success 201 {object} v1.ProjectBriefResp "StatusCreated"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/projects [post]
func (c *Controller) CreateProject(ctx *gin.Context) {
	var createProjectReq CreateProjectReq
	if err := ctx.ShouldBindJSON(&createProjectReq); err != nil {
		panic(err)
	}

	project, err := api.CreateProject(createProjectReq.Name, createProjectReq.Description)
	if err != nil {
		panic(err)
	}

	resp := ProjectBriefResp{}
	if err := utils.MarshalResponse(project, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Create project by id:%v name:%v", resp.Id, resp.Name)
	ctx.JSON(http.StatusCreated, resp)
}

// @Summary Get all projects
// @Description Api to get all projects
// @Tags PROJECT
// @Accept json
// @Produce json
// @Param page_num query integer true "Request page number" default(1)
// @Param page_size query integer true "Project number size" default(10)
// @Param query_field query string false "Fuzzy Query(field: name)"
// @Success 200 {object} v1.PaginationProjectResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/projects [get]
func (c *Controller) GetProjects(ctx *gin.Context) {
	var paginationReq GetPaginationReq
	if err := ctx.ShouldBind(&paginationReq); err != nil {
		panic(err)
	}

	projectPg, err := api.GetProjects(paginationReq.PageNum, paginationReq.PageSize, paginationReq.QueryField)
	if err != nil {
		panic(err)
	}

	resp := PaginationProjectResp{}
	if err := utils.MarshalResponse(projectPg, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get all projects, this pagination project number is: %v", len(resp.Items))
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get a project by id
// @Description Api to get a project by id
// @Tags PROJECT
// @Accept json
// @Produce json
// @Param id query integer true "Project ID"
// @Success 200 {object} v1.ProjectBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/projects/{id} [get]
func (c *Controller) GetProject(ctx *gin.Context) {
	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	project, err := api.GetProjectById(idReq.Id)
	if err != nil {
		panic(err)
	}

	resp := ProjectBriefResp{}
	if err := utils.MarshalResponse(project, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get project by id:%v", resp.Id)
	ctx.JSON(http.StatusOK, resp)

}

// @Summary Delete a project by id
// @Description Api to delete a project by id
// @Tags PROJECT
// @Accept json
// @Produce json
// @Param id body integer true "Project ID"
// @Success 204 {object} string "StatusNoContent"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/projects/{id} [delete]
func (c *Controller) DeleteProject(ctx *gin.Context) {
	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	if err := api.DeleteProjectById(idReq.Id); err != nil {
		panic(err)
	}
	c.Logger.Infof("Delete project by id:%v", idReq.Id)
	ctx.JSON(http.StatusNoContent, "")
}

// @Summary Update project description by id and body
// @Description api for update project description
// @Tags PROJECT
// @Accept json
// @Produce json
// @Param id path integer true "Project ID"
// @Param update_body body v1.UpdateProjectReq true "JSON type for update project description"
// @Success 200 {object} v1.ProjectBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/projects/{id} [put]
func (c *Controller) UpdateProject(ctx *gin.Context) {
	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	var descReq UpdateProjectReq
	if err := ctx.ShouldBindJSON(&descReq); err != nil {
		panic(err)
	}
	project, err := api.UpdateProject(idReq.Id, descReq.Description)
	if err != nil {
		panic(err)
	}

	resp := ProjectBriefResp{}
	if err := utils.MarshalResponse(project, &resp); err != nil {
		panic(err)
	}
	c.Logger.Infof("Update project's description by id:%v", idReq.Id)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get all apps by species project id
// @Description api for get all apps by project id
// @Tags PROJECT
// @Accept json
// @Produce json
// @Param id path integer true "Project ID"
// @Param page_num query integer true "Request page number" default(1)
// @Param page_size query integer true "Request page size" default(10)
// @Success 200 {object} v1.PaginateAppResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/projects/{id}/apps [get]
func (c *Controller) GetProjectApps(ctx *gin.Context) {
	var idReq IdReq // project id
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	var paginationReq GetPaginationReq
	if err := ctx.ShouldBind(&paginationReq); err != nil {
		panic(err)
	}

	_, err := api.GetProjectById(idReq.Id)
	if err != nil {
		panic(err)
	}
	// get apps by project id
	appPg, err := api.GetAppsByProjectId(idReq.Id, paginationReq.PageNum, paginationReq.PageSize, paginationReq.QueryField)
	if err != nil {
		panic(err)
	}

	resp := PaginateAppResp{}
	if err := utils.MarshalResponse(appPg, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get all apps with project id %v", idReq.Id)
	ctx.JSON(http.StatusOK, resp)
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
