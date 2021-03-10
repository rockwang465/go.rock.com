package models

// 单个应用的对应集群环境名称空间下的 配置文件(values.yaml)

type AppConf struct {
	Id           int64  `json:"id" gorm:"unique_index;not null"`
	Config       string `json:"config" gorm:"type:varchar(12800)"`
	AppId        int64  `json:"app_id" gorm:"not null"`
	ProjectEnvId int64  `json:"project_env_id" gorm:"not null"`
	Common
}

type AppConfPagination struct {
	PageNum  int64      `json:"page_num"`
	PageSize int64      `json:"page_size"`
	Total    int64      `json:"total"`
	Pages    int64      `json:"pages"`
	Items    []*AppConf `json:"items"`
}
