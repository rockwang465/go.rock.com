package server

import (
	"rock-platform/rock/server/log"
	"rock-platform/rock/server/routerEngine"
)

type Server struct {
	Logger       *log.Logger           `json:"logger"`
	RouterEngine *routerEngine.Routers `json:"router_engine"`
	//DBEngine
}

func GetServer() *Server {
	sv := &Server{
		Logger:       log.GetLogger(),
		RouterEngine: routerEngine.GetRouterEngine(),
	}
	return sv
}

// 初始化日志配置、中间件、路由、数据库、validator(不明白)
func (s *Server) InitServer() {
	s.Logger.InitLogger() // 初始化日志配置(日志级别、日志文件、日志分割、日志格式)

	//s.RouterEngine.Use()

	s.InitRouters() // 初始化路由(定义所有的url)
}
