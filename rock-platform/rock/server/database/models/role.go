package models

import "fmt"

type Role struct {
	Id          int64         `json:"id" xorm:"pk autoincr unique notnull 'id'"`
	Name        string        `json:"name" xorm:"unique notnull"`
	Description string        `json:"description"`
	Users       []*User       `json:"users" xorm:"-"`
	Permissions []*Permission `json:"permissions" xorm:"-"`
	//Common      `xorm:"extends"`
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
