package models

// 记录基于从运维平台(前端页面到后端drone)向指定集群环境的单个chart服务的版本发布的信息记录
// 也就是说，某个集群环境是否安装了某个服务的某个版本，是可以从instance表中来查看到的
// 当然如果是检测某个环境是否有这个版本用instance表去检查不是完全准确的，所以这里的作用主要是记录信息而已

type Instance struct {
	Id             int64  `json:"id" gorm:"unique_index;not null"`
	ClusterName    string `json:"cluster_name"`
	EnvNamespace   string `json:"env_namespace"`
	ProjectName    string `json:"project_name" gorm:"not null"`
	Name           string `json:"name"`       // helm deploy in cluster name, example: senseguard-td-result-consume-default
	ChartName      string `json:"chart_name"` // helm chart name, example: senseguard-td-result-consume
	ChartVersion   string `json:"chart_version"`
	LastDeployment int64  `json:"last_deployment" gorm:"not null"` // deployment_id
	AppId          int64  `json:"app_id" gorm:"not null"`
	EnvId          int64  `json:"env_id" gorm:"not null"`
	Common
}

type InstancePagination struct {
	PageNum  int64       `json:"page_num"`
	PageSize int64       `json:"page_size"`
	Total    int64       `json:"total"`
	Pages    int64       `json:"pages"`
	Items    []*Instance `json:"items"`
}
