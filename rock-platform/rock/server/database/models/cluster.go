package models

// user add a k8s cluster information
type Cluster struct {
	Id          int64  `json:"id" gorm:"unique_index;not null"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	Config      string `json:"config" gorm:"type:varchar(12800)"` // /etc/kubernetes/admin.conf
	Common
}

type ClusterPagination struct {
	PageNum  int64      `json:"page_num"`
	PageSize int64      `json:"page_size"`
	Total    int64      `json:"total"`
	Pages    int64      `json:"pages"`
	Items    []*Cluster `json:"items"`
}
