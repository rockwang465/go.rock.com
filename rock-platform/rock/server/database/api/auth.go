package api

import (
	"go.rock.com/rock-platform/rock/server/conf"
	"go.rock.com/rock-platform/rock/server/database"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/log"
	"time"
)

type UserMap map[string]interface{}

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

// increase the number of user failed login;
// and increase the time
func CountUserLoginFailedNumber(userId int64) error {
	db := database.GetDBEngine()
	config := conf.GetConfig()
	logger := log.GetLogger()
	var user = new(models.User)
	if err := db.Where("id = ?", userId).First(user).Error; err != nil {
		return err
	}

	// if password wrong, increase login retry count
	failLoginCount := user.LoginRetryCount + 1
	var userMap = UserMap{
		"login_retry_count": failLoginCount,
	}

	// if fail login count >= server.login-retry-count, increase block time
	maxLoginRetryCount := config.Viper.GetInt64("server.login-retry-count")
	if maxLoginRetryCount == 0 {
		logger.Warning("Not define server.login-retry-count, please check")
		maxLoginRetryCount = 3
	}
	if failLoginCount >= maxLoginRetryCount {
		loginBlockDuration := config.Viper.GetDuration("server.login-block-duration")
		if loginBlockDuration == 0 {
			logger.Warning("Not define server.login-block-duration, please check")
			loginBlockDuration = time.Minute * 5
		}
		userMap["login_block_until"] = time.Now().Add(loginBlockDuration)
	}

	// update to database
	if err := db.Model(&user).Update(userMap).Error; err != nil {
		return err
	}
	return nil
}

// reset fail login retry count
func ResetRetryCount(userId int64) error {
	db := database.GetDBEngine()
	var user = new(models.User)
	if err := db.Where("id = ?", userId).First(user).Error; err != nil {
		return err
	}

	if err := db.Model(user).Update("login_retry_count", 0).Error; err != nil {
		return err
	}
	return nil
}
