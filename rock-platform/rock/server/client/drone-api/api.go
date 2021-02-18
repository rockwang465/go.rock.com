package drone_api

//import "github.com/drone/drone-go/drone" // 顺义用的是这个项目的v0.8.4的tag，并修改了源码
import (
	"github.com/rockwang465/drone/drone-go/drone" // 所以不用被改动的源码，而是直接拿改好的源码放在自己的项目下
)

// use accessToken(gitlab access token) and username, generate drone token
func GetJwt(username, accessToken string) (*drone.Token, error) {
	client := getRawClient()
	return client.UserJwt(username, accessToken)
}

// get remote all repos by drone token
func GetRemoteRepos(jwt string) ([]*drone.RemoteRepo, error) {
	client, err := getClient(jwt)
	if err != nil {
		return nil, err
	}
	return client.RemoteRepos()
}

// use jwt(drone token) generate drone client, get gitlab project id repository
func SyncRemoteRepo(jwt string, gitlabProjectId int64) (*drone.RemoteRepo, error) {
	client, err := getClient(jwt)
	if err != nil {
		return nil, err
	}
	return client.SyncRemoteRepo(gitlabProjectId)
}

// curl http://gitlab.sz.sensetime.com/api/v4/projects/11222
// get repo information by drone repo id
func ActiveRepo(jwt string, repoId int64) (*drone.Repo, error) {
	client, err := getClient(jwt)
	if err != nil {
		return nil, err
	}
	return client.RepoPost(repoId)
}
