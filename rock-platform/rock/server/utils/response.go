package utils

import (
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/database/models"
)

func CalcPages(total, pageSize int64) int64 {
	var pages int64
	if total > pageSize {
		if total%pageSize > 0 {
			pages = (total / pageSize) + 1
		} else {
			pages = total / pageSize
		}
	} else {
		pages = 1
	}
	return pages
}

func GetConfCtx(ctx *gin.Context) (*models.ConfCtx, error) {
	c, isExist := ctx.Get("custom_config")
	if !isExist {
		return nil, NewRockError(404, 40400001, "config info doesn't exist in cookie")
	}
	conf, ok := c.(models.ConfCtx)
	if !ok {
		return nil, NewRockError(400, 40000006, "can't unmarshal config info from cookie")
	}

	return &conf, nil
}
