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

var skipLogPath = []string{"/health", "/swagger/index.html", "/swagger/swagger-ui.css",
	"/swagger/swagger-ui-standalone-preset.js", "/swagger/swagger-ui-bundle.js", "/swagger/swagger-ui.css.map",
	"/swagger/doc.json", "/swagger/swagger-ui-standalone-preset.js.map", "/swagger/swagger-ui-bundle.js.map",
	"/swagger/favicon-32x32.png", "/swagger/favicon-16x16.png"}

var skipAuthPath = []string{"/health", "/v1/auth/login", "/v1/auth/reset", "/v1/auth/pwd", "/swagger/index.html", "/swagger/swagger-ui.css",
	"/swagger/swagger-ui-standalone-preset.js", "/swagger/swagger-ui-bundle.js", "/swagger/swagger-ui.css.map",
	"/swagger/doc.json", "/swagger/swagger-ui-standalone-preset.js.map", "/swagger/swagger-ui-bundle.js.map",
	"/swagger/favicon-32x32.png", "/swagger/favicon-16x16.png"}

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
		middleware.AccessLog(skipLogPath...),
		middleware.Auth(skipAuthPath...),
		middleware.ErrorHandler(),
	)
	s.InitRouters()     // 初始化路由(定义所有的url)
	s.DBEngine.InitDB() // 同步库表
	s.initDBData()      // 初始化admin用户、role角色(管理员、开发者)
}

// use middleware
func (s *Server) addMiddleWare(mds ...gin.HandlerFunc) {
	s.RouterEngine.Use(mds...)
}

func (s *Server) initDBData() {
	e := database.GetDBEngine()
	roles := GetRolesInitData()
	existOrInsert(e, roles)
	records := GetUsersInitData()
	existOrInsert(e, records)
}
