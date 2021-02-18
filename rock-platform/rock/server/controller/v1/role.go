package v1

import (
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/database/api"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
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

type RoleBriefResp struct {
	Id          int64            `json:"id" binding:"required" example:"1"`
	Name        string           `json:"name" binding:"required" example:"admin_role"`
	Description string           `json:"description" binding:"required" example:"description for role"`
	CreatedAt   models.LocalTime `json:"created_at" binding:"required" example:"2018-10-09T14:57:23+08:00"`
	UpdatedAt   models.LocalTime `json:"updated_at" binding:"required" example:"2018-10-09T14:57:23+08:00"`
	Version     int              `json:"version" binding:"required" example:"1"`
}

type PaginateRoleResp struct {
	PageNum  int64            `json:"page_num" binding:"required" example:"1"`
	PageSize int64            `json:"page_size" binding:"required" example:"10"`
	Total    int64            `json:"total" binding:"required" example:"100"`
	Pages    int64            `json:"pages" binding:"required" example:"1"`
	Items    []*RoleBriefResp `json:"items" binding:"required"`
}

// @Summary Create role
// @Description api for create role
// @Tags ROLE
// @Accept json
// @Produce json
// @Param input_body body v1.CreateRoleReq true "JSON type input body"
// @Success 201 {object} v1.RoleBriefResp "StatusCreated"
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

	resp := RoleBriefResp{}
	if err := utils.MarshalResponse(role, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Create role by id:%v name:%v", role.Id, role.Name)
	ctx.JSON(http.StatusCreated, resp)
}

// @Summary Get role
// @Description api for get all roles
// @Tags ROLE
// @Accept json
// @Produce json
// @Param page_num query integer true "Request page number" default(1)
// @Param page_size query integer true "Request page size" default(10)
// @Param QueryField query string false "Fuzzy Query(field: name)"
// @Success 200 {object} v1.PaginateRoleResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/roles [get]
func (c *Controller) GetRoles(ctx *gin.Context) {
	var paginationReq GetPaginationReq
	if err := ctx.ShouldBind(&paginationReq); err != nil {
		panic(err)
	}
	rolePg, err := api.GetRoles(paginationReq.PageNum, paginationReq.PageSize, paginationReq.QueryField)
	if err != nil {
		panic(err)
	}
	resp := PaginateRoleResp{}
	if err := utils.MarshalResponse(rolePg, &resp); err != nil {
		panic(err)
	}
	c.Logger.Infof("Get all roles, this pagination role number is: %v", len(rolePg.Items))
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get an role by id
// @Description api for get an role by id
// @Tags ROLE
// @Accept json
// @Produce json
// @Param id path integer true "Role ID"
// @Success 200 {object} v1.RoleBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
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

	resp := RoleBriefResp{}
	if err := utils.MarshalResponse(role, &resp); err != nil {
		panic(err)
	}
	c.Logger.Infof("Get role name:%v by id:%v", role.Name, role.Id)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Delete an role by id
// @Description api for delete an role by id
// @Tags ROLE
// @Accept json
// @Produce json
// @Param id path integer true "Role ID"
// @Success 204 {object} string "StatusNoContent"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
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
	c.Logger.Infof("Delete role by id:%v", idReq.Id)
	ctx.JSON(http.StatusNoContent, "")
}

// @Summary Update role description by id and body
// @Description api for update role description
// @Tags ROLE
// @Accept json
// @Produce json
// @Param id path integer true "Role ID"
// @Param update_body body v1.UpdateRoleReq true "JSON type for update role description"
// @Success 200 {object} v1.RoleBriefResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
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

	resp := RoleBriefResp{}
	if err := utils.MarshalResponse(role, &resp); err != nil {
		panic(err)
	}
	c.Logger.Infof("Update role's description by id:%v", idReq.Id)
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get all users by species role id
// @Description api for get all users by species role id
// @Tags ROLE
// @Accept json
// @Produce json
// @Param id path integer true "Role ID"
// @Param page_num query integer true "Request page number" default(1)
// @Param page_size query integer true "Request page size" default(10)
// @Success 200 {object} v1.PaginateBriefUserResp "StatusOK"
// @Failure 400 {object} utils.HTTPError "StatusBadRequest"
// @Failure 404 {object} utils.HTTPError "StatusNotFound"
// @Failure 500 {object} utils.HTTPError "StatusInternalServerError"
// @Router /v1/roles/{id}/users [get]
func (c *Controller) GetRoleUsers(ctx *gin.Context) {
	var idReq IdReq // role id
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		panic(err)
	}

	var paginationReq GetPaginationReq
	if err := ctx.ShouldBind(&paginationReq); err != nil {
		panic(err)
	}

	_, err := api.GetRoleById(idReq.Id)
	if err != nil {
		panic(err)
	}
	// get users by role id, no query field
	userPg, err := api.GetRoleUsers(idReq.Id, paginationReq.PageNum, paginationReq.PageSize)
	if err != nil {
		panic(err)
	}

	resp := PaginateBriefUserResp{}
	if err := utils.MarshalResponse(userPg, &resp); err != nil {
		panic(err)
	}

	c.Logger.Infof("Get all users, this pagination role number is: %v", len(userPg.Items))
	ctx.JSON(http.StatusOK, resp)
}
