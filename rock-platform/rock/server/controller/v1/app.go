package v1

import (
	"github.com/gin-gonic/gin"
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

type AppPagination struct {
	PageNum int64           `json:"page_num" binding:"required" example:"1"`
	PerSize int64           `json:"per_size" binding:"required" example:"10"`
	Total   int64           `json:"total" binding:"required" example:"100"`
	Pages   int64           `json:"pages" binding:"required" example:"1"`
	Items   []*AppBriefResp `json:"items" binding:"required"`
}

type GetAppsPaginationReq struct {
	GetPaginationReq
	Id int64 `json:"id" form:"id" binding:"required,min=1" example:"1"`
}

// @Summary Create app
// @Description Api to create app
// @Tags APP
// @Accept json
// @Produce json
// @Param input_body body v1.CreateAppReq true "JSON type input body"
// @Success 201 {object} v1.AppBriefResp "StatusCreated"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/apps [post]
func (c *Controller) CreateApp(ctx *gin.Context) {
	var createApp CreateAppReq
	if err := ctx.ShouldBindJSON(&createApp); err != nil {
		panic(err)
	}

	app := new(models.App)
	var err error
	if createApp.GitlabProjectId == 0 {
		app, err = api.CreateApp(createApp.Name, "", "", createApp.Description, "", createApp.ProjectId, 0, 0)
		if err != nil {
			panic(err)
		}
	} else {
		// 后期补上:
		// fullName owner droneRepoId 是通过drone，使用用户的access token获取到的
		// gitlabAddr 也是通过drone， 通过drone的token，获取gitlab的地址等
		app, err = api.CreateApp(createApp.Name, "", "", createApp.Description, "", createApp.ProjectId, createApp.GitlabProjectId, 0)
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
// @Param per_size query integer true "App number size" default(10)
// @Param query_field query string false "Fuzzy Query(field: name)"
// @Success 200 {object} v1.ProjectPagination "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
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

	resp := AppPagination{}
	if err := utils.MarshalResponse(appPg, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get all apps, this pagination app number is: %v", len(resp.Items))
	ctx.JSON(http.StatusOK, resp)
}
