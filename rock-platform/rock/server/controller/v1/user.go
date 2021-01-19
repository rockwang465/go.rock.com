package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
)

// 对用户进行操作，如新增用户、删除用户、更新用户

type CreateUserReq struct {
	Name     string `json:"name" binding:"required" example:"admin_user"`
	Password string `json:"password" binding:"required" example:"********"`
	Email    string `json:"email" binding:"required" example:"admin_user@sensetime.com"`
	//RoleId   *RoleIdReq `json:"role_id" binding:"required"`  // 用顺义的这种定义，ctx.ShouldBind报错
	RoleId int64 `json:"role_id" binding:"required" example:"1"` // role表id=1
}

// @Summary Create user
// @Description Api to create user
// @Tags USER
// @Accept json
// @Produce json
// @Param input_body body v1.CreateUserReq true "JSON type input body"
// @Success 200 {object} api.UserDetailResp
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/users [post]
func (c *Controller) CreateUser(ctx *gin.Context) {
	var userReq CreateUserReq
	if err := ctx.ShouldBind(&userReq); err != nil {
		panic(err)
		return
	}
	if len(userReq.Password) < 6 {
		err := utils.NewRockError(400, 40000002, fmt.Sprintf("The password length is too short, greater than or equal 6")) // generate a error
		panic(err)
		return
	}

	user, role, err := api.CreateUser(userReq.Name, userReq.Password, userReq.Email, userReq.RoleId)
	if err != nil {
		panic(err)
		return
	}
	if err := utils.SendNewPwdEmail(user.Name, user.Email, userReq.Password); err != nil {
		panic(err)
		return
	}
	c.Logger.Debugf("Send create User(%s)'s email successfully", user.Name)

	token, err := utils.GenerateToken(user.Id, user.Name, user.Password, role.Name)
	if err != nil {
		panic(err)
		return
	}

	user, err = api.UpdateUserToken(user.Id, token)
	if err != nil {
		panic(err)
		return
	}

	resp, err := api.GetUserDetailResp(user.Id)
	if err != nil {
		panic(err)
		return
	}

	c.Logger.Infof("User %v register successful", user.Name)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get users
// @Description Api to get all users
// @Tags USER
// @Accept json
// @Produce json
// @Param PageNum query integer true "Request page number" default(1)
// @Param PageSize query integer true "Request page size" default(10)
// @Param QueryField query string false "Fuzzy Query(field: name)"
// @Success 200 {object} models.UserPagination
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/users [get]
func (c *Controller) GetUsers(ctx *gin.Context) {
	var paginationReq GetPaginationReq
	err := ctx.ShouldBind(&paginationReq)
	if err != nil {
		newErr := fmt.Sprintf("context bind failed, %v", err.Error())
		panic(newErr)
	}
	userPg, err := api.GetUsers(paginationReq.PageNum, paginationReq.PageSize, paginationReq.QueryField)
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, userPg)
}

// @Summary Get user with id
// @Description Api to get user with id
// @Tags USER
// @Accept json
// @Produce json
// @Param id path integer true "User ID"
// @Success 200 {object} api.UserDetailResp
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/users/{id} [get]
func (c *Controller) GetUser(ctx *gin.Context) {
	var getIdReq IdReq
	err := ctx.ShouldBindUri(&getIdReq)
	if err != nil {
		newErr := fmt.Sprintf("context bind failed, %v", err.Error())
		panic(newErr)
	}

	resp, err := api.GetUserDetailResp(getIdReq.Id)
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, resp)
}
