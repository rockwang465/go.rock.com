package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/conf"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
	"time"
)

type LoginUserInfo struct {
	Name     string `json:"name" binding:"required" example:"admin_user"`
	Password string `json:"password" binding:"required" example:"********"`
}

type ResetUserReq struct {
	Email string `json:"email" binding:"required,email" example:"admin_user@sensetime.com"`
}

type ResetPwdReq struct {
	Secret     string `json:"secret" binding:"required" example:"TR6UdhT7ebJOCC5N"`
	Email      string `json:"email" binding:"required,email" example:"admin_user@sensetime.com"`
	Password   string `json:"password" binding:"required" example:"********"`
	RePassword string `json:"re_password" binding:"required" example:"********"`
}

// @Summary Login rock platform with name and password
// @Description Api to login rock platform with name and password
// @Tags AUTH
// @Accept json
// @Produce json
// @Param input_body body v1.LoginUserInfo true "JSON type input body"
// @Success 201 {object} api.UserDetailResp
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/auth/login [post]
func (c *Controller) Login(ctx *gin.Context) {
	var userInfo LoginUserInfo
	if err := ctx.ShouldBind(&userInfo); err != nil {
		panic(err)
		return
	}

	// check user is exist
	user, err := api.GetUserDetailRespByName(userInfo.Name)
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

	// generate admin account token
	//tokenAdmin, err := middleware.GenerateToken(1, "admin", "3207ead4e092de77e022394b3204d755")
	//if err != nil {
	//	panic(err)
	//}else{
	//	fmt.Println("admin user token")
	//	fmt.Println(tokenAdmin)
	//}

	// generate jwt token
	token, err := utils.GenerateToken(user.Id, user.Name, user.Password, user.RoleName)
	if err != nil {
		panic(err)
		return
	}

	// update token to mysql
	var userData *models.User
	if user.Token != token {
		userData, err = api.UpdateUserToken(user.Id, token)
		if err != nil {
			panic(err)
			return
		}
	}

	// return user info
	resp, err := api.GetUserBriefResp(user.Id)
	if err != nil {
		panic(err)
		return
	}

	// save cookie
	ctx.SetCookie("token", userData.Token, utils.GetExpireDuration(), "/", "", false, true)
	c.Logger.Infof("User %v login successful", user.Name)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Logout rock platform
// @Description Api to logout rock platform
// @Tags AUTH
// @Accept json
// @Produce json
// @Success 204 {object} string "StatusNoContent"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/auth/logout [post]
func (c *Controller) Logout(ctx *gin.Context) {
	config, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	ctx.SetCookie("token", "", -1, "/", "", false, true)
	c.Logger.Infof("User %v logout successful", config.Username)
	ctx.JSON(http.StatusNoContent, "")
}

// @Summary Create reset email
// @Description Api to create reset email
// @Tags AUTH
// @Accept json
// @Produce json
// @Param input_body body v1.ResetUserReq true "JSON type input body"
// @Success 204 {object} string "StatusNoContent"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/auth/reset [post]
func (c *Controller) CreateResetEmail(ctx *gin.Context) {
	var resetUserReq ResetUserReq
	if err := ctx.ShouldBindJSON(&resetUserReq); err != nil {
		panic(err)
	}

	user, err := api.GetUserByEmail(resetUserReq.Email)
	if err != nil {
		panic(err)
	}

	secret := utils.GenerateSalt()
	config := conf.GetConfig()
	secretExpire := config.Viper.GetDuration("email.secret-expire")
	user, err = api.ResetSecretWithId(user.Id, secret, secretExpire)
	if err != nil {
		panic(err)
	}

	if err = utils.SendResetPwdEmail(user.Name, user.Email, secret, secretExpire); err != nil {
		panic(err)
	}
	c.Debugf("User(%s) send reset email successful", user.Name)

	ctx.SetCookie("token", "", -1, "/", "", false, true)

	c.Logger.Infof("User %v reset password by email %v successful", user.Name, resetUserReq.Email)
	ctx.JSON(http.StatusNoContent, "")
}

// @Summary Update user's password with secret
// @Description Api to update user password with secret
// @Tags AUTH
// @Accept json
// @Produce json
// @Param input_body body v1.ResetPwdReq true "JSON type input body"
// @Success 200 {object} string "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/auth/pwd [post]
func (c *Controller) UpdateUserPwdWithSecret(ctx *gin.Context) {
	var pwdReq ResetPwdReq
	err := ctx.ShouldBind(&pwdReq)
	if err != nil {
		panic(err)
	}

	if pwdReq.Password != pwdReq.RePassword {
		err = utils.NewRockError(400, 40000010, "Two input password is not the same")
		panic(err)
	}

	// check email
	has, err := api.HasEmail(pwdReq.Email)
	if err != nil {
		panic(err)
	}
	if !has {
		err = utils.NewRockError(400, 40000011, fmt.Sprintf("Email %v not found", pwdReq.Email))
		panic(err)
	}

	// check password is strong
	err = utils.CheckPwd(pwdReq.Password)
	if err != nil {
		panic(err)
	}

	user, err := api.UpdateUserPwdBySecret(pwdReq.Password, pwdReq.Email, pwdReq.Secret)
	if err != nil {
		panic(err)
	}

	userDetail, err := api.GetUserBriefResp(user.Id)
	if err != nil {
		panic(err)
	}
	token, err := utils.GenerateToken(user.Id, user.Name, user.Password, userDetail.RoleName)
	user, err = api.UpdateUserToken(user.Id, token)
	if err != nil {
		panic(err)
	}

	c.Logger.Info(fmt.Sprintf("User id:%v name:%v password reset successful", user.Id, user.Name))
	ctx.JSON(http.StatusOK, fmt.Sprintf("User %v password reset successful", user.Name))
}
