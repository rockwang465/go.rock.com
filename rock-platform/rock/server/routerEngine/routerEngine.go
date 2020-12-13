package routerEngine

import "github.com/gin-gonic/gin"

type Routers struct {
	*gin.Engine
}

var SingleRouters *Routers

func GetRouterEngine() *Routers {
	if SingleRouters == nil {
		SingleRouters = &Routers{
			gin.New(),
		}
	}
	return SingleRouters
}
