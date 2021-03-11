package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/clients/drone-api"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
)

type CreateBuildReq struct {
	Name string `json:"name" binding:"required" example:"master/v1.0.0"`
	Type string `json:"type" binding:"required,oneof=branch tag" example:"branch/tag"`
	Envs []*Env `json:"envs,omitempty" binding:"omitempty"`
}

type Env struct {
	Key   string `json:"key" binding:"required" example:"key"`
	Value string `json:"value" binding:"required" example:"value"`
}

type BuildDetailResp struct {
	BuildBriefResp
	Procs []*RootProc `json:"procs" binding:"omitempty"`
}

type RootProc struct {
	ProcCommon
	Children []*ProcCommon `json:"children,omitempty" binding:"required"`
}

type ProcCommon struct {
	ID       int64  `json:"id" binding:"required" example:"1"`
	PID      int64  `json:"pid" binding:"required" example:"1"`
	PPID     int64  `json:"ppid" binding:"required" example:"1"`
	PGID     int64  `json:"pgid" binding:"required" example:"1"`
	Name     string `json:"name" binding:"required" example:"clone"`
	State    string `json:"state" binding:"required" example:"success/failure"`
	ExitCode int64  `json:"exit_code" binding:"required" example:"0"`
	Started  int64  `json:"start_time" binding:"required" example:"1542277389"`
	Stopped  int64  `json:"end_time" binding:"required" example:"1542277389"`
	Machine  string `json:"machine" binding:"required" example:"2b545f3330cf"`
}

type BuildBriefResp struct {
	Id               int64  `json:"id" example:"1"`
	RepoId           int64  `json:"repo_id" example:"1"`
	AppId            int64  `json:"app_id,omitempty" example:"1"` // extendBuildField func add this field
	ConsoleProjectId int64  `json:"console_project_id" example:"1"`
	Number           int64  `json:"number" example:"1"`
	Status           string `json:"status" example:"pending/success/failure"`
	Commit           string `json:"commit" example:"26b4808f0d35ac8f4621490166d683e255d9fed4"`
	Branch           string `json:"branch" example:"master"`
	Message          string `json:"message" example:"fix template issue\n"`
	Author           string `json:"author" example:"someone"`
	AuthorEmail      string `json:"author_email" example:"someone@sensetime.com"`
	CreatedAt        int64  `json:"created_at" example:"1614655426"`
	StartedAt        int64  `json:"started_at" example:"1614655426"`
	FinishedAt       int64  `json:"finished_at" example:"1614655426"`
	EnqueuedAt       int64  `json:"enqueued_at" example:"1614655426"`
	//Parent           int64  `json:"parent" example:"1"`
	//Event            string `json:"event" example:"deployment"`
	//Error            string `json:"error" example:""`
	//DeployTo         int64  `json:"deploy_to" example:""`
	//Ref              string `json:"ref" example:"refs/heads/master"`
	//Refspec          string `json:"refspec" example:""`
	//Remote           string `json:"remote" example:""`
	//Title            string `json:"title" example:""`
	//Timestamp        int64  `json:"timestamp" example:"0"`
	//Sender           string `json:"sender" example:""`
	//AuthorAvatar     string `json:"author_avatar" example:"https://www.gravatar.com/avatar/44f1af844f14d167aaa69014a5176353.jpg?s=128"`
	//LinkUrl          string `json:"link_url" example:"https://gitlab.sz.sensetime.com/galaxias/charts/senseguard-guest-management"`
	//ReviewedBy       string `json:"reviewed_by" example:""`
	//ReviewedAt       int64  `json:"reviewed_at" example:"0"`
}

type GetBuildsPaginationReq struct {
	PerPage    int64  `json:"per_page" form:"per_page" binding:"required,min=1" example:"10"`
	Page       int64  `json:"page" form:"page" binding:"required,min=1" example:"1"`
	QueryField string `json:"query_field" form:"query_field" binding:"omitempty" example:""`
}

type GetBuildPaginationReq struct {
	PerPage          int64 `json:"per_page" form:"per_page" binding:"required,min=1" example:"10"`                    // page_size
	Page             int64 `json:"page" form:"page" binding:"required,min=1" example:"1"`                             // page_num
	ConsoleProjectId int64 `json:"console_project_id" form:"console_project_id" binding:"omitempty,min=1" example:""` // console_project_id为project_id,当你需要过滤查某个项目下所有的构建信息时，可以通过下拉菜单选择某个项目下的所有构建信息
}

