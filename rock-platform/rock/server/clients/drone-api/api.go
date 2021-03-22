package drone_api

//import "github.com/drone/drone-go/drone" // 之前用的是这个项目的v0.8.4的tag，并修改了源码
import (
	"github.com/rockwang465/drone/drone-go/drone" // 所以不用被改动的源码，而是直接拿改好的源码放在自己的项目下
)

// 所有使用drone client(当前api.go)对drone server(放在linux上的二进制服务)请求，保存的数据都放在了mysql数据库中。
// 也就是说drone server是基于mysql保存数据的。

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

// create a custom docker registry
func CreateRegistry(jwt, addr, user, pwd string) (*drone.Registry, error) {
	client, err := getClient(jwt)
	if err != nil {
		return nil, err
	}
	return client.RegistryCustomCreate(addr, user, pwd)
}

// get all docker registry list
func GetRegistries(jwt string) ([]*drone.Registry, error) {
	client, err := getClient(jwt)
	if err != nil {
		return nil, err
	}
	return client.RegistryCustomList()
}

// get an registry by IP address
func GetRegistry(jwt, addr string) (*drone.Registry, error) {
	client, err := getClient(jwt)
	if err != nil {
		return nil, err
	}

	return client.CustomRegistry(addr)
}

// delete an registry by IP address
func DeleteRegistry(jwt, addr string) error {
	client, err := getClient(jwt)
	if err != nil {
		return err
	}
	return client.CustomRegistryDelete(addr)
}

// update an registry info(username,password) by IP address
func UpdateRegistry(jwt, addr, user, pwd string) (*drone.Registry, error) {
	client, err := getClient(jwt)
	if err != nil {
		return nil, err
	}
	return client.CustomRegistryUpdate(addr, user, pwd)
}

// create a secret
func CreateSecret(jwt, name, value string) (*drone.Secret, error) {
	client, err := getClient(jwt)
	if err != nil {
		return nil, err
	}
	return client.SecretCustomCreate(name, value)
}

// get all secrets
func GetSecrets(jwt string) ([]*drone.Secret, error) {
	client, err := getClient(jwt)
	if err != nil {
		return nil, err
	}
	return client.SecretCustomList()
}

// get a secret by name
func GetSecret(jwt, name string) (*drone.Secret, error) {
	client, err := getClient(jwt)
	if err != nil {
		return nil, err
	}
	return client.CustomSecret(name)
}

// delete a secret by name
func DeleteSecret(jwt, name string) error {
	client, err := getClient(jwt)
	if err != nil {
		return err
	}
	return client.CustomSecretDelete(name)
}

// update a secret info(value) by name
func UpdateSecret(jwt, name, value string) (*drone.Secret, error) {
	client, err := getClient(jwt)
	if err != nil {
		return nil, err
	}
	return client.CustomSecretUpdate(name, value)
}

// get app branches by app gitlab project id
func RepoBranches(jwt string, repoId int64) ([]*drone.Branch, error) {
	client, err := getClient(jwt)
	if err != nil {
		return nil, err
	}
	return client.RepoBranches(repoId)
}

// get app tags by app gitlab project id
func RepoTags(jwt string, repoId int64) ([]*drone.Tag, error) {
	client, err := getClient(jwt)
	if err != nil {
		return nil, err
	}
	return client.RepoTags(repoId)
}

// create a build task by triggerType(app branch/tag), triggerName(branch/tag name), repoId(drone_repo_id), rockProjectId(project_id), env(app_id, project_env_id)
func CreateBuild(jwt, triggerType, triggerName string, repoId, rockProjectId int64, envs map[string]string) (*drone.Build, error) {
	client, err := getClient(jwt)
	if err != nil {
		return nil, err
	}
	return client.CreateBuild(repoId, rockProjectId, triggerType, triggerName, envs)
}

// get global build list info by page(page_num), perPage(page_size), console_project_id
func GetCustomGlobalBuilds(jwt string, page, perPage, cId int64) (*drone.PaginateBuild, error) {
	client, err := getClient(jwt)
	if err != nil {
		return nil, err
	}
	return client.GetCustomGlobalBuildList(page, perPage, cId)
}

// get all builds info by repoId, page(page_num), perPage(page_size)
func GetCustomBuilds(jwt string, repoId, page, perPage int64) (*drone.PaginateBuild, error) {
	client, err := getClient(jwt)
	if err != nil {
		return nil, err
	}
	return client.GetCustomBuildList(repoId, page, perPage)
}

// get specific build info by repo_id and build_number
func GetCustomBuild(jwt string, repoId int64, buildNumber int) (*drone.Build, error) {
	client, err := getClient(jwt)
	if err != nil {
		return nil, err
	}
	return client.CustomBuild(repoId, buildNumber)
}

// get specific app build logs by repo_id and buildNumber and logNumber
func GetCustomLogs(jwt string, repoId int64, buildNumber, logNumber int) ([]*drone.Log, error) {
	client, err := getClient(jwt)
	if err != nil {
		return nil, err
	}
	return client.BuildCustomLogs(repoId, buildNumber, logNumber)
}
