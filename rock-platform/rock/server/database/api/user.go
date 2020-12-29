package api

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"go.rock.com/rock-platform/rock/server/database"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
)

func CreateUser(username, password, email string) (*models.User, error) {
	//fmt.Printf("%s,%s,%s", username, password, email)
	db := database.GetDBEngine()
	// get user , if exists return error
	var user = &models.User{}
	db.Where("name = ?", username).First(&user)
	if user.Id != 0 {
		err := utils.NewRockError(400, 40000001, fmt.Sprintf("user with name(%v) is alerady exist", username)) // generate a error
		return nil, err
	}

	// password encrypt
	salt := utils.GenerateSalt()                   // get salt
	encryptPwd := utils.EncryptPwd(password, salt) // encrypt
	var User = &models.User{
		Name:     username,
		Password: encryptPwd,
		Email:    email,
		Salt:     salt,
		//Token:           "???",
		//LoginRetryCount: 0,
		//LoginBlockUntil: nil,
		//RoleId:          nil,
	}

	// insert to table
	if err := db.Create(User).Error; err != nil {
		dbErr, ok := err.(*mysql.MySQLError)
		if ok {
			if dbErr.Number == 1062 {
				err = utils.NewRockError(400, 40000001, fmt.Sprintf("user with name(%v) is alerady exist", username))
			}
		}
		return nil, err
	}
	return User, nil
}

func GetUserByName(username string) (*models.User, error) {
	db := database.GetDBEngine()
	var user = new(models.User)
	if err := db.Where("name = ?", username).First(user).Error; err != nil {
		if err.Error() == "record not found" {
			err = utils.NewRockError(400, 40000004, fmt.Sprintf("user with name(%v) is not found", username))
			return nil, err
		}
		return nil, err
	}
	return user, nil
}
