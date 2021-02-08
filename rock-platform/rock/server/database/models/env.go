package models

type Env struct {
	Id          int64  `json:"id" gorm:"unique_index;not null"`
	Description string `json:"description"`
	Namespace   string `json:"namespace" gorm:"not null"`
	ClusterId   int64  `json:"cluster_id" gorm:"not null"`
	Common
}

//type BriefEnv struct {
//	Id          int64  `json:"id" gorm:"unique_index;not null"`
//	Description string `json:"description"`
//	Namespace   string `json:"namespace" gorm:"not null"`
//	Common
//}
//
//func (*BriefEnv) TableName() string {
//	return "env"
//}

type EnvPagination struct {
	PageNum  int64  `json:"page_num"`
	PageSize int64  `json:"per_size"`
	Total    int64  `json:"total"`
	Pages    int64  `json:"pages"`
	Items    []*Env `json:"items"`
}

//type BriefEnvPagination struct {
//	PageNum  int64       `json:"page_num"`
//	PageSize int64       `json:"per_size"`
//	Total    int64       `json:"total"`
//	Pages    int64       `json:"pages"`
//	Items    []*BriefEnv `json:"items"`
//}
