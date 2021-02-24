package v1

import (
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/clients/drone-api"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
)

type CreateRegistryReq struct {
	Address  string `json:"address"  binding:"required" example:"10.151.3.75"`
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"*********"`
}

type RegistryBriefResp struct {
	Address  string `json:"address"  binding:"required" example:"10.151.3.75"`
	Username string `json:"username" binding:"required" example:"admin"`
}

type RegistryByAddrReq struct {
	Address string `json:"address" uri:"address" binding:"required" example:"10.151.3.75"`
}

type UpdateRegistryReq struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"admin"`
}

// @Summary Create registry
// @Description Api to create registry
// @Tags REGISTRY
// @Accept json
// @Produce json
// @Param input_body body v1.CreateRegistryReq true "JSON type input body"
// @Success 201 {object} v1.RegistryBriefResp
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/registries [post]
func (c *Controller) CreateRegistry(ctx *gin.Context) {
	cfgCtx, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	var req CreateRegistryReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}

	registry, err := drone_api.CreateRegistry(cfgCtx.DroneToken, req.Address, req.Username, req.Password)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%#v\n", registry) // &drone.Registry{ID:0, Address:"10.151.3.75", Username:"admin", Password:"Sense*******", Email:"", Token:""}

	resp := RegistryBriefResp{}
	if err := utils.MarshalResponse(registry, &resp); err != nil {
		panic(err)
	}

	c.Infof("Registry with docker image repo address: %s created successfully in drone", req.Address)
	ctx.JSON(http.StatusCreated, resp)
}

// @Summary Get all registries
// @Description Api to get all registries
// @Tags REGISTRY
// @Accept json
// @Produce json
// @Success 200 {array} v1.RegistryBriefResp
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/registries [get]
func (c *Controller) GetRegistries(ctx *gin.Context) {
	cfgCtx, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	registries, err := drone_api.GetRegistries(cfgCtx.DroneToken)
	if err != nil {
		panic(err)
	}

	resp := make([]*RegistryBriefResp, 0)
	if err := utils.MarshalResponse(registries, &resp); err != nil {
		panic(err)
	}

	c.Infof("Get all registries, length is %v", len(resp))
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get an registry
// @Description Api to get an registry
// @Tags REGISTRY
// @Accept json
// @Produce json
// @Param address path string true "Registry Address"
// @Success 200 {object} v1.RegistryBriefResp
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/registries/{address} [get]
func (c *Controller) GetRegistry(ctx *gin.Context) {
	cfgCtx, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	var uriReq RegistryByAddrReq
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	registry, err := drone_api.GetRegistry(cfgCtx.DroneToken, uriReq.Address)
	if err != nil {
		panic(err)
	}

	resp := RegistryBriefResp{}
	if err := utils.MarshalResponse(registry, &resp); err != nil {
		panic(err)
	}

	c.Infof("Get an registry by address %v", uriReq.Address)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Delete an registry by address
// @Description Api to delete an registry by address
// @Tags REGISTRY
// @Accept json
// @Produce json
// @Param address path string true "Registry Address"
// @Success 204 {object} string "StatusNoContent"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/registries/{address} [delete]
func (c *Controller) DeleteRegistry(ctx *gin.Context) {
	cfgCtx, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	var uriReq RegistryByAddrReq
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	if err := drone_api.DeleteRegistry(cfgCtx.DroneToken, uriReq.Address); err != nil {
		panic(err)
	}

	c.Infof("Delete an registry by address %v", uriReq.Address)
	ctx.JSON(http.StatusNoContent, "")
}

// @Summary Update registry info with id and body
// @Description Api for update registry info
// @Tags REGISTRY
// @Accept json
// @Produce json
// @Param address path string true "Registry Address"
// @Param update_body body v1.UpdateRegistryReq true "JSON type update registry info
// @Success 204 {object} v1.RegistryBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/registries/{address} [put]
func (c *Controller) UpdateRegistry(ctx *gin.Context) {
	cfgCtx, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	var uriReq RegistryByAddrReq
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	var req UpdateRegistryReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}

	registry, err := drone_api.UpdateRegistry(cfgCtx.DroneToken, uriReq.Address, req.Username, req.Password)
	if err != nil {
		panic(err)
	}

	resp := RegistryBriefResp{}
	if err := utils.MarshalResponse(registry, &resp); err != nil {
		panic(err)
	}

	c.Infof("Update registry's info by address %v", uriReq.Address)
	ctx.JSON(http.StatusOK, resp)
}
