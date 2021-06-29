package v1

// record namespace and cluster_id for a k8s cluster

import (
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
)

type CreateEnvReq struct {
	Namespace   string `json:"namespace" binding:"required" example:"namespace of k8s cluster"`
	Description string `json:"description" binding:"omitempty,max=100" example:"description for env"`
	ClusterId   int64  `json:"cluster_id" binding:"required,min=1" example:"1"`
}

type EnvBriefResp struct {
	Id          int64            `json:"id" example:"1"`
	Namespace   string           `json:"namespace" example:"namespace of k8s cluster"`
	Description string           `json:"description" example:"description for env"`
	ClusterId   int64            `json:"cluster_id" example:"1"`
	CreatedAt   models.LocalTime `json:"created_at" example:"2021-01-28 20:20:20"`
	UpdatedAt   models.LocalTime `json:"updated_at" example:"2021-01-28 20:20:20"`
	Version     int64            `json:"version" example:"1"`
}

type PaginationEnvResp struct {
	PageNum  int64           `json:"page_num" binding:"required" example:"1"`
	PageSize int64           `json:"page_size" binding:"required" example:"10"`
	Total    int64           `json:"total" binding:"required" example:"100"`
	Pages    int64           `json:"pages" binding:"required" example:"1"`
	Items    []*EnvBriefResp `json:"items" binding:"required"`
}

type UpdateEnvReq struct {
	Description string `json:"description" binding:"omitempty,max=100" example:"description for env"`
}

// @Summary Create env
// @Description Api to create env
// @Tags ENV
// @Accept json
// @Produce json
// @Param input_body body v1.CreateEnvReq true "JSON type input body"
// @Success 201 {object} v1.EnvBriefResp "StatusCreated"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/envs [post]
func (c *Controller) CreateEnv(ctx *gin.Context) {
	var req CreateEnvReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}

	env, err := api.CreateEnv(req.Namespace, req.Description, req.ClusterId)
	if err != nil {
		panic(err)
	}

	resp := EnvBriefResp{}
	if err := utils.MarshalResponse(env, &resp); err != nil {
		panic(err)
	}
	c.Logger.Infof("Create env by id:%v namespace:%v", resp.Id, resp.Namespace)
	ctx.JSON(http.StatusCreated, resp)
}

// @Summary Get all envs
// @Description Api to get all env
// @Tags ENV
// @Accept json
// @Produce json
// @Param page_num query integer true "Request page number" default(1)
// @Param page_size query integer true "Env number size" default(10)
// @Success 200 {object} v1.PaginationEnvResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/envs [get]
func (c *Controller) GetEnvs(ctx *gin.Context) {
	var paginationReq GetPaginationReq
	if err := ctx.ShouldBind(&paginationReq); err != nil {
		panic(err)
	}

	// get all env, no query field
	envPg, err := api.GetEnvs(paginationReq.PageNum, paginationReq.PageSize)
	if err != nil {
		panic(err)
	}

	resp := PaginationEnvResp{}
	if err := utils.MarshalResponse(envPg, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get all env, this pagination env number is: %v", len(resp.Items))
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get an env by id
// @Description api for get an env by id
// @Tags ENV
// @Accept json
// @Produce json
// @Param id path integer true "Env ID"
// @Success 200 {object} v1.EnvBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/envs/{id} [get]
func (c *Controller) GetEnv(ctx *gin.Context) {
	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	env, err := api.GetEnvById(idReq.Id)
	if err != nil {
		panic(err)
	}

	resp := EnvBriefResp{}
	if err := utils.MarshalResponse(env, &resp); err != nil {
		panic(err)
	}
	c.Logger.Infof("Get env by id:%v", idReq.Id)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get an env by id
// @Description api for get an env by id
// @Tags ENV
// @Accept json
// @Produce json
// @Param id path integer true "Env ID"
// @Success 204 {object} string "StatusNoContent"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/envs/{id} [delete]
func (c *Controller) DeleteEnv(ctx *gin.Context) {
	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}
	if err := api.DeleteEnvById(idReq.Id); err != nil {
		panic(err)
	}
	c.Logger.Infof("Delete env by id:%v", idReq.Id)
	ctx.JSON(http.StatusNoContent, "")
}

// @Summary Update env description by id and body
// @Description api for update env description
// @Tags ENV
// @Accept json
// @Produce json
// @Param id path integer true "Env ID"
// @Param update_body body v1.UpdateEnvReq true "JSON type for update env description"
// @Success 200 {object} v1.EnvBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/envs/{id} [put]
func (c *Controller) UpdateEnv(ctx *gin.Context) {
	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	var descReq UpdateEnvReq
	if err := ctx.ShouldBindJSON(&descReq); err != nil {
		panic(err)
	}
	env, err := api.UpdateEnv(idReq.Id, descReq.Description)
	if err != nil {
		panic(err)
	}

	resp := EnvBriefResp{}
	if err := utils.MarshalResponse(env, &resp); err != nil {
		panic(err)
	}
	c.Logger.Infof("Update env's description by id:%v", idReq.Id)
	ctx.JSON(http.StatusOK, resp)
}
