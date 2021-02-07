package middleWare

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/rockwang465/drone/drone-go/drone"
	"go.rock.com/rock-platform/rock/server/log"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
)

// 返回panic的报错给用户界面，而不是直接被panic了。需要在路由中Use。
func ErrorHandler() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		logger := log.GetLogger()

		defer func() {
			if err := recover(); err != nil {
				//response.Response(ctx, 500, 500001, fmt.Sprint(err), nil)
				switch e := err.(type) {
				// RockError 为自定义的错误格式，所有严重报错部分都可以用到这里的格式，然后panic出来，就会到recover这里了。
				case *utils.RockError:
					logger.Errorf("Rock Error: %v", e)
					ctx.JSON(e.HttpCode, gin.H{"error": e.Error(), "error_code": e.ErrCode})
					ctx.Abort()
					return
				case *mysql.MySQLError:
					logger.Errorf("Mysql Error: with num %v and message is: %v", e.Number, e.Message)
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error":      fmt.Sprintf("Mysql Error: with num %v and message is: %v", e.Number, e.Message),
						"error_code": 50000002,
					})
					ctx.Abort()
					return
				case *drone.DroneError:
					// 如果error内无明显的报错信息，可以去drone-server服务的日志里查看
					logger.Errorf("Drone Error: %v", e)
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error":      e.Error(),
						"error_code": e.ErrCode,
					})
					ctx.Abort()
					return
					// case *k8sErr.StatusError:
					// validator.ValidationErrors:
				default:
					logger.Error(err)
					// 这里要把default出来的类型打印一下，不然rock不确定为什么会走到default和type是什么也不知道。
					logger.Infof("type is :[%T]", e)
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err, "error_code": 5000001, "switch": "default"})
					ctx.Abort()
					return
				}
			}
		}()
		ctx.Next()
	}
}
