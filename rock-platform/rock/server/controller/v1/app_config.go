package v1

import (
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
)

type UpdateAppConfReq struct {
	ProjectEnvId int64  `json:"project_env_id" binding:"required,min=1" example:"1"`
	Config       string `json:"config" binding:"required" example:"the app helm chart values.yaml config content"`
}

type AppConfBriefResp struct {
	Id        int64            `json:"id" binding:"required" example:"1"`
	Config    string           `json:"config" binding:"required" example:"app config content"`
	CreatedAt models.LocalTime `json:"created_at" binding:"required" example:"2018-10-09T14:57:23+08:00"`
	UpdatedAt models.LocalTime `json:"updated_at" binding:"required" example:"2018-10-09T14:57:23+08:00"`
	Version   int64            `json:"version" binding:"required" example:"1"`
}

type AppConfDetailResp struct {
	AppConfBriefResp
	AppId        int64 `json:"app_id" example:"1"`
	ProjectEnvId int64 `json:"project_env_id" example:"1"`
}

type AppProjectEnvIdReq struct {
	ProjectEnvId int64 `json:"pe_id" form:"pe_id" binding:"required,min=1" example:"1"`
}

// @Summary Update or create specific app conf
// @Description api for update or create specific app conf by app_id and project_env_id
// @Tags APP
// @Accept json
// @Produce json
// @Param id path integer true "App ID"
// @Param update_body body v1.UpdateAppConfReq true "JSON type for update or create app info"
// @Success 200 {object} v1.AppConfDetailResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/apps/{id}/config [put]
func (c *Controller) UpdateOrCreateAppConf(ctx *gin.Context) {
	// 通过 project_env_id + app_id 来定位 app_conf 表的数据，然后将配置(values.yaml)更新到app_conf表中
	var uriReq IdReq // app_id
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	var req UpdateAppConfReq // sed s/$/"\\\n"/ override.yaml | tr -d "\n"
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}

	// verify the YAML format
	if err := utils.YamlValidator(req.Config); err != nil {
		panic(err)
	}

	appConf, err := api.UpdateOrCreateAppConfById(uriReq.Id, req.ProjectEnvId, req.Config)
	if err != nil {
		panic(err)
	}

	resp := AppConfDetailResp{}
	if err := utils.MarshalResponse(appConf, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Update or create app_conf by app_id(%v) and project_env_id(%v)", uriReq.Id, req.ProjectEnvId)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary get specific app config
// @Description api for get specific app config by app_id and project_env_id
// @Tags APP
// @Accept json
// @Produce json
// @Param id path integer true "App ID"
// @Param pe_id query integer true "Project env ID"
// @Success 200 {object} v1.AppConfDetailResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/apps/{id}/config [get]
func (c *Controller) GetAppConf(ctx *gin.Context) {
	var uriReq IdReq // app_conf_id
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	var paramsReq AppProjectEnvIdReq // project_env_id
	if err := ctx.ShouldBind(&paramsReq); err != nil {
		panic(err)
	}

	appConf, err := api.GetAppConfByAppAndProjectEnvId(uriReq.Id, paramsReq.ProjectEnvId)
	if err != nil {
		panic(err)
	}

	resp := AppConfDetailResp{}
	if err := utils.MarshalResponse(appConf, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get app_conf config by app_id(%v) and project_env_id(%v)", uriReq.Id, paramsReq.ProjectEnvId)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary delete specific app config
// @Description api for delete specific app config by app_id and project_env_id
// @Tags APP
// @Accept json
// @Produce json
// @Param id path integer true "App ID"
// @Param pe_id query integer true "Project env ID"
// @Success 204 {object} string "StatusNoContent"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/apps/{id}/config [delete]
func (c *Controller) DeleteAppConf(ctx *gin.Context) {
	var uriReq IdReq // app_conf_id
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	var paramsReq AppProjectEnvIdReq // project_env_id
	if err := ctx.ShouldBind(&paramsReq); err != nil {
		panic(err)
	}
	if err := api.DeleteAppConfByProjectAndAppId(uriReq.Id, paramsReq.ProjectEnvId); err != nil {
		panic(err)
	}

	c.Logger.Infof("Delete app_conf config by app_id(%v) and project_env_id(%v)", uriReq.Id, paramsReq.ProjectEnvId)
	ctx.JSON(http.StatusNoContent, "")
}
