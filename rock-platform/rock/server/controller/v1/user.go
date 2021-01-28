package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
)

// 对用户进行操作，如新增用户、删除用户、更新用户

type CreateUserReq struct {
	Name     string `json:"name" binding:"required" example:"admin_user"`
	Password string `json:"password" binding:"required" example:"********"`
	Email    string `json:"email" binding:"required,email" example:"admin_user@sensetime.com"`
	//RoleId   *RoleIdReq `json:"role_id" binding:"required"`  // 用顺义的这种定义，ctx.ShouldBind报错
	RoleId int64 `json:"role_id" binding:"required" example:"1"` // role表id=1
}

type UserBriefResp struct {
	Id        int64            `json:"id" binding:"required" example:"1"`
	Name      string           `json:"name" binding:"required" example:"admin_role"`
	Email     string           `json:"email" binding:"required" example:"admin@sensetime.com"`
	RoleId    int64            `json:"role_id"  binding:"required" example:"1"`
	CreatedAt models.LocalTime `json:"created_at" binding:"required" example:"2018-10-09T14:57:23+08:00"`
	UpdatedAt models.LocalTime `json:"updated_at" binding:"required" example:"2018-10-09T14:57:23+08:00"`
	Version   int              `json:"version" binding:"required" example:"1"`
}

type PaginateBriefUserResp struct {
	PageNum int64            `json:"page_num" binding:"required" example:"1"`
	PerSize int64            `json:"per_size" binding:"required" example:"10"`
	Total   int64            `json:"total" binding:"required" example:"100"`
	Pages   int64            `json:"pages" binding:"required" example:"1"`
	Items   []*UserBriefResp `json:"items" binding:"required"`
}

//type UserDetailResp struct {
//	UserBriefResp
//	RoleId *RoleBriefResp `json:"role_id" binding:"required"`
//}
//
//type PaginateDetailUserResp struct {
//	PageNum int64            `json:"page_num" binding:"required" example:"1"`
//	PerSize int64            `json:"per_size" binding:"required" example:"10"`
//	Total   int64            `json:"total" binding:"required" example:"100"`
//	Pages   int64            `json:"pages" binding:"required" example:"1"`
//	Items   []*UserDetailResp `json:"items" binding:"required"`
//}

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

	err := utils.CheckPwd(userReq.Password)
	if err != nil {
		panic(err)
	}
	has, err := api.HasEmail(userReq.Email)
	if err != nil {
		panic(err)
	}
	if has {
		err = utils.NewRockError(http.StatusBadRequest, 40000009, fmt.Sprintf("Email %v was registered", userReq.Email))
		panic(err)
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

	resp, err := api.GetUserBriefResp(user.Id)
	if err != nil {
		panic(err)
		return
	}

	c.Logger.Infof("User created with id: %v, name: %v", resp.Id, resp.Name)
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

	pgUser := PaginateBriefUserResp{}
	if err := utils.MarshalResponse(userPg, &pgUser); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get all users, this pagination user number is: %v", len(pgUser.Items))
	ctx.JSON(http.StatusOK, pgUser)
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
	var uriIdReq IdReq
	err := ctx.ShouldBindUri(&uriIdReq)
	if err != nil {
		newErr := fmt.Sprintf("context bind failed, %v", err.Error())
		panic(newErr)
	}

	resp, err := api.GetUserBriefResp(uriIdReq.Id)
	if err != nil {
		panic(err)
	}
	c.Logger.Infof("Get user with id: %v", uriIdReq.Id)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Delete user with id
// @Description Api to delete user with id
// @Tags USER
// @Accept json
// @Produce json
// @Param id path integer true "User ID"
// @Success 200 {object} string
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/users/{id} [delete]
func (c *Controller) DeleteUser(ctx *gin.Context) {
	var uriIdReq IdReq
	err := ctx.ShouldBindUri(&uriIdReq)
	if err != nil {
		newErr := fmt.Sprintf("context bind failed, %v", err.Error())
		panic(newErr)
	}

	username, err := api.DeleteUserById(uriIdReq.Id)
	if err != nil {
		panic(err)
	}

	c.Logger.Infof("Delete user: %v with id: %v", username, uriIdReq.Id)
	ctx.JSON(http.StatusOK, fmt.Sprintf("Delete user successful"))
}

// @Summary Update user password with id and old password
// @Description Api to update user password with id and old password
// @Tags USER
// @Accept json
// @Produce json
// @Param Id path integer true "User ID"
// @Param update_body body v1.UpdateUserPwdReq true "JSON body for update user info"
// @Success 200 {object} api.UserDetailResp
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/users/{id}/password [put]
func (c *Controller) UpdateUserPwd(ctx *gin.Context) {
	config, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	var uriIdReq IdReq
	err = ctx.ShouldBindUri(&uriIdReq)
	if err != nil {
		newErr := fmt.Sprintf("context bind failed, %v", err.Error())
		panic(newErr)
	}

	var pwdReq UpdateUserPwdReq
	err = ctx.ShouldBindJSON(&pwdReq)
	if err != nil {
		newErr := fmt.Sprintf("context bind failed, %v", err.Error())
		panic(newErr)
	}

	oldPwd := pwdReq.OldPassword
	newPwd := pwdReq.NewPassword
	user, err := api.UpdateUserPwdById(uriIdReq.Id, oldPwd, newPwd, config.Role)
	if err != nil {
		panic(err)
	}

	if config.Username == user.Name {
		ctx.SetCookie("token", user.Token, utils.GetExpireDuration(), "/", "", false, true)
	}

	resp, err := api.GetUserBriefResp(uriIdReq.Id)
	if err != nil {
		panic(err)
	}

	c.Logger.Infof("Update user with id: %v", uriIdReq.Id)
	ctx.JSON(http.StatusOK, resp)
}
