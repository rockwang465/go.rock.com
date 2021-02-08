package v1

// record namespace and cluster_id for a k8s cluster

//import "github.com/gin-gonic/gin"

type CreateEnvReq struct {
	Namespace   string `json:"namespace" binding:"required" example:"namespace of k8s cluster"`
	Description string `json:"description" binding:"omitempty,max=250" example:"description for env"`
	ClusterId   int64  `json:"cluster_id" binding:"required"`
}

//func (c *Controller) CreateEnv(ctx *gin.Context) {
//	var req CreateEnvReq
//}
