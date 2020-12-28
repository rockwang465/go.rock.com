package models

import "fmt"

type Role struct {
	Id          int64         `json:"id" gorm:"unique_index;not null"` // primary_key, AUTO_INCREMENT
	Name        string        `json:"name" gorm:"unique_index;not null"`
	Description string        `json:"description"`
	Users       []*User       `json:"users" gorm:"-"` // - 忽略这个字段,不进行映射
	Permissions []*Permission `json:"permissions" gorm:"-"`
	Common
}

type RolePagination struct {
	Page    int64   `json:"page"`
	PerPage int64   `json:"per_page"`
	Total   int64   `json:"total"`
	Pages   int64   `json:"pages"`
	Items   []*Role `json:"items"`
}

func (r *Role) String() string {
	return fmt.Sprintf("Role:%v", r.Name)
}
