package v1

import (
	"github.com/gin-gonic/gin"
	api "go.rock.com/rock-platform/rock/server/database/api"
	"net/http"
)

// 用户注册页面

type RegisterUserInfo struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func (c *Controller) Register(ctx *gin.Context) {
	var rui RegisterUserInfo

	if err := ctx.ShouldBind(&rui); err != nil {
		panic(err)
		return
	}

	if len(rui.Password) < 6 {
		panic("The password length is too short, greater than or equal 6")
		return
	}

	userInfo, err := api.RegistryCreateUser(rui.Username, rui.Password, rui.Email)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":   "register is success",
		"user_info": userInfo,
	})
}
