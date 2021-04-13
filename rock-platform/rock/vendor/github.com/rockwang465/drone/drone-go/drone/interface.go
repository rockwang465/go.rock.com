package drone

import "net/http"

// Client is used to communicate with a Drone server.
type Client interface {
	// SetClient sets the http.Client.
	SetClient(*http.Client)

	// SetAddress sets the server address.
	SetAddress(string)

	// Self returns the currently authenticated user.
	Self() (*User, error)

	// User returns a user by login.
	User(string) (*User, error)

	// UserJwt auth returns jwt token.
	UserJwt(username, token string) (*Token, error)

	// UserList returns a list of all registered users.
	UserList() ([]*User, error)

	// UserPost creates a new user account.
	UserPost(*User) (*User, error)

	// UserPatch updates a user account.
	UserPatch(*User) (*User, error)

	// UserDel deletes a user account.
	UserDel(string) error

	// Repo returns a repository by name.
	Repo(string, string) (*Repo, error)

	// RepoBranches returns branch list of that repo in the gitlab.
	RepoBranches(int64) ([]*Branch, error)

	// RepoTags returns tag list of that repo in the gitlab.
	RepoTags(int64) ([]*Tag, error)

	// RepoList returns a list of all repositories to which the user has explicit
	// access in the host system.
	RepoList() ([]*Repo, error)

	// RepoList returns a list of all repositories to which the user has explicit
	// access in the host system.
	RemoteRepos() ([]*RemoteRepo, error)

	// SyncRemoteRepo sync gitlab data into drone
	SyncRemoteRepo(int64) (*RemoteRepo, error)

	// RepoPost activates a repository.
	RepoPost(int64) (*Repo, error)

	// RepoPatch updates a repository.
	RepoPatch(string, string, *RepoPatch) (*Repo, error)

	// RepoMove moves the repository
	RepoMove(string, string, string) error

	// RepoChown updates a repository owner.
	RepoChown(string, string) (*Repo, error)

	// RepoRepair repairs the repository hooks.
	RepoRepair(string, string) error

	// RepoDel deletes a repository.
	RepoDel(string, string) error

	// Build returns a repository build by number.
	Build(string, string, int) (*Build, error)

	// CreateBuild trigger a repository build.
	CreateBuild(int64, int64, string, string, map[string]string) (*Build, error)

	// Build returns a repository build by number.
	CustomBuild(int64, int) (*Build, error)

	// BuildLast returns the latest repository build by branch. An empty branch
	// will result in the default branch.
	BuildLast(string, string, string) (*Build, error)

	// BuildList returns a list of recent builds for the
	// the specified repository.
	BuildList(string, string) ([]*Build, error)

	// BuildList returns a list of recent builds for the
	// the specified repository.
	GetCustomBuildList(int64, int64, int64) (*PaginateBuild, error)

	// BuildList returns a list of recent builds for the
	// the specified repository.
	GetCustomGlobalBuildList(int64, int64, int64) (*PaginateBuild, error)

	// BuildQueue returns a list of enqueued builds.
	BuildQueue() ([]*Activity, error)

	// BuildStart re-starts a stopped build.
	BuildStart(string, string, int, map[string]string) (*Build, error)

	// BuildStop stops the specified running job for given build.
	BuildStop(string, string, int, int) error

	// BuildApprove approves a blocked build.
	BuildApprove(string, string, int) (*Build, error)

	// BuildDecline declines a blocked build.
	BuildDecline(string, string, int) (*Build, error)

	// BuildKill force kills the running build.
	BuildKill(string, string, int) error

	// Deploy triggers a deployment for an existing build using the specified
	// target environment.
	Deploy(string, string, int, string, map[string]string) (*Build, error)

	// Registry returns a registry by hostname.
	Registry(owner, name, hostname string) (*Registry, error)

	// Custom Registry returns a registry by hostname.
	CustomRegistry(name string) (*Registry, error)

	// RegistryList returns a list of all repository registries.
	RegistryList(owner, name string) ([]*Registry, error)

	// RegistryList returns a list of all repository registries.
	RegistryCustomList() ([]*Registry, error)

	// RegistryCreate creates a registry.
	RegistryCreate(owner, name string, registry *Registry) (*Registry, error)

	// CustomRegistryCreate creates a registry.
	RegistryCustomCreate(addr, user, pwd string) (*Registry, error)

	// RegistryUpdate updates a registry.
	RegistryUpdate(owner, name string, registry *Registry) (*Registry, error)

	// RegistryUpdate updates a registry.
	CustomRegistryUpdate(addr, user, pwd string) (*Registry, error)

	// RegistryDelete deletes a registry.
	RegistryDelete(owner, name, hostname string) error

	// CustomRegistryDelete deletes a registry by name.
	CustomRegistryDelete(name string) error

	// Secret returns a secret by name.
	Secret(owner, name, secret string) (*Secret, error)

	// Secret returns a secret by name.
	CustomSecret(name string) (*Secret, error)

	// SecretList returns a list of all repository secrets.
	SecretList(owner, name string) ([]*Secret, error)

	// SecretList returns a list of all repository secrets.
	SecretCustomList() ([]*Secret, error)

	// SecretCreate creates a registry.
	SecretCreate(owner, name string, secret *Secret) (*Secret, error)

	// CustomRegistryCreate creates a registry.
	SecretCustomCreate(name, value string) (*Secret, error)

	// SecretUpdate updates a registry.
	SecretUpdate(owner, name string, secret *Secret) (*Secret, error)

	// SecretUpdate updates a registry.
	CustomSecretUpdate(name, value string) (*Secret, error)

	// SecretDelete deletes a secret.
	SecretDelete(owner, name, secret string) error

	// SecretDelete deletes a secret.
	CustomSecretDelete(name string) error

	// Server returns the named servers details.
	Server(name string) (*Server, error)

	// ServerList returns a list of all active build servers.
	ServerList() ([]*Server, error)

	// BuildLogs returns a list of logs of a number of build.
	BuildLogs(owner, name string, num, job int) ([]*Log, error)

	// BuildLogs returns a list of logs of a number of build.
	BuildCustomLogs(repoId int64, num, job int) ([]*Log, error)
}