type PaginateBuildResp struct {
	Page    int64             `json:"page" binding:"required" example:"1"`      // page_num
	PerPage int64             `json:"per_page" binding:"required" example:"10"` // page_size
	Total   int64             `json:"total" binding:"required" example:"100"`
	Pages   int64             `json:"pages" binding:"required" example:"1"`
	Items   []*BuildBriefResp `json:"items" binding:"required"`
}

type GetAppBuildReq struct {
	Id          int64 `json:"id" uri:"id" binding:"required,min=1" example:"1"`                      // app id
	BuildNumber int   `json:"build_number" uri:"build_number" binding:"omitempty,min=0" example:"1"` // build number id
}

type GetAppBuildLogReq struct {
	Id          int64 `json:"id" uri:"id" binding:"required,min=1" example:"1"`                      // app id
	BuildNumber int   `json:"build_number" uri:"build_number" binding:"omitempty,min=0" example:"1"` // build number id
	LogNumber   int   `json:"log_number" uri:"log_number" binding:"omitempty,min=0" example:"1"`     // log number id
}

// @Summary Trigger specific app's branch or tag build process
// @Description Api to trigger specific app's branch or tag build process
// @Tags APP
// @Accept json
// @Produce json
// @Param id path integer true "App ID"
// @Param input_body body v1.CreateBuildReq true "JSON type input body"
// @Success 200 {object} v1.BuildBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/apps/{id}/builds [post]
func (c *Controller) CreateAppBuild(ctx *gin.Context) {
	// 拉取最新代码编译发版逻辑(CreateAppBuild 当前函数):
	//      前端进行单个服务发版，将相关参数传给CreateApp Api，
	//      CreateApp Api 通过drone-go模块将相关参数发给 drone-server，
	//      drone-server将任务下发给drone-agent，
	//      drone-agent 拉取该应用的源码，根据 .drone.yaml(pipeline)定义进行任务执行。
	//      当执行 .drone.yaml 最后一步(deploy_to_env)部署应用到指定环境时，会运行infra-drone-plugins中的python脚本，
	//      通过admin jwt token(galaxias_api_token)调用运维平台的 CreateDeployment Api 进行应用部署到指定环境。

	// 准备工作: 配置.drone.yaml中需要的secret
	//          docker_username = admin  // harbor用户
	//          docker_password = Se*****5  // harbor密码
	//          galaxias_api_token = eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFkbWluIiwiZHJvbmVfdG9rZW4iOiIiLCJwYXNzd29yZCI6IjMyMDdlYWQ0ZTA5MmRlNzdlMDIyMzk0YjMyMDRkNzU1Iiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNjE0MDc3NjAzLCJpYXQiOjE2MTQwNzE1NDMsImlzcyI6IlJvY2sgV2FuZyIsInN1YiI6IkxvZ2luIHRva2VuIn0.nAMR3xjGZ-4etgyVT2qfiUx2oEZhKM_iRs8lui1vTJ4  // 运维平台admin用户jwt token

	// 发布中心 -> 提交构建 -> 出现以下5个选项:
	//          (点击提交构建,请求http://10.151.3.xx:8888/v1/projects?page=1&per_page=1000&fq= 拿到所有的project，用于下面 [选择项目] 的下拉展示)
	// 选择项目: idea-aurora (下拉菜单,选择project名称)
	//          (点击选择项目,请求http://10.151.3.xx:8888/v1/projects/28/apps?page=1&per_page=1000&fq= 拿到该project下的所有app,用于 [选择应用] 的下拉展示)
	//          (点击选择项目,请求http://10.151.3.xx:8888/v1/projects/28/project-envs?page=1&per_page=1000&fq= 拿到该project下的所有project_env(固定project对应的项目环境),用于 [项目环境] 的下拉展示)
	// 选择应用: aurora-auth (下拉菜单,选择app名称)
	// 选择类型: Branch (下拉菜单,选择 Branch/Tag)
	//          (点击选择类型,请求http://10.151.3.xx:8888/v1/apps/147/branches 拿到该app所有的Branch或Tag)
	// Branch/Tag: master (下拉菜单,选择分支或者Tag名称)
	// 项目环境: 10.151.3.99-default(10.151.3.99-default) (选择部署到哪台环境)
	// 发布: 发布按钮进行该应用的发布
	//          (点击发布,请求http://10.151.3.xx:32001/v1/apps/147/builds 当前在写的这个api)
	//          http://10.151.3.xx:8888/v1/apps/147
	//          http://10.151.3.xx:8888/v1/apps/147/builds/84
	//          http://10.151.3.xx:8888/v1/apps/147/builds/84/logs/2
	//          http://10.151.3.xx:8888/v1/apps/147/builds/84
	//          http://10.151.3.xx:8888/v1/apps/147/builds/84/logs/3
	//          http://10.151.3.xx:8888/v1/apps/147/builds/84
	//          http://10.151.3.xx:8888/v1/apps/147/builds/84/logs/4
	//          http://10.151.3.xx:8888/v1/apps/147/builds/84
	//          http://10.151.3.xx:8888/v1/apps/147/builds/84/logs/5

	cfgCtx, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	var req CreateBuildReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}
	// CreateBuildReq 示例:
	// {
	//  name: "dev",
	//  type: "branch",
	//  envs: [
	//    {key: "GALAXIAS_APP_ID", value: "147"}, // app id
	//    {key: "GALAXIAS_PROJECT_ENV_ID", value: "297"}] // project_env_id
	//}

	var uriReq IdReq // app id
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	app, err := api.GetAppById(uriReq.Id)
	if err != nil {
		panic(err)
	}
	if app.GitlabProjectId == 0 {
		err := utils.NewRockError(404, 40400010, fmt.Sprintf(
			"App(%v) can't be built because app haven't associated gitlab project correctly", uriReq.Id))
		panic(err)
	}

	envs := make(map[string]string, len(req.Envs))
	for _, env := range req.Envs {
		envs[env.Key] = env.Value
	}

	build, err := drone_api.CreateBuild(cfgCtx.DroneToken, req.Type, req.Name, app.DroneRepoId, app.ProjectId, envs)
	if err != nil {
		panic(err)
	}

	resp := BuildBriefResp{}
	if err := utils.MarshalResponse(build, &resp); err != nil {
		panic(err)
	}
	c.Infof("Create app(id: %v, name: %v)'s %s type build with name %s", app.Id, app.Name, req.Type, req.Name)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get global builds info
