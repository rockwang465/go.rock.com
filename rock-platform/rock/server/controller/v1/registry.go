package v1

import "github.com/gin-gonic/gin"

// 用户注册页面

type RegisterInfo struct {
	Username  string `json:"username" binding:"required, max=40"`
	Password  string `json:"password" binding:"required"`
	EmailAddr string `json:"email_addr" binding:"required"`
}

func Register(ctx *gin.Context) {
	var ri RegisterInfo

	if err := ctx.ShouldBind(ri); err != nil {

	}

}
