package v1

import (
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/client/k8s"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
)

type CreateClusterReq struct {
	Name        string `json:"name" binding:"required" example:"test-cluster"`
	Description string `json:"description" binding:"omitempty,max=100" example:"description for a k8s cluster"`
	Config      string `json:"config" binding:"required" example:"k8s config file"`
}

type ClusterBriefResp struct {
	Id          int64            `json:"id" example:"1"`
	Name        string           `json:"name" example:"test-cluster"`
	Description string           `json:"description" example:"description for a k8s cluster"`
	CreatedAt   models.LocalTime `json:"created_at" example:"2021-01-28 20:20:20"`
	UpdateAt    models.LocalTime `json:"updated_at" example:"2021-01-28 20:20:20"`
}

type PaginationClusterResp struct {
	PageNum int64               `json:"page_num" binding:"required" example:"1"`
	PerSize int64               `json:"per_size" binding:"required" example:"10"`
	Total   int64               `json:"total" binding:"required" example:"100"`
	Pages   int64               `json:"pages" binding:"required" example:"1"`
	Items   []*ClusterBriefResp `json:"items" binding:"required"`
}

type UpdateClusterReq struct {
	Description string `json:"description" binding:"omitempty,max=100" example:"description the cluster"`
	Config      string `json:"config" binding:"required" example:"k8s config file"`
}

// @Summary Create cluster
// @Description Api to create k8s cluster
// @Tags CLUSTER
// @Accept json
// @Produce json
// @Param input_body body v1.CreateClusterReq true "JSON type input body"
// @Success 201 {object} v1.ClusterBriefResp "StatusCreated"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/clusters [post]
func (c *Controller) CreateCluster(ctx *gin.Context) {
	// postman存放/etc/kubernetes/admin.conf,需要先从yaml转为字符串方式，操作方法如下:
	// sed s/$/"\\\n"/ /etc/kubernetes/admin.conf | tr -d "\n"
	var req CreateClusterReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}

	// check k8s status
	if err := k8s.IsK8sHealth(req.Config); err != nil {
		c.Logger.Error("use k8s config request failed, ", err)
		err = utils.NewRockError(400, 40000019, "the k8s config is not correct, please check it")
		panic(err)
	}

	cluster, err := api.CreateCluster(req.Name, req.Description, req.Config)
	if err != nil {
		panic(err)
	}
	resp := ClusterBriefResp{}
	if err := utils.MarshalResponse(cluster, &resp); err != nil {
		panic(err)
	}
	c.Logger.Infof("Create cluster by id:%v name:%v", resp.Id, resp.Name)
	ctx.JSON(http.StatusCreated, resp)
}

// @Summary Get all clusters
// @Description Api to get all k8s clusters
// @Tags CLUSTER
// @Accept json
// @Produce json
// @Param page_num query integer true "Request page number" default(1)
// @Param per_size query integer true "Cluster number size" default(10)
// @Param query_field query string false "Fuzzy Query(field: name)"
// @Success 200 {object} v1.PaginationClusterResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/clusters [get]
func (c *Controller) GetClusters(ctx *gin.Context) {
	var paginationReq GetPaginationReq
	if err := ctx.ShouldBind(&paginationReq); err != nil {
		panic(err)
	}

	clusterPg, err := api.GetClusters(paginationReq.PageNum, paginationReq.PageSize, paginationReq.QueryField)
	if err != nil {
		panic(err)
	}

	resp := PaginationClusterResp{}
	if err := utils.MarshalResponse(clusterPg, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get all clusters, this pagination cluster number is: %v", len(resp.Items))
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get a cluster by id
// @Description Api to get a cluster by id
// @Tags CLUSTER
// @Accept json
// @Produce json
// @Param id query integer true "Cluster ID"
// @Success 200 {object} v1.ClusterBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/clusters/{id} [get]
func (c *Controller) GetCluster(ctx *gin.Context) {
	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	cluster, err := api.GetClusterById(idReq.Id)
	if err != nil {
		panic(err)
	}

	resp := ClusterBriefResp{}
	if err := utils.MarshalResponse(cluster, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get cluster by id:%v", resp.Id)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Delete a cluster by id
// @Description Api to delete a cluster by id
// @Tags CLUSTER
// @Accept json
// @Produce json
// @Param id body integer true "Cluster ID"
// @Success 204 {object} string "StatusNoContent"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/clusters/{id} [delete]
func (c *Controller) DeleteCluster(ctx *gin.Context) {
	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	if err := api.DeleteClusterById(idReq.Id); err != nil {
		panic(err)
	}
	c.Logger.Infof("Delete cluster by id:%v", idReq.Id)
	ctx.JSON(http.StatusNoContent, "")
}

// @Summary Update cluster description by id and body
// @Description api for update cluster description
// @Tags CLUSTER
// @Accept json
// @Produce json
// @Param id path integer true "Cluster ID"
// @Param update_body body v1.UpdateClusterReq true "JSON type for update cluster description"
// @Success 200 {object} v1.ClusterBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/clusters/{id} [put]
func (c *Controller) UpdateCluster(ctx *gin.Context) {
	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	var req UpdateClusterReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}

	// check k8s status
	if err := k8s.IsK8sHealth(req.Config); err != nil {
		c.Logger.Error("use k8s config request failed, ", err)
		err = utils.NewRockError(400, 40000019, "the k8s config is not correct, please check it")
	}

	cluster, err := api.UpdateCluster(idReq.Id, req.Description, req.Config)
	if err != nil {
		panic(err)
	}

	resp := ClusterBriefResp{}
	if err := utils.MarshalResponse(cluster, &resp); err != nil {
		panic(err)
	}
	c.Logger.Infof("Update cluster's config by id:%v", idReq.Id)
	ctx.JSON(http.StatusOK, resp)
}
