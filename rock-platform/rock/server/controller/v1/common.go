package v1

type GetPaginationReq struct {
	PageSize   int64  `json:"page_size" form:"page_size" binding:"required,min=1" example:"20"`
	PageNum    int64  `json:"page_num" form:"page_num" binding:"required,min=1" example:"1"`
	QueryField string `json:"query_field" form:"query_field" binding:"omitempty" example:"rock"` // omitempty: allow empty
}

type IdReq struct {
	Id int64 `json:"id" uri:"id" binding:"required,min=1" example:"1"`
}

type UpdateUserPwdReq struct {
	OldPassword string `json:"old_password" binding:"required" example:"********"`
	NewPassword string `json:"new_password" binding:"required" example:"********"`
}
