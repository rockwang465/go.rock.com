package middleWare

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/log"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
	"strings"
)

const (
	AdminRole            string = "admin"
	SystemToolsAdminRole        = "system_tools_admin"
	DeployerRole                = "deployer"
	DeveloperRole               = "developer"
)

//func Auth(skipPath []string, ctx *gin.Context) {
func Auth(skipPath ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var skipUrl map[string]bool
		if length := len(skipPath); length > 0 {
			skipUrl = make(map[string]bool, length)
			for _, path := range skipPath {
				skipUrl[path] = true
			}
		}

		reqUrl := ctx.Request.URL.Path
		_, ok := skipUrl[reqUrl]
		if !ok { // if not in skip auth path, then need to auth
			logger := log.GetLogger()
			token, _ := ctx.Cookie("token")
			if token == "" {
				auth := ctx.GetHeader("Authorization")
				if !strings.HasPrefix(auth, "Bearer ") {
					ctx.JSON(http.StatusUnauthorized, "Permission denied, user doesn't login, please login first ")
					logger.Warning("Permission denied, not Bearer int prefix")
					ctx.Abort()
					return
				}
				token = strings.TrimPrefix(auth, "Bearer ")
			}

			// parse token
			jwtToken, claim, err := utils.ParseToken(token)
			if err != nil {
				newErr := fmt.Sprintf("Parse token failed , %v", err)
				ctx.JSON(http.StatusUnauthorized, newErr)
				logger.Error(newErr)
				ctx.Abort()
				return
			}
			if !jwtToken.Valid {
				newErr := fmt.Sprintf("Token is valid")
				ctx.JSON(http.StatusUnauthorized, newErr)
				logger.Warn(newErr)
				ctx.Abort()
				return
			}

			// query user info
			user, err := api.HasUserById(claim.UserId)
			if err != nil {
				newErr := fmt.Sprintf("User token is not valid, %v", err)
				ctx.JSON(http.StatusNotFound, "User token is not valid, please login later")
				logger.Warn(newErr)
				ctx.Abort()
				return
			}

			// compare token
			if user.Token != token {
				newErr := "User info is malformed, please login again"
				ctx.JSON(http.StatusUnauthorized, newErr)
				logger.Warn(newErr)
				ctx.SetCookie("token", "", -1, "/", "", false, true)
				ctx.Abort()
				return
			}

			// authorize success, set cookie
			var cfgCtx = models.ConfCtx{
				UserId:   claim.UserId,
				Username: claim.Username,
				Role:     claim.Role,
			}
			ctx.Set("custom_config", cfgCtx)
			ctx.SetCookie("token", user.Token, utils.GetExpireDuration(), "/", "", false, true)

			ctx.Next()
		}
	}
}

// check context, must be admin role
func IsAdmin(ctx *gin.Context) {
	config, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	// verify claim by id
	if config.Role != AdminRole && config.Role != SystemToolsAdminRole {
		newWarn := fmt.Sprintf("Permission denied, only system_tools_admin or admin role can do this operation")
		err := utils.NewRockError(http.StatusUnauthorized, 40100001, newWarn)
		panic(err)
	}
}

// check context, is user self or admin role
func IsUserSelfOrAdmin(ctx *gin.Context) {
	config, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	// check is user self
	user, err := api.GetUserDetailResp(config.UserId)
	if err != nil {
		panic(err)
	}
	if user.Id == config.UserId {
		return
	}

	// else must be admin role
	if user.RoleName == AdminRole {
		return
	} else {
		newWarn := fmt.Sprintf("Permission denied, only user self or admin role can do this operation")
		err := utils.NewRockError(http.StatusUnauthorized, 40100001, newWarn)
		panic(err)
	}
}
