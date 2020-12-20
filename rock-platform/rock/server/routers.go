package server

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.rock.com/rock-platform/rock/server/controller/v1"
	"net/http"
)

// 所有路由api都定义在这里
func (s *Server) InitRouters() {
	router := s.RouterEngine
	ctlv1 := v1.GetController() // 先自动加载log模块，且所有方法写在Controller结构体下

	// use ginSwagger middleware to
	// Rock找的文档: https://juejin.im/post/6844904198211895303
	// 官方文档: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
	// 将swagger添加到路由中：
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1Root := router.Group("/v1")
	{
		v1Root.POST("/register", ctlv1.Register)
	}

	// 健康检查接口
	router.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})
}
