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
		userApi := v1Root.Group("/users")
		{
			userApi.POST("", middleware.IsAdmin, ctlv1.CreateUser)
			userApi.GET("", middleware.IsAdmin, ctlv1.GetUsers)
			userApi.GET("/:id", middleware.IsUserSelfOrAdmin, ctlv1.GetUser)
			userApi.DELETE("/:id", middleware.IsAdmin, ctlv1.DeleteUser)
			userApi.PUT("/:id/access", middleware.IsUserSelfOrAdmin, ctlv1.UpdateUserAccessToken) // 关联token按钮,更新access token
			userApi.PUT("/:id/password", middleware.IsUserSelfOrAdmin, ctlv1.UpdateUserPwd)
			//userApi.PUT("/:id/roles", middleware.IsAdmin, ctlv1.UpdateUserRole)
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
			projectApi.DELETE("/:id", middleware.IsAdmin, ctlv1.DeleteProject)
			projectApi.PUT("/:id", middleware.IsAdmin, ctlv1.UpdateProject)
			projectApi.GET("/:id/apps", ctlv1.GetProjectApps)
			//projectApi.POST("/:id/project-envs", middleware.IsAdmin, ctlv1.CreateProjectEnv)
			//projectApi.GET("/:id/project-envs", ctlv1.GetProjectEnvs)
			//projectApi.DELETE("/:id/project-envs/:pe_id", middleware.IsAdmin, ctlv1.DeleteProjectEnv)
			//projectApi.GET("/:id/project-envs/:pe_id", ctlv1.GetProjectEnv)
			//projectApi.PUT("/:id/project-envs/:pe_id", middleware.IsAdmin, ctlv1.UpdateProjectEnv)
		}

		appApi := v1Root.Group("/apps")
		{
			appApi.POST("", ctlv1.CreateApp)
			appApi.GET("", ctlv1.GetApps)
			appApi.GET("/:id", ctlv1.GetApp)
			appApi.DELETE("/:id", ctlv1.DeleteApp)
			appApi.PUT("/:id", ctlv1.UpdateApp)
			//appApi.PUT("/:id/gitlab", ctlv1.UpdateAppGitlabProject)
			//appApi.GET("/:id/builds", ctlv1.GetAppBuilds)
			//appApi.POST("/:id/builds", ctlv1.CreateAppBuild)
			//appApi.GET("/:id/branches", ctlv1.GetAppBranches)
			//appApi.GET("/:id/tags", ctlv1.GetAppTags)
			//appApi.GET("/:id/charts", ctlv1.GetAppChartVersions)
			//appApi.GET("/:id/instances", ctlv1.GetAppInstances)
			//appApi.GET("/:id/builds/:build_number", ctlv1.GetAppBuild)
			//appApi.GET("/:id/builds/:build_number/logs/:log_number", ctlv1.GetBuildLogs)
			//appApi.DELETE("/:id/config", ctlv1.DeleteAppConf)
			//appApi.GET("/:id/config", ctlv1.GetAppConf)
			//appApi.PUT("/:id/config", ctlv1.UpdateOrCreateAppConf)
		}

		clusterApi := v1Root.Group("/clusters")
		{
			clusterApi.POST("", middleware.IsAdmin, ctlv1.CreateCluster)
			clusterApi.GET("", ctlv1.GetClusters)
			clusterApi.GET("/:id", ctlv1.GetCluster)
			clusterApi.DELETE("/:id", middleware.IsAdmin, ctlv1.DeleteCluster)
			clusterApi.PUT("/:id", middleware.IsAdmin, ctlv1.UpdateCluster)
			//clusterApi.GET("/:id/envs", ctlv1.GetClusterEnvs)
			//clusterApi.GET("/:id/nodes", ctlv1.GetClusterNodes) // 获取节点信息
			//clusterApi.GET("/:id/nodes/:name", ctlv1.GetClusterNode)
			//clusterApi.GET("/:id/jobs", ctlv1.GetClusterJobs)
			//clusterApi.GET("/:id/namespaces/:namespace/jobs/:name", ctlv1.GetClusterJob)
			//clusterApi.GET("/:id/license-status", middleware.IsSystemAdminOrAdmin, ctlv1.GetLicenseStatus)
			//clusterApi.GET("/:id/license-c2v", middleware.IsSystemAdminOrAdmin, ctlv1.GetC2vFile)
			//clusterApi.GET("/:id/license-fingerprint", middleware.IsSystemAdminOrAdmin, ctlv1.GetFingerprintFile)
			//clusterApi.POST("/:id/license-online", middleware.IsSystemAdminOrAdmin, ctlv1.ActiveOnline)
			//clusterApi.POST("/:id/license-offline", middleware.IsSystemAdminOrAdmin, ctlv1.ActiveOffline)
			//clusterApi.GET("/:id/license-clics", middleware.IsSystemAdminOrAdmin, ctlv1.GetClientLicenses)
		}

		envApi := v1Root.Group("/envs")
		{
			envApi.POST("", ctlv1.CreateEnv)
			envApi.GET("", ctlv1.GetEnvs)
			envApi.GET("/:id", ctlv1.GetEnv)
			envApi.DELETE("/:id", ctlv1.DeleteEnv)
			envApi.PUT("/:id", ctlv1.UpdateEnv)
			//envApi.POST("/:id/jobs", ctlv1.CreateEnvJob)
			//envApi.POST("/:id/cronjobs", ctlv1.CreateEnvCronJob)
		}

		authApi := v1Root.Group("/auth")
		{
			authApi.POST("/login", ctlv1.Login)
			authApi.POST("/logout", ctlv1.Logout)
			authApi.POST("/reset", ctlv1.CreateResetEmail)
			authApi.PUT("/pwd", ctlv1.UpdateUserPwdWithSecret)
		}

		repoApi := v1Root.Group("/repos")
		{
			repoApi.GET("", ctlv1.GetRemoteRepos)
		}

	}

	// 健康检查接口
	router.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})
}
