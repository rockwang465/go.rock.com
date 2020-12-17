package models

type Permission struct {
	Id          int64   `json:"id" xorm:"pk autoincr unique notnull 'id'"`
	Name        string  `json:"name" xorm:"unique notnull"`
	Description string  `json:"description"`
	Roles       []*Role `json:"roles" xorm:"-"`
	//Common      `xorm:"extends"`
}

//type PermissionPagination struct {
//	Page    int64         `json:"page"`
//	PerPage int64         `json:"per_page"`
//	Total   int64         `json:"total"`
//	Pages   int64         `json:"pages"`
//	Items   []*Permission `json:"items"`
//}
//
//type PermissionRoleMap struct {
//	*Permission           `xorm:"extends"`
//	*RolePermissionMapper `xorm:"extends"`
//	*Role                 `xorm:"extends"`
//}
//
//type RolePermissionMap struct {
//	*Role                 `xorm:"extends"`
//	*RolePermissionMapper `xorm:"extends"`
//	*Permission           `xorm:"extends"`
//}
