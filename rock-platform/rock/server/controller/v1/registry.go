package v1

import (
	"github.com/gin-gonic/gin"
	api "go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/database/models"
	"net/http"
)

// 用户注册页面

type RegisterUserInfo struct {
	Username string `json:"username" binding:"required" example:"admin_user"`
	Password string `json:"password" binding:"required" example:"********"`
	Email    string `json:"email" binding:"required" example:"admin_user@sensetime.com"`
}

type UserDetailResp struct {
	Id        int64            `json:"id" example:"1"`
	Username  string           `json:"username" example:"admin_user"`
	Email     string           `json:"email" example:"admin_user@sensetime.com"`
	CreatedAt models.LocalTime `json:"created_at" example:"2020-12-20 15:15:22"`
	UpdatedAt models.LocalTime `json:"updated_at" example:"2020-12-20 15:15:22"`
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
// @Router /v1/register [post]
func (c *Controller) Register(ctx *gin.Context) {
	var user RegisterUserInfo

	if err := ctx.ShouldBind(&user); err != nil {
		panic(err)
		return
	}

	if len(user.Password) < 6 {
		panic("The password length is too short, greater than or equal 6")
		return
	}

	userInfo, err := api.RegistryCreateUser(user.Username, user.Password, user.Email)
	if err != nil {
		panic(err)
	}

	var resp *UserDetailResp
	resp = &UserDetailResp{
		Id:        userInfo.Id,
		Username:  userInfo.Name,
		Email:     userInfo.Email,
		CreatedAt: userInfo.CreatedAt,
		UpdatedAt: userInfo.UpdatedAt,
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message":   "register is success",
		"user_info": resp,
	})
}
