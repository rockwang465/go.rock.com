package v1

import (
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/clients/drone-api"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
)

type CreateAppReq struct {
	Name            string `json:"name" binding:"required" example:"test_app"`                            // 新建的app的名称
	Description     string `json:"description" binding:"omitempty,max=100" example:"description the app"` // 描述这个app
	ProjectId       int64  `json:"project_id" binding:"required,min=1" example:"1"`                       // 属于哪个project(id关联)
	GitlabProjectId int64  `json:"gitlab_project_id" binding:"omitempty,min=1" example:"1"`               // 当前用户在gitlab上拥有的project名称(即app的名称)
}

type AppBriefResp struct {
	Id            int64            `json:"id" example:"1"`
	Name          string           `json:"name" example:"test_app"`
	FullName      string           `json:"full_name"  example:"senseguard/test_app"`
	Owner         string           `json:"owner"  example:"1"`
	Description   string           `json:"description"  example:"description the app"`
	GitlabAddress string           `json:"gitlab_address"  example:"http://gitlab.sensetime.com"`
	ProjectId     int64            `json:"project_id"  example:"1"`
	DroneRepoId   int64            `json:"drone_repo_id" example:"1"`
	CreatedAt     models.LocalTime `json:"created_at" example:"2021-01-28 20:20:20"`
	UpdateAt      models.LocalTime `json:"updated_at" example:"2021-01-28 20:20:20"`
	Version       int64            `json:"version" example:"1"`
}

type PaginateAppResp struct {
	PageNum  int64           `json:"page_num" binding:"required" example:"1"`
	PageSize int64           `json:"page_size" binding:"required" example:"10"`
	Total    int64           `json:"total" binding:"required" example:"100"`
	Pages    int64           `json:"pages" binding:"required" example:"1"`
	Items    []*AppBriefResp `json:"items" binding:"required"`
}

type GetAppsPaginationReq struct {
	GetPaginationReq
	Id int64 `json:"id" form:"id" binding:"required,min=1" example:"1"`
}

type UpdateAppReq struct {
	Description string `json:"description" binding:"omitempty,max=100" example:"description for app"`
}

type UpdateAppGitlabAddressReq struct {
	GitlabProjectId int64 `json:"gitlab_project_id" binding:"required,min=1" example:"1"`
}

// @Summary Create app
// @Description Api to create app
// @Tags APP
// @Accept json
// @Produce json
// @Param input_body body v1.CreateAppReq true "JSON type input body"
// @Success 201 {object} v1.AppBriefResp "StatusCreated"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/apps [post]
func (c *Controller) CreateApp(ctx *gin.Context) {
	cfgCtx, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	var createApp CreateAppReq
	if err := ctx.ShouldBindJSON(&createApp); err != nil {
		panic(err)
	}

	app := new(models.App)
	if createApp.GitlabProjectId == 0 {
		app, err = api.CreateApp(createApp.Name, "", "", createApp.Description, "", createApp.ProjectId, 0, 0)
		if err != nil {
			panic(err)
		}
	} else {
		remote, err := drone_api.SyncRemoteRepo(cfgCtx.DroneToken, createApp.GitlabProjectId)
		if err != nil {
			panic(err)
		}
		// example: createApp.GitlabProjectId:9616, remote.ID:9 (drone repo id)
		repo, err := drone_api.ActiveRepo(cfgCtx.DroneToken, remote.ID)
		if err != nil {
			err := utils.NewRockError(403, 40300002, "Permission deny, because you don't have gitlab project master permission")
			panic(err)
		}

		app, err = api.CreateApp(createApp.Name, repo.FullName, repo.Owner, createApp.Description, remote.Clone, createApp.ProjectId, createApp.GitlabProjectId, repo.ID)
		if err != nil {
			panic(err)
		}
	}

	resp := AppBriefResp{}
	if err := utils.MarshalResponse(app, &resp); err != nil {
		panic(err)
	}
	c.Logger.Infof("Create project by id:%v name:%v", resp.Id, resp.Name)
	ctx.JSON(http.StatusCreated, resp)
}