// @Description Api to get global builds info
// @Tags BUILD
// @Accept json
// @Produce json
// @Param page query integer true "Request page number" default(1)
// @Param per_page query integer true "App number per page " default(10)
// @Param console_project_id query integer false "Console project id"
// @Success 200 {array} v1.BuildBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/builds [get]
func (c *Controller) GetGlobalBuilds(ctx *gin.Context) {
	// console_project_id 为 project_id, 但暂未发现使用场景,有待确认 ?????
	var paginationReq GetBuildPaginationReq
	if err := ctx.ShouldBind(&paginationReq); err != nil {
		panic(err)
	}

	cfgCtx, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	buildPagination, err := drone_api.GetCustomGlobalBuilds(cfgCtx.DroneToken, paginationReq.Page, paginationReq.PerPage, paginationReq.ConsoleProjectId)
	if err != nil {
		panic(err)
	}

	apps, err := api.GetAppsList()
	if err != nil {
		panic(err)
	}

	briefApps := []*models.BriefApp{}
	if err := utils.MarshalResponse(apps, &briefApps); err != nil {
		panic(err)
	}

	resp := PaginateBuildResp{}
	if err := utils.MarshalResponse(buildPagination, &resp); err != nil {
		panic(err)
	}

	extendBuildField(resp, briefApps)

	c.Logger.Infof("Get global build list, total builds number is %v", len(buildPagination.Items))
	ctx.JSON(http.StatusOK, resp)
	// {
	//    "page": 1,
	//    "per_page": 10,
	//    "total": 2,
	//    "pages": 1,
	//    "items": [
	//        {
	//            "id": 2,
	//            "repo_id": 11,
	//            "app_id": 13,
	//            "console_project_id": 5,  // console_project_id为project_id,当你需要过滤查某个项目下所有的构建信息时，可以通过下拉菜单选择某个项目下的所有构建信息
	//            "number": 2,  // 构建任务的id，可以基于此id去drone_api中查询此任务的详细信息(查构建日志)
	//            "status": "pending",
	//            "commit": "26b4808f0d35ac8f4621490166d683e255d9fed4",
	//            "branch": "master",
	//            "message": "add senseguard-guest-management service\n",
	//            "author": "someone",
	//            "author_email": "someone@sensetime.com",
	//            "created_at": 1614676798,
	//            "started_at": 0,
	//            "finished_at": 0,
	//            "enqueued_at": 1614676798
	//        },
	//        ... ...
	//    ]
	//}
}

