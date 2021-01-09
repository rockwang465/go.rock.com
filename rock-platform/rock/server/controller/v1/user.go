package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	api "go.rock.com/rock-platform/rock/server/database/api"
	middleWare "go.rock.com/rock-platform/rock/server/middleware"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
)

// 对用户进行操作，如新增用户、删除用户、更新用户

type CreateUserReq struct {
	Name     string `json:"name" binding:"required" example:"admin_user"`
	Password string `json:"password" binding:"required" example:"********"`
	Email    string `json:"email" binding:"required" example:"admin_user@sensetime.com"`
	//RoleId   *RoleIdReq `json:"role_id" binding:"required"`  // 用顺义的这种定义，ctx.ShouldBind报错
	RoleId int64 `json:"role_id" binding:"required" example:"1"`
}

// @Summary Create user
// @Description Api to create user
// @Tags USER
// @Accept  json
// @Produce  json
// @Param input_body body v1.RegisterUserInfo true  "JSON type input body"
// @Success 201 {object} v1.UserDetailResp
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
	fmt.Println("api.CreateUser:", userReq.Name, userReq.Password, userReq.Email, userReq.RoleId)

	user, err := api.CreateUser(userReq.Name, userReq.Password, userReq.Email, userReq.RoleId)
	if err != nil {
		panic(err)
		return
	}

	if err := utils.SendNewPwdEmail(user.Name, user.Email, userReq.Password); err != nil {
		panic(err)
		return
	}
	c.Logger.Debugf("Send create User(%s)'s email successfully", user.Name)

	token, err := middleWare.GenerateToken(user.Id, user.Name, user.Password)
	if err != nil {
		panic(err)
		return
	}

	user, err = api.UpdateUserToken()

	//resp := &UserDetailResp{}
	//err = utils.MarshalResponse(user, resp)
	//if err != nil {
	//	panic(err)
	//	return
	//}
	c.Logger.Infof("User %v register successful", user.Name)
	//ctx.JSON(http.StatusOK, resp)
	ctx.JSON(http.StatusOK, user)
}
