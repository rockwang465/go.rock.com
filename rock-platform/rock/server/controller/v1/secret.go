package v1

import (
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/clients/drone-api"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
)

type CreateSecretReq struct {
	Name  string `json:"name"  binding:"required" example:"docker_user"`
	Value string `json:"value" binding:"required" example:"admin"`
}

type SecretBriefResp struct {
	Name string `json:"name"  example:"docker_user"`
}

type SecretByNameReq struct {
	Name string `json:"name" uri:"name" binding:"required" example:"docker_user"`
}
type UpdateSecretReq struct {
	Value string `json:"value" binding:"required" example:"admin"`
}

// @Summary Create secret
// @Description Api to create secret
// @Tags SECRET
// @Accept json
// @Produce json
// @Param input_body body v1.CreateSecretReq true "JSON type input body"
// @Success 201 {object} v1.SecretBriefResp
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/secrets [post]
func (c *Controller) CreateSecret(ctx *gin.Context) {
	// secret 是给 .drone.yaml 中需要的基础变量，最终给 drone-agent 进行镜像拉取 和 运维平台api(CreateDeployment应用发版到指定环境)请求
	// docker_username = admin  // harbor用户
	// docker_password = Se*****5  // harbor密码
	// galaxias_api_token = eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFkbWluIiwiZHJvbmVfdG9rZW4iOiIiLCJwYXNzd29yZCI6IjMyMDdlYWQ0ZTA5MmRlNzdlMDIyMzk0YjMyMDRkNzU1Iiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNjE0MDc3NjAzLCJpYXQiOjE2MTQwNzE1NDMsImlzcyI6IlJvY2sgV2FuZyIsInN1YiI6IkxvZ2luIHRva2VuIn0.nAMR3xjGZ-4etgyVT2qfiUx2oEZhKM_iRs8lui1vTJ4  // 运维平台admin用户jwt token
	cfgCtx, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	var req CreateSecretReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}

	secret, err := drone_api.CreateSecret(cfgCtx.DroneToken, req.Name, req.Value)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%#v\n", secret) // &drone.Secret{ID:1, Name:"docker_name", Value:"", Images:[]string(nil), Events:[]string(nil)}

	resp := SecretBriefResp{}
	if err := utils.MarshalResponse(secret, &resp); err != nil {
		panic(err)
	}

	c.Infof("Secret with name: %s created successfully in drone", req.Name)
	ctx.JSON(http.StatusCreated, resp)
}

// @Summary Get all secrets
// @Description Api to get all secrets
// @Tags SECRET
// @Accept json
// @Produce json
// @Success 200 {array} v1.SecretBriefResp
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/secrets [get]
func (c *Controller) GetSecrets(ctx *gin.Context) {
	cfgCtx, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	secrets, err := drone_api.GetSecrets(cfgCtx.DroneToken)
	if err != nil {
		panic(err)
	}

	resp := []*SecretBriefResp{}
	if err := utils.MarshalResponse(secrets, &resp); err != nil {
		panic(err)
	}

	c.Infof("Get all secrets, length is %v", len(resp))
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get an secret
// @Description Api to get an secret
// @Tags SECRET
// @Accept json
// @Produce json
// @Param name path string true "Secret Name"
// @Success 200 {object} v1.SecretBriefResp
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/secret/{name} [get]
func (c *Controller) GetSecret(ctx *gin.Context) {
	cfgCtx, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	var uriReq SecretByNameReq
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	sercret, err := drone_api.GetSecret(cfgCtx.DroneToken, uriReq.Name)
	if err != nil {
		panic(err)
	}

	resp := SecretBriefResp{}
	if err := utils.MarshalResponse(sercret, &resp); err != nil {
		panic(err)
	}

	c.Infof("Get an secret by name %v", uriReq.Name)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Delete an secret by name
// @Description Api to delete an secret by name
// @Tags SECRET
// @Accept json
// @Produce json
// @Param name path string true "Secret Name"
// @Success 204 {object} string "StatusNoContent"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/secret/{name} [delete]
func (c *Controller) DeleteSecret(ctx *gin.Context) {
	cfgCtx, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	var uriReq SecretByNameReq
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	if err := drone_api.DeleteSecret(cfgCtx.DroneToken, uriReq.Name); err != nil {
		panic(err)
	}

	c.Infof("Delete an secret by name %v", uriReq.Name)
	ctx.JSON(http.StatusNoContent, "")
}

// @Summary Delete an secret by name
// @Description Api to delete an secret by name
// @Tags SECRET
// @Accept json
// @Produce json
// @Param name path string true "Secret Name"
// @Param update_body body v1.UpdateSecretReq true "JSON type update secret info"
// @Success 204 {object} v1.SecretBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/secret/{name} [put]
func (c *Controller) UpdateSecret(ctx *gin.Context) {
	cfgCtx, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	var uriReq SecretByNameReq
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	var req UpdateSecretReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}

	secret, err := drone_api.UpdateSecret(cfgCtx.DroneToken, uriReq.Name, req.Value)
	if err != nil {
		panic(err)
	}

	resp := SecretBriefResp{}
	if err := utils.MarshalResponse(secret, &resp); err != nil {
		panic(err)
	}

	c.Infof("Update secret's info by name %v", uriReq.Name)
	ctx.JSON(http.StatusOK, resp)
}
