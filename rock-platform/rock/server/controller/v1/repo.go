package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/rockwang465/drone/drone-go/drone"
	"go.rock.com/rock-platform/rock/server/client/drone-api"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
)

type DroneRepoBriefResp struct {
	ProjectId int64  `json:"project_id" example:"1"`
	Owner     string `json:"owner" example:"admin"`
	Name      string `json:"name" example:"infra-console"`
	FullName  string `json:"full_name" example:"fis-infra/infra-console"`
	Clone     string `json:"clone_url" example:"http://gitlab.sz.sensetime.com/fis-infra/infra-console.git"`
	IsAdded   bool   `json:"is_added" example:"false"`
}

// get current use has gitlab repos by drone token

// @Summary Get remote all repos
// @Description Api to get remote(gitlab) all repos
// @Tags REPO
// @Accept json
// @Produce json
// @Success 200 {object} v1.DroneRepoBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/repos [get]
func (c *Controller) GetRemoteRepos(ctx *gin.Context) {
	cfgCtx, err := utils.GetConfCtx(ctx)
	if err != nil {
		panic(err)
	}

	// get current user gitlab apps by drone token
	remoteRepo, err := drone_api.GetRemoteRepos(cfgCtx.DroneToken)
	if err != nil {
		panic(err)
	}
	//k:0, v:&drone.RemoteRepo{ID:0, ProjectId:22919, Owner:"idea-aurora", Name:"aurora-automobile-process-service-api", FullName:"idea-aurora/service-api/aurora-automobile-process-service-api", Link:"https://gitlab.sz.sensetime.com/idea-aurora/service-api/aurora-automobile-process-service-api", Clone:"https://gitlab.sz.sensetime.com/idea-aurora/service-api/aurora-automobile-process-service-api.git", Branch:"master", Visibility:"", IsPrivate:true, IsTrusted:false, IsGated:false, Active:false, AllowPull:false, AllowPush:false, AllowDeploy:false, AllowTag:false, LastBuild:0, Config:"", IsAdded:false}
	//k:1, v:&drone.RemoteRepo{ID:0, ProjectId:22918, Owner:"idea-aurora", Name:"aurora-pedestrian-process-service-api", FullName:"idea-aurora/service-api/aurora-pedestrian-process-service-api", Link:"https://gitlab.sz.sensetime.com/idea-aurora/service-api/aurora-pedestrian-process-service-api", Clone:"https://gitlab.sz.sensetime.com/idea-aurora/service-api/aurora-pedestrian-process-service-api.git", Branch:"master", Visibility:"", IsPrivate:true, IsTrusted:false, IsGated:false, Active:false, AllowPull:false, AllowPush:false, AllowDeploy:false, AllowTag:false, LastBuild:0, Config:"", IsAdded:false}
	//k:2, v:&drone.RemoteRepo{ID:0, ProjectId:22915, Owner:"idea-aurora", Name:"symphony-research", FullName:"idea-aurora/symphony-research", Link:"https://gitlab.sz.sensetime.com/idea-aurora/symphony-research", Clone:"https://gitlab.sz.sensetime.com/idea-aurora/symphony-research.git", Branch:"master", Visibility:"", IsPrivate:true, IsTrusted:false, IsGated:false, Active:false, AllowPull:false, AllowPush:false, AllowDeploy:false, AllowTag:false, LastBuild:0, Config:"", IsAdded:false}
	//k:3, v:&drone.RemoteRepo{ID:0, ProjectId:22768, Owner:"idea-aurora", Name:"aurora-automobile-process-service", FullName:"idea-aurora/service/aurora-automobile-process-service", Link:"https://gitlab.sz.sensetime.com/idea-aurora/service/aurora-automobile-process-service", Clone:"https://gitlab.sz.sensetime.com/idea-aurora/service/aurora-automobile-process-service.git", Branch:"master", Visibility:"", IsPrivate:true, IsTrusted:false, IsGated:false, Active:false, AllowPull:false, AllowPush:false, AllowDeploy:false, AllowTag:false, LastBuild:0, Config:"", IsAdded:false}
	//k:4, v:&drone.RemoteRepo{ID:0, ProjectId:22767, Owner:"idea-aurora", Name:"aurora-pedestrian-process-service", FullName:"idea-aurora/service/aurora-pedestrian-process-service", Link:"https://gitlab.sz.sensetime.com/idea-aurora/service/aurora-pedestrian-process-service", Clone:"https://gitlab.sz.sensetime.com/idea-aurora/service/aurora-pedestrian-process-service.git", Branch:"master", Visibility:"", IsPrivate:true, IsTrusted:false, IsGated:false, Active:false, AllowPull:false, AllowPush:false, AllowDeploy:false, AllowTag:false, LastBuild:0, Config:"", IsAdded:false}
	//k:5, v:&drone.RemoteRepo{ID:0, ProjectId:22764, Owner:"SenseNebula-Guard-std", Name:"v2.2.0-to-v2.3.0", FullName:"SenseNebula-Guard-std/version-upgrade/v2.2.0-to-v2.3.0", Link:"https://gitlab.sz.sensetime.com/SenseNebula-Guard-std/version-upgrade/v2.2.0-to-v2.3.0", Clone:"https://gitlab.sz.sensetime.com/SenseNebula-Guard-std/version-upgrade/v2.2.0-to-v2.3.0.git", Branch:"master", Visibility:"", IsPrivate:true, IsTrusted:false, IsGated:false, Active:false, AllowPull:false, AllowPush:false, AllowDeploy:false, AllowTag:false, LastBuild:0, Config:"", IsAdded:false}

	// get all apps in mysql
	apps, err := api.GetAppsList()
	if err != nil {
		panic(err)
	}

	// compare gitlab apps.FullName and mysql apps.FullName
	diffRepos := markIsAdded(remoteRepo, apps)

	resp := []*DroneRepoBriefResp{}
	if err := utils.MarshalResponse(diffRepos, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("User id:%v name:%v get all remote repo list", cfgCtx.UserId, cfgCtx.Username)
	ctx.JSON(http.StatusOK, resp)
}

// compare apps.FullName is in remoteRepo.FullName, if exist then remoteRepo.isAdded = true
func markIsAdded(remoteRepo []*drone.RemoteRepo, apps []*models.App) []*drone.RemoteRepo {
	remoteReposMap := make(map[string]*drone.RemoteRepo)
	for _, repo := range remoteRepo {
		remoteReposMap[repo.FullName] = repo
	}

	for _, app := range apps {
		rp, ok := remoteReposMap[app.FullName]
		if ok {
			rp.IsAdded = true
		}
	}

	remoteReposList := make([]*drone.RemoteRepo, 0)
	for _, repo := range remoteReposMap {
		remoteReposList = append(remoteReposList, repo)
	}

	return remoteReposList
}