// Get app_id from apps, And add app_id filed to resp.
// Just for add app_id to resp.
func extendBuildField(resp PaginateBuildResp, apps []*models.BriefApp) {
	repoAppMap := make(map[int64]*models.BriefApp)
	for _, app := range apps {
		if app.DroneRepoId != 0 {
			//fmt.Println("app.DroneRepoId:", app.DroneRepoId) // 11
			//fmt.Println("app:", app)                         // app: &{13 senseguard-guest-management galaxias/charts/senseguard-guest-management sense nebula guard chart: senseguard-guest-management https://gitlab.sz.sensetime.com/galaxias/charts/senseguard-guest-management.git 11 23296 {{0 63750279559 <nil>} {0 63750279559 <nil>} <nil> 0}}
			repoAppMap[app.DroneRepoId] = app
		}
	}

	for _, build := range resp.Items {
		// build: {
		//   "id": 2,
		//   "repo_id": 11, // repo_id == app.DroneRepoId
		//   "app_id": 13,
		//   "console_project_id": 5,
		//   "number": 2,
		//   "status": "pending",
		//   "enqueued_at": 1614676798,
		//   "created_at": 1614676798,
		//   "started_at": 0,
		//   "finished_at": 0,
		//   "commit": "26b4808f0d35ac8f4621490166d683e255d9fed4",
		//   "branch": "master",
		//   "author": "someone",
		//   "author_email": "someone@sensetime.com"
		// },
		app, ok := repoAppMap[build.RepoId] // get app by repoAppMap[drone_repo_id]
		if !ok {
			continue
		}
		build.AppId = app.Id // add an AppId field to resp
	}
}

