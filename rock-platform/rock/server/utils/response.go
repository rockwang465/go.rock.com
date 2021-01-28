package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.rock.com/rock-platform/rock/server/database/models"
	"net/http"
	"regexp"
)

const (
	LevelD = iota
	LevelC
	LevelB
	LevelA
	LevelS

	PwdMinLength = 6
	PwdMaxLength = 20
	PwdMinLevel  = 2
)

// check password
func CheckPwd(password string) error {
	if len(password) < PwdMinLength {
		err := NewRockError(http.StatusBadRequest, 42200002, fmt.Sprintf("The password length is too short, greater than or equal 6")) // generate a error
		return err
	}
	if len(password) > PwdMaxLength {
		err := NewRockError(http.StatusBadRequest, 40000002, fmt.Sprintf("The password length is too long, less than or equal 20")) // generate a error
		return err
	}

	var level int = LevelD
	patternList := []string{`[0-9]+`, `[A-Z]+`, `[a-z]+`, `[~!@#$%^&*?_-]+`}
	for _, pattern := range patternList {
		match, err := regexp.MatchString(pattern, password)
		if err != nil {
			return err
		}
		if match {
			level += 1
		}
	}

	if level < PwdMinLevel {
		err := NewRockError(http.StatusBadRequest, 40000007, "Password not strong")
		return err
	}
	return nil
}

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

func MarshalResponse(src, dest interface{}) error {
	byteSrc, err := json.Marshal(src)
	if err != nil {
		return err
	}
	fmt.Println(string(byteSrc))

	if err := json.Unmarshal(byteSrc, dest); err != nil {
		return err
	}
	fmt.Println(dest)
	return nil
}
