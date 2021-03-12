package models

// 单个服务部署(helm chart)的信息记录

type Deployment struct {
	Id           int64  `json:"id" gorm:"unique_index;not null"`
	Name         string `json:"name" gorm:"not null"`
	Description  string `json:"description"`
	ChartName    string `json:"chart_name"`
	ChartVersion string `json:"chart_version"`
	AppId        int64  `json:"app_id"`
	EnvId        int64  `json:"env_id"`
	//Executor     *User  `json:"executor"` // User table
	Common
}

type DeploymentPagination struct {
	PageNum  int64         `json:"page_num"`
	PageSize int64         `json:"page_size"`
	Total    int64         `json:"total"`
	Pages    int64         `json:"pages"`
	Items    []*Deployment `json:"items"`
}