// @Summary Get specific app's all builds info
// @Description Api to get all builds
// @Tags APP
// @Accept json
// @Produce json
// @Param id path integer true "App ID"
// @Param page query integer true "Request page number" default(1)
// @Param per_page query integer true "App number per page " default(10)
// @Success 200 {array} v1.BuildBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/apps/{id}/builds [get]
func (c *Controller) GetAppBuilds(ctx *gin.Context) {
	var paginationReq GetBuildsPaginationReq
	if err := ctx.ShouldBind(&paginationReq); err != nil {
		panic(err)
	}

	var uriReq IdReq // app_id
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	cfgCtx, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	app, err := api.GetAppById(uriReq.Id)
	if err != nil {
		panic(err)
	}

	buildPagination, err := drone_api.GetCustomBuilds(cfgCtx.DroneToken, app.DroneRepoId, paginationReq.Page, paginationReq.PerPage)
	if err != nil {
		panic(err)
	}

	resp := PaginateBuildResp{}
	if err := utils.MarshalResponse(buildPagination, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get builds of project_id(%v)'s App(%v), total builds number is %v", app.ProjectId, app.Name, len(buildPagination.Items))
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get build by id
// @Description Api for get build by id
// @Tags APP
// @Accept json
// @Produce json
// @Param id path integer true "App ID"
// @Param build_number path integer true "Build number"
// @Success 200 {object} v1.BuildDetailResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/apps/{id}/builds/{build_number} [get]
func (c *Controller) GetAppBuild(ctx *gin.Context) {
	// 发布中心 -> 构建历史 -> 点击"查看",会出现下面描述的 "三步", 来查选择的某个构建历史信息(这里触发 GetBuildLogs 获取详细日志信息)
	// 第一步: GET 请求 /v1/apps/{app_id} 确认此构建历史的app的app_id存在 -> 此为api GetApp 做的事情
	//        示例: http://10.151.3.xx:88888/v1/apps/203
	// 第二步: GET 请求 /v1/apps/{app_id}/builds/{build_number} -> 此为当前api GetAppBuild 做的事情
	//        示例: http://10.151.3.xx:88888/v1/apps/203/builds/27
	// {
	//    "id": 34209,
	//    "repo_id": 188,
	//    "console_project_id": 28,
	//    "number": 27,  // build_number
	//    "status": "failure",
	//    "commit": "afc621d5e9ee8660760d7bb73699ec22325ed369",
	//    "branch": "temp",
	//    "message": "Merge branch 'temp' of https://gitlab.sz.sensetime.com/idea-aurora/service/aurora-pedestrian-process-service into temp\n",
	//    "author": "someone",
	//    "author_email": "someone@sensetime.com",
	//    "created_at": 1614770469,
	//    "started_at": 1614770469,
	//    "finished_at": 1614770474,
	//    "enqueued_at": 1614770469,
	//    "procs": [
	//        {
	//            "id": 170729,
	//            "pid": 1,
	//            "ppid": 0,
	//            "pgid": 1,
	//            "name": "",
	//            "state": "failure",  // 当前构建任务的整体情况: 失败。 其他状态有: success/failure/running/pending/error
	//            "exit_code": 1,
	//            "start_time": 1614770469,
	//            "end_time": 1614770474,
	//            "machine": "fis-devops-0",
	//            "children": [
	//                {
	//                    "id": 170730,
	//                    "pid": 2,
	//                    "ppid": 1,
	//                    "pgid": 2,
	//                    "name": "clone",
	//                    "state": "failure", // 当前构建任务的第一步: 失败。 其他状态有: success/failure/running/pending
	//                    "exit_code": 1,
	//                    "start_time": 1614770469,
	//                    "end_time": 1614770474,
	//                    "machine": "fis-devops-0"
	//                },
	//                {
	//                    "id": 170731,
	//                    "pid": 3,
	//                    "ppid": 1,
	//                    "pgid": 3,
	//                    "name": "build-and-push-image",
	//                    "state": "skipped",  // 当前构建任务的第二步: 还没开始 (因为第一步失败了)
	//                    "exit_code": 0,
	//                    "start_time": 0,
	//                    "end_time": 0,
	//                    "machine": ""
	//                },
	//                {
	//                    "id": 170732,
	//                    "pid": 4,
	//                    "ppid": 1,
	//                    "pgid": 4,
	//                    "name": "package_and_upload_chart",
	//                    "state": "skipped",  // 同上
	//                    "exit_code": 0,
	//                    "start_time": 0,
	//                    "end_time": 0,
	//                    "machine": ""
	//                },
	//                {
	//                    "id": 170733,
	//                    "pid": 5,
	//                    "ppid": 1,
	//                    "pgid": 5,
	//                    "name": "deploy_to_env",
	//                    "state": "skipped",
	//                    "exit_code": 0,
	//                    "start_time": 0,
	//                    "end_time": 0,
	//                    "machine": ""
	//                }
	//            ]
	//        }
	//    ]
	//}
	// 第三步: GET 请求 /v1/apps/{app_id}/builds/{build_number}/logs/{log_number} -> 此为api GetBuildLogs 做的事情
	//        但是,这里log_number 具体为多少,需要前端根据 第二步 返回的数据进行判断.
	//        具体前端逻辑如下:
	//        1 前端先检测 procs.state == "success/failure/running/pending/error等"，则展示总状态为对应状态.
	//        2 然后再 for 循环检测 procs.children, 查找每个 procs.children下的state:
	//           2.1 假设现在循环到 pid == 1, 则请求页面: (示例)http://10.151.3.xx:88888/v1/apps/174/builds/184/logs/1 (最后的log_number为1,表示请求第一阶段的日志)
	//               如果procs.children.state == running 或 failure 则break，返回当前信息到前端页面.
	//               如果procs.children.state == success 则展示当前的日志给前端，然后去查看 pid == 2的日志(即(示例)http://10.151.3.xx:88888/v1/apps/174/builds/184/logs/2)
	//           2.2 定时任务: 当 procs.children.state == running 或 pending 时,每隔2000毫秒再执行一下上面 2.1 的操作.
	//           2.3 当循环到最后一个pid(pid == 5)时, procs.children.state == success 还是成功,那么展示对应日志,并结束所有请求.(即(示例)http://10.151.3.xx:88888/v1/apps/174/builds/184/logs/5)
	//        【注意】: A. 如果 procs.state == "pending",是drone-agent没有起,所以会一直pending.因为当前的任务不知道往哪里调.
	//                 B. 其他的报错,一般都是java源码进行docker build失败了,所以需要让开发同事先确认代码/Dockerfile没有问题,才能发版.

	var uriReq GetAppBuildReq // app id & build number
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	cfgCtx, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	app, err := api.GetAppById(uriReq.Id)
	if err != nil {
		panic(err)
	}

	build, err := drone_api.GetCustomBuild(cfgCtx.DroneToken, app.DroneRepoId, uriReq.BuildNumber)
	if err != nil {
		panic(err)
	}
	// build:
	// {
	//    "id": 1,
	//    "repo_id": 11,
	//    "console_project_id": 5,
	//    "number": 1,
	//    "parent": 0,
	//    "event": "deployment",
	//    "status": "pending",
	//    "error": "",
	//    "enqueued_at": 1614655426,
	//    "created_at": 1614655426,
	//    "started_at": 0,
	//    "finished_at": 0,
	//    "deploy_to": "",
	//    "commit": "26b4808f0d35ac8f4621490166d683e255d9fed4",
	//    "branch": "master",
	//    "ref": "refs/heads/master",
	//    "refspec": "",
	//    "remote": "",
	//    "title": "",
	//    "message": "add senseguard-guest-management service\n",
	//    "timestamp": 0,
	//    "sender": "",
	//    "author": "someone",
	//    "author_avatar": "https://www.gravatar.com/avatar/44f1af844f14d167aaa69014a5176353.jpg?s=128",
	//    "author_email": "someone@sensetime.com",
	//    "link_url": "https://gitlab.sz.sensetime.com/galaxias/charts/senseguard-guest-management",
	//    "reviewed_by": "",
	//    "reviewed_at": 0,
	//    "procs": [
	//        {
	//            "id": 1,
	//            "pid": 1,
	//            "ppid": 0,
	//            "pgid": 1,
	//            "name": "",
	//            "state": "pending",
	//            "exit_code": 0,
	//            "children": [
	//                {
	//                    "id": 2,
	//                    "pid": 2,
	//                    "ppid": 1,
	//                    "pgid": 2,
	//                    "name": "clone",
	//                    "state": "pending",
	//                    "exit_code": 0
	//                },
	//                {
	//                    "id": 3,
	//                    "pid": 3,
	//                    "ppid": 1,
	//                    "pgid": 3,
	//                    "name": "build-and-push-image",
	//                    "state": "pending",
	//                    "exit_code": 0
	//                },
	//                {
	//                    "id": 4,
	//                    "pid": 4,
	//                    "ppid": 1,
	//                    "pgid": 4,
	//                    "name": "package_and_upload_chart",
	//                    "state": "pending",
	//                    "exit_code": 0
	//                },
	//                {
	//                    "id": 5,
	//                    "pid": 5,
	//                    "ppid": 1,
	//                    "pgid": 5,
	//                    "name": "deploy_to_env",
	//                    "state": "pending",
	//                    "exit_code": 0
	//                }
	//            ]
	//        }
	//    ]
	//}
	resp := BuildDetailResp{}
	if err := utils.MarshalResponse(build, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get build by app_id(%v) and build_number(%v)", uriReq.Id, uriReq.BuildNumber)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get specific app build log by log number
// @Description Api for get specific app build log by log number
// @Tags APP
// @Accept json
// @Produce json
// @Param id path integer true "App ID"
// @Param build_number path integer true "Build number"
// @Param log_number path integer true "Log number"
// @Success 200 {object} v1.BuildDetailResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/apps/{id}/builds/{build_number}/logs/{log_number} [get]
func (c *Controller) GetBuildLogs(ctx *gin.Context) {
	// 详细逻辑见: GetAppBuild 函数中的注释
	var uriReq GetAppBuildLogReq
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		panic(err)
	}

	cfgCtx, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	app, err := api.GetAppById(uriReq.Id)
	if err != nil {
		panic(err)
	}

	logs, err := drone_api.GetCustomLogs(cfgCtx.DroneToken, app.DroneRepoId, uriReq.BuildNumber, uriReq.LogNumber)
	if err != nil {
		panic(err)
	}

	c.Logger.Infof("Get project_id(%v) app(%v)'s build logs by app_id(%v) and build_number(%v) and log_number(%v)",
		app.ProjectId, app.Name, uriReq.Id, uriReq.BuildNumber, uriReq.LogNumber)
	ctx.JSON(http.StatusOK, logs)
}
