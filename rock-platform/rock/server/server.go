package server

import (
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/database"
	"go.rock.com/rock-platform/rock/server/log"
	middleware "go.rock.com/rock-platform/rock/server/middleware"
	"go.rock.com/rock-platform/rock/server/routerEngine"
)

type Server struct {
	Logger       *log.Logger           `json:"logger"`
	RouterEngine *routerEngine.Routers `json:"router_engine"`
	DBEngine     *database.DBEngine    `json:"db_engine"`
}

var SingleServer *Server

func GetServer() *Server {
	if SingleServer == nil {
		SingleServer = &Server{
			Logger:       log.GetLogger(),                // 实例化logrus.Logger对象
			RouterEngine: routerEngine.GetRouterEngine(), // 实例化一个没有中间件的空白路由(r := gin.New()代替r := gin.Default())
			DBEngine:     database.GetDBEngine(),         // 实例化gorm数据库
		}
	}
	return SingleServer
}

// 初始化日志配置、中间件、路由、数据库、validator(不明白)
func (s *Server) InitServer() {
	s.Logger.InitLogger() // 初始化日志配置(日志级别、日志文件、日志分割、日志格式)

	s.addMiddleWare(
		middleware.ErrorHandler(),
	)

	s.InitRouters()     // 初始化路由(定义所有的url)
	s.DBEngine.InitDB() // 同步库表
}

// use middleware
func (s *Server) addMiddleWare(mds ...gin.HandlerFunc) {
	s.RouterEngine.Use(mds...)
}
