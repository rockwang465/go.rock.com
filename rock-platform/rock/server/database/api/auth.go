package api

import (
	"go.rock.com/rock-platform/rock/server/database"
	"go.rock.com/rock-platform/rock/server/database/models"
)

// update user token
func UpdateUserToken(id int64, token string) (*models.User, error) {
	db := database.GetDBEngine()
	var user = new(models.User)
	if err := db.Where("id = ?", id).First(user).Error; err != nil {
		return nil, err
	}

	if err := db.Model(user).Update("token", token).Error; err != nil {
		return nil, err
	}
	return user, nil
}
