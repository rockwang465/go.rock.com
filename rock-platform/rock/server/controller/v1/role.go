package v1

import (
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/database/api"
	"net/http"
)

//type RoleIdReq struct {
//	Id int64 `json:"id" binding:"required,min=1" example:"1"`
//}

type CreateRoleReq struct {
	Name        string `json:"name" binding:"required" example:"admin_role"`
	Description string `json:"description" binding:"omitempty,max=100" example:"description for role"`
}

type QueryReq struct {
	QueryField string `json:"query_field" binding:"omitempty" example:"admin"`
}

type UpdateRoleReq struct {
	Description string `json:"description" binding:"omitempty,max=100" example:"description for role"`
}

//type RoleBriefResp struct {
//	Id          int64      `json:"id" binding:"required" example:"1"`
//	Name        string     `json:"name" binding:"required" example:"admin_role"`
//	Description string     `json:"description" binding:"required" example:"description for role"`
//	CreatedAt   *time.Time `json:"created_at" binding:"required" example:"2018-10-09T14:57:23+08:00"`
//	UpdatedAt   *time.Time `json:"updated_at" binding:"required" example:"2018-10-09T14:57:23+08:00"`
//	Version     int        `json:"version" binding:"required" example:"1"`
//}

// @Summary Create role
// @Description api for create role
// @Tags ROLE
// @Accept json
// @Produce json
// @Param input_body body v1.CreateRoleReq true "JSON type input body"
// @Success 201 {object} models.Role "StatusCreated"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/roles [post]
func (c *Controller) CreateRole(ctx *gin.Context) {
	var roleReq CreateRoleReq
	if err := ctx.ShouldBindJSON(&roleReq); err != nil {
		panic(err)
	}

	_, err := api.GetRoleByName(roleReq.Name)
	if err != nil {
		c.Logger.Info("GetRoleByName err.Error: ", err.Error())
		if err.Error() != "record not found" {
			panic(err)
		}
	}

	if err = api.CreateRole(roleReq.Name, roleReq.Description); err != nil {
		panic(err)
	}

	role, err := api.GetBriefRoleByName(roleReq.Name)
	if err != nil {
		panic(err)
	}
	c.Logger.Infof("Create role with id:%v name:%v", role.Id, role.Name)
	ctx.JSON(http.StatusCreated, role)
}

// @Summary Get role
// @Description api for get all roles
// @Tags ROLE
// @Accept json
// @Produce json
// @Param page_num query integer true "Request page number" default(1)
// @Param page_size query integer true "Request page size" default(10)
// @Param QueryField query string false "Fuzzy Query(field: name)"
// @Success 200 {object} models.RolePagination "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/roles [get]
func (c *Controller) GetRoles(ctx *gin.Context) {
	var paginationReq GetPaginationReq
	err := ctx.ShouldBind(&paginationReq)
	if err != nil {
		panic(err)
	}
	rolePg, err := api.GetRoles(paginationReq.PageNum, paginationReq.PageSize, paginationReq.QueryField)
	if err != nil {
		panic(err)
	}
	c.Logger.Infof("Get all roles, this pagination role number is: %v", len(rolePg.Items))
	ctx.JSON(http.StatusOK, rolePg)
}

// @Summary Get an role by id
// @Description api for get an role by id
// @Tags ROLE
// @Accept json
// @Produce json
// @Param id path integer true "Role ID"
// @Success 200 {object} models.Role "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/roles/{id} [get]
func (c *Controller) GetRole(ctx *gin.Context) {
	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	role, err := api.GetRoleById(idReq.Id)
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, role)
}

// @Summary Get an role by id
// @Description api for get an role by id
// @Tags ROLE
// @Accept json
// @Produce json
// @Param id path integer true "Role ID"
// @Success 204 {object} string "StatusNoContent"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/roles/{id} [delete]
func (c *Controller) DeleteRole(ctx *gin.Context) {
	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}
	if err := api.DeleteRoleById(idReq.Id); err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusNoContent, "")
}

// @Summary Update role description by id and body
// @Description api for update role description
// @Tags ROLE
// @Accept json
// @Produce json
// @Param id path integer true "Role ID"
// @Param description body string true "Role Description"
// @Param update_body body v1.UpdateRoleReq true "JSON type for update role description"
// @Success 200 {object} models.Role "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/roles/{id} [put]
func (c *Controller) UpdateRole(ctx *gin.Context) {
	var idReq IdReq
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	var descReq UpdateRoleReq
	if err := ctx.ShouldBindJSON(&descReq); err != nil {
		panic(err)
	}
	role, err := api.UpdateRole(idReq.Id, descReq.Description)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, role)
}
