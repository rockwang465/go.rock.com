package models

import "fmt"

type Role struct {
	Id          int64         `json:"id" gorm:"unique_index;not null"` // primary_key, AUTO_INCREMENT
	Name        string        `json:"name" gorm:"not null"`
	Description string        `json:"description"`
	Users       []*User       `json:"users" gorm:"-"` // - ignore field, not mapping
	Permissions []*Permission `json:"permissions" gorm:"-"`
	Common
}

type RolePagination struct {
	PageNum  int64   `json:"page_num"`
	PageSize int64   `json:"page_size"`
	Total    int64   `json:"total"`
	Pages    int64   `json:"pages"`
	Items    []*Role `json:"items"`
}

func (r *Role) String() string {
	return fmt.Sprintf("Role:%v", r.Name)
}
