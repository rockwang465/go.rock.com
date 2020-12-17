package api

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"go.rock.com/rock-platform/rock/server/database"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
)

func RegistryCreateUser(username, password, email string) (*models.User, error) {
	fmt.Println("RegistryCreateUser ")
	fmt.Printf("%s,%s,%s", username, password, email)

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
