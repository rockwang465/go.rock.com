package models

type ProjectEnv struct {
	Id          int64  `json:"id" gorm:"unique_index;not null"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	EnvId       int64  `json:"env_id" gorm:"not null"`
	ProjectId   int64  `json:"project_id" gorm:"not null"`
	Common
}

type BriefProjectEnv struct {
	Id          int64  `json:"id" gorm:"unique_index;not null"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	EnvId       int64  `json:"env_id" gorm:"not null"`
	ProjectId   int64  `json:"project_id" gorm:"not null"`
	Common
}

type ProjectEnvPagination struct {
	PageNum  int64              `json:"page_num"`
	PageSize int64              `json:"page_size"`
	Total    int64              `json:"total"`
	Pages    int64              `json:"pages"`
	Items    []*BriefProjectEnv `json:"items"`
}

//func (*BriefProjectEnv) TableName() string {
//	return "project_env"
//}
