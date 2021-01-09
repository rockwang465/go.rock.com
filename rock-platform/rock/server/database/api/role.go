package api

import (
	"go.rock.com/rock-platform/rock/server/database/models"
)

type RoleBriefResp struct {
	Id          int64            `json:"id" binding:"required" example:"1"`
	Name        string           `json:"name" binding:"required" example:"admin_role"`
	Description string           `json:"description" binding:"required" example:"description for role"`
	CreatedAt   models.LocalTime `json:"created_at" binding:"required" example:"2018-10-09T14:57:23+08:00"`
	UpdatedAt   models.LocalTime `json:"updated_at" binding:"required" example:"2018-10-09T14:57:23+08:00"`
	Version     int              `json:"version" binding:"required" example:"1"`
}
