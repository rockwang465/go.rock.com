package v1

import (
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
)

type CreateProjectReq struct {
	Name        string `json:"name" binding:"required" example:"senseguard-example"`
	Description string `json:"description" binding:"omitempty,max=100" example:"This is an example"`
}

type ProjectBriefResp struct {
	Id          int64            `json:"id" example:"1"`
	Name        string           `json:"name" example:"senseguard-example"`
	Description string           `json:"description" example:"This is an example"`
	CreatedAt   models.LocalTime `json:"created_at" example:"2021-01-28 20:20:20"`
	UpdateAt    models.LocalTime `json:"updated_at" example:"2021-01-28 20:20:20"`
}

type ProjectPagination struct {
	PageNum int64               `json:"page_num" binding:"required" example:"1"`
	PerSize int64               `json:"per_size" binding:"required" example:"10"`
	Total   int64               `json:"total" binding:"required" example:"100"`
	Pages   int64               `json:"pages" binding:"required" example:"1"`
	Items   []*ProjectBriefResp `json:"items" binding:"required"`
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
// @Param per_size query integer true "Project number size" default(10)
// @Param query_field query string false "Fuzzy Query(field: name)"
// @Success 200 {object} v1.ProjectPagination "StatusOK"
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

	resp := ProjectPagination{}
	if err := utils.MarshalResponse(projectPg, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get all projects, this pagination project number is: %v", len(resp.Items))
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get a project
// @Description Api to get a project
// @Tags PROJECT
// @Accept json
// @Produce json
// @Param id query integer true "Project ID"
// @Success 200 {object} v1.ProjectBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/project/{id} [get]
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
