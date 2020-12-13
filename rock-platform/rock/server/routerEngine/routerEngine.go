package routerEngine

import "github.com/gin-gonic/gin"

type Routers struct {
	*gin.Engine
}

func GetRouterEngine() *Routers {
	return &Routers{
		gin.New(),
	}
}
