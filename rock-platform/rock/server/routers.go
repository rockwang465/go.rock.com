package server

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.rock.com/rock-platform/rock/server/controller/v1"
	middleware "go.rock.com/rock-platform/rock/server/middleware"
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
		registryApi := v1Root.Group("/users")
		{
			registryApi.POST("", middleware.IsAdmin, ctlv1.CreateUser)
			registryApi.GET("", middleware.IsAdmin, ctlv1.GetUsers)
			registryApi.GET("/:id", middleware.IsUserSelfOrAdmin, ctlv1.GetUser)
			registryApi.DELETE("/:id", middleware.IsAdmin, ctlv1.DeleteUser)
			registryApi.PUT("/:id/password", middleware.IsUserSelfOrAdmin, ctlv1.UpdateUserPwd)
		}

		roleApi := v1Root.Group("/roles")
		{
			roleApi.POST("", middleware.IsAdmin, ctlv1.CreateRole)
			roleApi.GET("", middleware.IsAdmin, ctlv1.GetRoles)
			roleApi.GET("/:id", middleware.IsAdmin, ctlv1.GetRole)
			roleApi.DELETE("/:id", middleware.IsAdmin, ctlv1.DeleteRole)
			roleApi.PUT("/:id", middleware.IsAdmin, ctlv1.UpdateRole)
			roleApi.GET("/:id/users", middleware.IsAdmin, ctlv1.GetRoleUsers)
			//roleApi.GET("/:id/permissions", middleware.IsAdmin, ctlv1.GetRolePermissions)
			//roleApi.PUT("/:id/permissions", middleware.IsAdmin, ctlv1.ManagerRolePermissions)
		}

		//permApi := v1Root.Group("/permissions")
		//{
		//	permApi.POST("", ctlv1.CreatePermission)
		//	permApi.GET("", ctlv1.GetPermissions)
		//	permApi.GET("/:id", ctlv1.GetPermission)
		//	permApi.DELETE("/:id", ctlv1.DeletePermission)
		//	permApi.PUT("/:id", ctlv1.UpdatePermission)
		//	permApi.GET("/:id/roles", ctlv1.GetPermissionRoles)
		//}

		projectApi := v1Root.Group("/projects")
		{
			projectApi.POST("", middleware.IsAdmin, ctlv1.CreateProject)
			projectApi.GET("", ctlv1.GetProjects)
			projectApi.GET("/:id", ctlv1.GetProject)
			//projectApi.DELETE("/:id", middleware.IsAdmin, ctlv1.DeleteProject)
			//projectApi.PUT("/:id", middleware.IsAdmin, ctlv1.UpdateProject)
			//projectApi.GET("/:id/apps", ctlv1.GetProjectApps)
			//projectApi.POST("/:id/project-envs", middleware.IsAdmin, ctlv1.CreateProjectEnv)
			//projectApi.GET("/:id/project-envs", ctlv1.GetProjectEnvs)
			//projectApi.DELETE("/:id/project-envs/:pe_id", middleware.IsAdmin, ctlv1.DeleteProjectEnv)
			//projectApi.GET("/:id/project-envs/:pe_id", ctlv1.GetProjectEnv)
			//projectApi.PUT("/:id/project-envs/:pe_id", middleware.IsAdmin, ctlv1.UpdateProjectEnv)
		}

		authApi := v1Root.Group("/auth")
		{
			authApi.POST("/login", ctlv1.Login)
			authApi.POST("/logout", ctlv1.Logout)
			authApi.POST("/reset", ctlv1.CreateResetEmail)
			authApi.POST("/pwd", ctlv1.UpdateUserPwdWithSecret)
		}

	}

	// 健康检查接口
	router.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})
}
