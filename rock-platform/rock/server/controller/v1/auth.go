package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/database/api"
	middleware "go.rock.com/rock-platform/rock/server/middleware"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
	"time"
)

// 对用户是否有权限操作进行验证的模块

type LoginUserInfo struct {
	Username string `json:"username" binding:"required" example:"admin_user"`
	Password string `json:"password" binding:"required" example:"********"`
}

// @Summary Login
// @Description Api to login
// @Tags USER
// @Accept  json
// @Produce  json
// @Param input_body body v1.LoginUserInfo true  "JSON type input body"
// @Success 201 {object} v1.UserDetailResp
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/login [post]
func (c *Controller) Login(ctx *gin.Context) {
	var userInfo LoginUserInfo
	if err := ctx.ShouldBind(&userInfo); err != nil {
		panic(err)
		return
	}

	// check user is exist
	user, err := api.GetUserByName(userInfo.Username)
	if err != nil {
		panic(err)
		return
	}

	// 1. if block time: return error
	// 2. if password is wrong:
	//     retry count + 1
	//     2.1 if user retry count > max retry login count (viper.GetInt64(server.login-retry-count)):
	//           increase block time(default increase 5 minutes)
	// 3. if password is right:
	//     3.1 reset retry count
	//     3.2 generate token

	// block from multiple login fail user
	if user.LoginBlockUntil != nil && time.Now().Before(*user.LoginBlockUntil) {
		err := utils.NewRockError(401, 40100004, fmt.Sprintf("User %v failed login too many times, block to %v", user.Name, user.LoginBlockUntil))
		panic(err)
		return
	}

	// verify password
	encryptPwd := utils.EncryptPwd(userInfo.Password, user.Salt)
	if encryptPwd != user.Password {
		if err := api.CountUserLoginFailedNumber(user.Id); err != nil { // increase the number of user failed login count and time
			panic(err)
			return
		}
		err := utils.NewRockError(400, 40000003, "password incorrect")
		panic(err)
		return
	} else {
		err := api.ResetRetryCount(user.Id)
		if err != nil {
			panic(err)
		}
	}

	// generate jwt token
	token, err := middleware.GenerateToken(user.Id, user.Name, user.Password)
	if err != nil {
		panic(err)
		return
	}

	// update token to mysql
	if user.Token != token {
		user, err = api.UpdateUserToken(user.Id, token)
		if err != nil {
			panic(err)
			return
		}
	}

	// save cookie
	ctx.SetCookie("token", user.Token, utils.GetExpireDuration(), "/", "", false, true)
	c.Logger.Infof("User %v login successful", user.Name)

	resp := &UserDetailResp{
		Id:        user.Id,
		Username:  user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	ctx.JSON(http.StatusOK, gin.H{"user_info": resp})
}
