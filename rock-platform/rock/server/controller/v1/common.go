package v1

type GetPaginationReq struct {
	PageSize   int64  `json:"page_size" binding:"required" example:"20"`
	PageNum    int64  `json:"page_num" binding:"required" example:"1"`
	QueryField string `json:"query_field" binding:"required" example:"rock"`
}