// @Summary Get all apps
// @Description Api to get all apps
// @Tags APP
// @Accept json
// @Produce json
// @Param page_num query integer true "Request page number" default(1)
// @Param page_size query integer true "App number size" default(10)
// @Param query_field query string false "Fuzzy Query(field: name)"
// @Success 200 {object} v1.PaginateAppResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/apps [get]
func (c *Controller) GetApps(ctx *gin.Context) {
	var paginationReq GetAppsPaginationReq
	if err := ctx.ShouldBind(&paginationReq); err != nil {
		panic(err)
	}

	appPg, err := api.GetApps(paginationReq.PageNum, paginationReq.PageSize, paginationReq.QueryField, paginationReq.Id)
	if err != nil {
		panic(err)
	}

	resp := PaginateAppResp{}
	if err := utils.MarshalResponse(appPg, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get all apps, this pagination app number is: %v", len(resp.Items))
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get an app by id
// @Description api for get an app by id
// @Tags APP
// @Accept json
// @Produce json
// @Param id path integer true "App ID"
// @Success 200 {object} v1.AppBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/apps/{id} [get]
func (c *Controller) GetApp(ctx *gin.Context) {
	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	app, err := api.GetAppById(idReq.Id)
	if err != nil {
		panic(err)
	}

	resp := AppBriefResp{}
	if err := utils.MarshalResponse(app, &resp); err != nil {
		panic(err)
	}
	c.Logger.Infof("Get app name:%v by id:%v", app.Name, app.Id)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get an app by id
// @Description api for get an app by id
// @Tags APP
// @Accept json
// @Produce json
// @Param id path integer true "App ID"
// @Success 204 {object} string "StatusNoContent"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/apps/{id} [delete]
func (c *Controller) DeleteApp(ctx *gin.Context) {
	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}
	if err := api.DeleteAppById(idReq.Id); err != nil {
		panic(err)
	}
	c.Logger.Infof("Delete app by id:%v", idReq.Id)
	ctx.JSON(http.StatusNoContent, "")
}

// @Summary Update app description by id and body
// @Description api for update app description
// @Tags APP
// @Accept json
// @Produce json
// @Param id path integer true "App ID"
// @Param update_body body v1.UpdateAppReq true "JSON type for update app description"
// @Success 200 {object} v1.RoleBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/apps/{id} [put]
func (c *Controller) UpdateApp(ctx *gin.Context) {
	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	var descReq UpdateAppReq
	if err := ctx.ShouldBindJSON(&descReq); err != nil {
		panic(err)
	}
	app, err := api.UpdateApp(idReq.Id, descReq.Description)
	if err != nil {
		panic(err)
	}

	resp := RoleBriefResp{}
	if err := utils.MarshalResponse(app, &resp); err != nil {
		panic(err)
	}
	c.Logger.Infof("Update app's description by id:%v", idReq.Id)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Update app gitlab address
// @Description api for update app gitlab address
// @Tags APP
// @Accept json
// @Produce json
// @Param id path integer true "App ID"
// @Param update_body body v1.UpdateAppGitlabAddressReq true "JSON type for update app gitlab address"
// @Success 200 {object} v1.AppBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/apps/{id}/gitlab [put]
func (c *Controller) UpdateAppGitlabProject(ctx *gin.Context) {
	// 应用管理 -> 应用 -> 查看单个应用 -> 关联代码仓库:(此时会自动请求: /v1/repos 拿到所有的gitlab地址,渲染到下拉菜单,并让用户选择需要更改为哪个gitlab地址)
	// 用户 选择一个gitlab 地址后，则开始更新应用的gitlab地址。即当前api的操作。
	cfgCtx, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	var req UpdateAppGitlabAddressReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}

	remote, err := drone_api.SyncRemoteRepo(cfgCtx.DroneToken, req.GitlabProjectId)
	if err != nil {
		panic(err)
	}
	repo, err := drone_api.ActiveRepo(cfgCtx.DroneToken, remote.ID)
	if err != nil {
		err := utils.NewRockError(403, 40300002, "Permission deny, because you don't have gitlab project master permission")
		panic(err)
	}

	app, err := api.UpdateAppGitlabAddressById(idReq.Id, repo.FullName, repo.Owner, remote.Clone, req.GitlabProjectId, repo.ID)
	if err != nil {
		panic(err)
	}

	resp := AppBriefResp{}
	if err := utils.MarshalResponse(app, &resp); err != nil {
		panic(err)
	}
	c.Logger.Infof("Update app gitlab address(%v) by id:%v", remote.Clone, idReq.Id)
	ctx.JSON(http.StatusOK, resp)
}
