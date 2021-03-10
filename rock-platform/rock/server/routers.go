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
		userApi.POST("", middleware.IsAdmin, ctlv1.CreateUser)
		userApi.GET("", middleware.IsAdmin, ctlv1.GetUsers)
		userApi.GET("/:id", middleware.IsUserSelfOrAdmin, ctlv1.GetUser)
		userApi.DELETE("/:id", middleware.IsAdmin, ctlv1.DeleteUser)
		userApi.PUT("/:id/access", middleware.IsUserSelfOrAdmin, ctlv1.UpdateUserAccessToken) // 关联token按钮,更新access token
		userApi.PUT("/:id/password", middleware.IsUserSelfOrAdmin, ctlv1.UpdateUserPwd)
		userApi.PUT("/:id/roles", middleware.IsAdmin, ctlv1.UpdateUserRole)

		roleApi := v1Root.Group("/roles")
		roleApi.POST("", middleware.IsAdmin, ctlv1.CreateRole)
		roleApi.GET("", middleware.IsAdmin, ctlv1.GetRoles)
		roleApi.GET("/:id", middleware.IsAdmin, ctlv1.GetRole)
		roleApi.DELETE("/:id", middleware.IsAdmin, ctlv1.DeleteRole)
		roleApi.PUT("/:id", middleware.IsAdmin, ctlv1.UpdateRole)
		roleApi.GET("/:id/users", middleware.IsAdmin, ctlv1.GetRoleUsers)
		//roleApi.GET("/:id/permissions", middleware.IsAdmin, ctlv1.GetRolePermissions)
		//roleApi.PUT("/:id/permissions", middleware.IsAdmin, ctlv1.ManagerRolePermissions)

		//permApi := v1Root.Group("/permissions")
		//permApi.POST("", ctlv1.CreatePermission)
		//permApi.GET("", ctlv1.GetPermissions)
		//permApi.GET("/:id", ctlv1.GetPermission)
		//permApi.DELETE("/:id", ctlv1.DeletePermission)
		//permApi.PUT("/:id", ctlv1.UpdatePermission)
		//permApi.GET("/:id/roles", ctlv1.GetPermissionRoles)

		projectApi := v1Root.Group("/projects")
		projectApi.POST("", middleware.IsAdmin, ctlv1.CreateProject)
		projectApi.GET("", ctlv1.GetProjects)
		projectApi.GET("/:id", ctlv1.GetProject)
		projectApi.DELETE("/:id", middleware.IsAdmin, ctlv1.DeleteProject)
		projectApi.PUT("/:id", middleware.IsAdmin, ctlv1.UpdateProject)
		projectApi.GET("/:id/apps", ctlv1.GetProjectApps)
		projectApi.POST("/:id/project-envs", middleware.IsAdmin, ctlv1.CreateProjectEnv) // 对指定project_id增加项目环境(project_env表)
		projectApi.GET("/:id/project-envs", ctlv1.GetProjectEnvs)
		//projectApi.DELETE("/:id/project-envs/:pe_id", middleware.IsAdmin, ctlv1.DeleteProjectEnv) // 暂时无法继续写，缺少app_conf表 // 删除指定project(project_id)下指定项目环境(project_env_id)
		projectApi.GET("/:id/project-envs/:pe_id", ctlv1.GetProjectEnv)                        // 查看指定project(project_id)下指定项目环境(project_env_id)的信息
		projectApi.PUT("/:id/project-envs/:pe_id", middleware.IsAdmin, ctlv1.UpdateProjectEnv) // 更新指定project(project_id)下指定项目环境(project_env_id)的信息

		buildApi := v1Root.Group("/builds")
		buildApi.GET("", ctlv1.GetGlobalBuilds) // console_project_id未发现使用场景? 获取所有的构建任务的简单记录(用于页面批量展示)

		nodeApi := v1Root.Group("/nodes")
		nodeApi.GET("", ctlv1.GetGlobalNodes)

		//jobApi := v1Root.Group("/jobs")  // 暂时弃用
		//jobApi.GET("", ctlv1.GetGlobalJobs)

		appApi := v1Root.Group("/apps")
		appApi.POST("", ctlv1.CreateApp)
		appApi.GET("", ctlv1.GetApps)
		appApi.GET("/:id", ctlv1.GetApp)
		appApi.DELETE("/:id", ctlv1.DeleteApp)
		appApi.PUT("/:id", ctlv1.UpdateApp)
		appApi.PUT("/:id/gitlab", ctlv1.UpdateAppGitlabProject) // 修改应用的gitlab地址
		appApi.GET("/:id/builds", ctlv1.GetAppBuilds)           // 通过app_id获取该应用的所有构建记录
		appApi.POST("/:id/builds", ctlv1.CreateAppBuild)        // 构建一个新任务(单个服务发版)
		appApi.GET("/:id/branches", ctlv1.GetAppBranches)       // 获取应用所有分支
		appApi.GET("/:id/tags", ctlv1.GetAppTags)               // 获取应用所有tag
		appApi.GET("/:id/charts", ctlv1.GetAppChartVersions)    // 通过app_id获取该应用的所有chart版本
		//appApi.GET("/:id/instances", ctlv1.GetAppInstances)
		appApi.GET("/:id/builds/:build_number", ctlv1.GetAppBuild)                   // 通过app_id和build_number查看指定任务的详细构建记录
		appApi.GET("/:id/builds/:build_number/logs/:log_number", ctlv1.GetBuildLogs) // 通过app_id和build_number和log_number查看指定任务的详细构建日志
		appApi.DELETE("/:id/config", ctlv1.DeleteAppConf)
		appApi.GET("/:id/config", ctlv1.GetAppConf)
		appApi.PUT("/:id/config", ctlv1.UpdateOrCreateAppConf) // 更新/创建一个应用的配置(values.yaml)

		chartApi := v1Root.Group("/charts")
		chartApi.GET("", ctlv1.GetAllCharts)
		chartApi.GET("/:name", ctlv1.GetNamedChartVersions)                        // 获取指定服务的所有版本号
		chartApi.GET("/:name/versions/:version", ctlv1.GetNamedChartVersion)       // 获取指定服务的指定版本号
		chartApi.DELETE("/:name/versions/:version", ctlv1.DeleteNamedChartVersion) // 删除指定服务的指定版本号

		envApi := v1Root.Group("/envs")
		envApi.POST("", ctlv1.CreateEnv)
		envApi.GET("", ctlv1.GetEnvs)
		envApi.GET("/:id", ctlv1.GetEnv)
		envApi.DELETE("/:id", ctlv1.DeleteEnv)
		envApi.PUT("/:id", ctlv1.UpdateEnv)
		//envApi.POST("/:id/jobs", ctlv1.CreateEnvJob)
		//envApi.POST("/:id/cronjobs", ctlv1.CreateEnvCronJob)

		clusterApi := v1Root.Group("/clusters")
		clusterApi.POST("", middleware.IsAdmin, ctlv1.CreateCluster)
		clusterApi.GET("", ctlv1.GetClusters)
		clusterApi.GET("/:id", ctlv1.GetCluster)
		clusterApi.DELETE("/:id", middleware.IsAdmin, ctlv1.DeleteCluster)
		clusterApi.PUT("/:id", middleware.IsAdmin, ctlv1.UpdateCluster)
		clusterApi.GET("/:id/envs", ctlv1.GetClusterEnvs)        // 获取指定cluster_id的env信息
		clusterApi.GET("/:id/nodes", ctlv1.GetClusterNodes)      // 获取指定cluster_id的节点信息
		clusterApi.GET("/:id/nodes/:name", ctlv1.GetClusterNode) // 获取指定节点名称的节点信息
		//clusterApi.GET("/:id/jobs", ctlv1.GetClusterJobs)
		//clusterApi.GET("/:id/namespaces/:namespace/jobs/:name", ctlv1.GetClusterJob)
		//clusterApi.GET("/:id/license-status", middleware.IsSystemAdminOrAdmin, ctlv1.GetLicenseStatus)
		//clusterApi.GET("/:id/license-c2v", middleware.IsSystemAdminOrAdmin, ctlv1.GetC2vFile)
		//clusterApi.GET("/:id/license-fingerprint", middleware.IsSystemAdminOrAdmin, ctlv1.GetFingerprintFile)
		//clusterApi.POST("/:id/license-online", middleware.IsSystemAdminOrAdmin, ctlv1.ActiveOnline)
		//clusterApi.POST("/:id/license-offline", middleware.IsSystemAdminOrAdmin, ctlv1.ActiveOffline)
		//clusterApi.GET("/:id/license-clics", middleware.IsSystemAdminOrAdmin, ctlv1.GetClientLicenses)

		registryApi := v1Root.Group("/registries") // 镜像仓库信息管理
		registryApi.POST("", middleware.IsAdmin, ctlv1.CreateRegistry)
		registryApi.GET("", middleware.IsAdmin, ctlv1.GetRegistries)
		registryApi.GET("/:address", middleware.IsAdmin, ctlv1.GetRegistry)
		registryApi.DELETE("/:address", middleware.IsAdmin, ctlv1.DeleteRegistry)
		registryApi.PUT("/:address", middleware.IsAdmin, ctlv1.UpdateRegistry)

		secretApi := v1Root.Group("/secrets") // 密钥管理
		secretApi.POST("", middleware.IsAdmin, ctlv1.CreateSecret)
		secretApi.GET("", middleware.IsAdmin, ctlv1.GetSecrets)
		secretApi.GET("/:name", middleware.IsAdmin, ctlv1.GetSecret)
		secretApi.DELETE("/:name", middleware.IsAdmin, ctlv1.DeleteSecret)
		secretApi.PUT("/:name", middleware.IsAdmin, ctlv1.UpdateSecret)

		deployApi := v1Root.Group("/deployments")
		deployApi.POST("", ctlv1.CreateDeployment) // 基于固定版本号进行单个应用发版
		//deployApi.GET("", ctlv1.GetDeployments)
		//deployApi.GET("/:id", ctlv1.GetDeployment)
		//deployApi.DELETE("/:id", ctlv1.DeleteDeployment)
		//deployApi.PUT("/:id", ctlv1.UpdateDeployment)

		//instanceApi := v1Root.Group("/instances")
		//instanceApi.GET("", ctlv1.GetInstances)
		//instanceApi.GET("/:id", ctlv1.GetInstance)
		//instanceApi.GET("/:id/configs", ctlv1.GetInstanceConfig)
		//instanceApi.GET("/:id/deployments", ctlv1.GetInstanceDeployment)
		//instanceApi.GET("/:id/logs", ctlv1.GetInstanceLog)
		//instanceApi.GET("/:id/logfile", ctlv1.GetInstanceLogFile)
		//instanceApi.GET("/:id/pods", ctlv1.GetInstancePods)
		//instanceApi.GET("/:id/scale", ctlv1.GetInstanceScale)
		//instanceApi.PUT("/:id/scale", ctlv1.UpdateInstanceScale)
		//instanceApi.DELETE("/:id", ctlv1.DeleteInstance)

		authApi := v1Root.Group("/auth")
		authApi.POST("/login", ctlv1.Login)
		authApi.POST("/logout", ctlv1.Logout)
		authApi.POST("/reset", ctlv1.CreateResetEmail)
		authApi.PUT("/pwd", ctlv1.UpdateUserPwdWithSecret)

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
