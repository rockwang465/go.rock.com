package api

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"go.rock.com/rock-platform/rock/server/database"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
	"net/http"
	"time"
)

type UserDetailResp struct {
	Id        int64            `json:"id" example:"1"`
	Name      string           `json:"name" example:"admin_user"`
	Email     string           `json:"email" example:"admin_user@sensetime.com"`
	CreatedAt models.LocalTime `json:"created_at" example:"2020-12-20 15:15:22"`
	UpdatedAt models.LocalTime `json:"updated_at" example:"2020-12-20 15:15:22"`
	//Version int64 `json:"version" example:"1"`

	RoleId          int64            `json:"role_id"`
	RoleName        string           `json:"role_name"`
	RoleDescription string           `json:"role_description" binding:"required" example:"description for role"`
	RoleCreatedAt   models.LocalTime `json:"role_created_at" example:"2020-12-20 15:15:22"`
	RoleUpdatedAt   models.LocalTime `json:"role_updated_at" example:"2020-12-20 15:15:22"`
	RoleVersion     int              `json:"role_version" binding:"required" example:"1"`
}

type UserFullResp struct {
	Id              int64            `json:"id" example:"1"`
	Name            string           `json:"name" example:"admin_user"`
	Password        string           `json:"password" example:"********"`
	Email           string           `json:"email" example:"admin_user@sensetime.com"`
	Salt            string           `json:"salt" example:"salt secret"`
	Token           string           `json:"token" example:"user token"`
	CreatedAt       models.LocalTime `json:"created_at" example:"2020-12-20 15:15:22"`
	UpdatedAt       models.LocalTime `json:"updated_at" example:"2020-12-20 15:15:22"`
	LoginBlockUntil *time.Time       `json:"login_block_until" example:"2020-12-20 15:15:22"`
	//LoginRetryCount int64            `json:"login_retry_count" example:"1"`
	//Version         int64            `json:"version" example:"1"`

	RoleId          int64            `json:"role_id"`
	RoleName        string           `json:"role_name"`
	RoleDescription string           `json:"role_description" binding:"required" example:"description for role"`
	RoleCreatedAt   models.LocalTime `json:"role_created_at" example:"2020-12-20 15:15:22"`
	RoleUpdatedAt   models.LocalTime `json:"role_updated_at" example:"2020-12-20 15:15:22"`
	RoleVersion     int              `json:"role_version" binding:"required" example:"1"`
}

func CreateUser(username, password, email string, roleId int64) (*models.User, *models.Role, error) {
	db := database.GetDBEngine()
	// get user , if exists return error
	var user = &models.User{}
	db.Where("name = ?", username).First(&user)
	if user.Id != 0 {
		err := utils.NewRockError(400, 40000001, fmt.Sprintf("user with name(%v) is alerady exist", username)) // generate a error
		return nil, nil, err
	}

	// password encrypt
	salt := utils.GenerateSalt()                   // get salt
	encryptPwd := utils.EncryptPwd(password, salt) // encrypt

	// verify role_id in role table id
	role := new(models.Role)
	if err := db.First(role, roleId).Error; err != nil {
		return nil, nil, err
	}

	var User = &models.User{
		Name:     username,
		Password: encryptPwd,
		Email:    email,
		Salt:     salt,
		RoleId:   roleId,
		//Token:           "???",
		//LoginRetryCount: 0,
		//LoginBlockUntil: nil,
	}

	// insert to table
	if err := db.Create(User).Error; err != nil {
		dbErr, ok := err.(*mysql.MySQLError)
		if ok {
			if dbErr.Number == 1062 {
				err = utils.NewRockError(400, 40000001, fmt.Sprintf("user with name(%v) is alerady exist", username))
			}
		}
		return nil, nil, err
	}
	return User, role, nil
}

func GetUsers(pageNum, pageSize int64, filedName string) (*models.UserPagination, error) {
	db := database.GetDBEngine()
	query := "%" + filedName + "%"
	Users := make([]*models.User, 0)

	var count int64
	if err := db.Order("name desc").Offset((pageNum-1)*pageSize).Where("name like ?", query).Find(&Users).Count(&count).Error; err != nil {
		fmt.Println("err: 1 count")
		return nil, err
	}

	if err := db.Order("name desc").Offset((pageNum-1)*pageSize).Where("name like ?", query).Limit(pageSize).Find(&Users).Error; err != nil {
		return nil, err
	}

	pages := utils.CalcPages(count, pageSize)
	var userPagination = &models.UserPagination{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    count,
		Pages:    pages,
		Items:    Users,
	}
	return userPagination, nil
}

func GetUserDetailResp(userId int64) (*UserDetailResp, error) {
	db := database.GetDBEngine()
	resp := new(UserDetailResp)
	if err := db.Raw("SELECT a.id as id, a.name as name, a.email as email, a.created_at as created_at, a.updated_at as updated_at, b.id as role_id, b.name as role_name, b.description as role_description, b.created_at as role_created_at, b.updated_at as role_updated_at, b.version as role_version from user a LEFT JOIN role b ON a.role_id = b.id where a.id = ? ORDER BY id ASC LIMIT 1", userId).Scan(resp).Error; err != nil {
		return nil, err
	}
	return resp, nil
}

func GetUserDetailRespByName(username string) (*UserDetailResp, error) {
	db := database.GetDBEngine()
	resp := new(UserDetailResp)
	if err := db.Raw("SELECT a.id as id, a.name as name, a.email as email, a.created_at as created_at, a.updated_at as updated_at, b.id as role_id, b.name as role_name, b.description as role_description, b.created_at as role_created_at, b.updated_at as role_updated_at, b.version as role_version from user a LEFT JOIN role b ON a.role_id = b.id where a.name = ? ORDER BY id ASC LIMIT 1", username).Scan(resp).Error; err != nil {
		return nil, err
	}
	return resp, nil
}

func GetUserFullRespByName(username string) (*UserFullResp, error) {
	db := database.GetDBEngine()
	resp := new(UserFullResp)
	if err := db.Raw("SELECT a.id as id, a.name as name, a.password as password, a.email as email, a.salt as salt, a.token as token, a.created_at as created_at, a.updated_at as updated_at, a.login_retry_count as login_retry_count, b.id as role_id, b.name as role_name, b.description as role_description, b.created_at as role_created_at, b.updated_at as role_updated_at, b.version as role_version from user a LEFT JOIN role b ON a.role_id = b.id where a.name = ? ORDER BY id ASC LIMIT 1", username).Scan(resp).Error; err != nil {
		return nil, err
	}
	return resp, nil
}

//func GetUserByName(username string) (*models.User, error) {
//	db := database.GetDBEngine()
//	var user = new(models.User)
//	if err := db.Where("name = ?", username).First(user).Error; err != nil {
//		if err.Error() == "record not found" {
//			err = utils.NewRockError(400, 40000004, fmt.Sprintf("user with name(%v) was not found", username))
//			return nil, err
//		}
//		return nil, err
//	}
//	return user, nil
//}

// get user by id, if not found return error
func HasUserById(userId int64) (*models.User, error) {
	db := database.GetDBEngine()
	var user = new(models.User)
	if err := db.Where("id = ?", userId).First(user).Error; err != nil {
		if err.Error() == "record not found" {
			err := utils.NewRockError(400, 40000005, fmt.Sprintf("User with id(%v) was not found", userId))
			return nil, err
		}
		return nil, err
	}
	return user, nil
}

// get user by name, if not found return error
func HasUserByName(username string) (*models.User, error) {
	db := database.GetDBEngine()
	var user = new(models.User)
	if err := db.Where("name = ?", username).First(user).Error; err != nil {
		if err.Error() == "record not found" {
			err = utils.NewRockError(400, 40000004, fmt.Sprintf("User with name(%v) was not found", username))
			return nil, err
		}
		return nil, err
	}
	return user, nil
}

// delete user by id, if not found return error
func DeleteUserById(userId int64) (string, error) {
	db := database.GetDBEngine()
	var user = new(models.User)

	user, err := HasUserById(userId)
	if err != nil {
		return "", err
	}

	if user.Name == "admin" {
		err = utils.NewRockError(http.StatusBadRequest, 40400002, fmt.Sprintf("Admin user can't be deleted"))
		return "", err
	}

	username := user.Name

	if err := db.Where("id = ?", userId).Delete(user).Error; err != nil {
		return "", err
	}
	return username, nil
}

// update user password by id, if not found return error
func UpdateUserPwdById(id int64, oldPwd, newPwd, role string) (*models.User, error) {
	user, err := HasUserById(id)
	if err != nil {
		return nil, err
	}

	if oldPwd == newPwd {
		err = utils.NewRockError(http.StatusBadRequest, 40000008, "The new password cannot be the same as the old one")
		return nil, err
	}
	err = utils.CheckPwd(newPwd)
	if err != nil {
		return nil, err
	}

	encOldPwd := utils.EncryptPwd(oldPwd, user.Salt)
	if encOldPwd != user.Password {
		err = utils.NewRockError(http.StatusBadRequest, 40000003, "password incorrect")
		return nil, err
	}

	encNewPwd := utils.EncryptPwd(newPwd, user.Salt)
	token, err := utils.GenerateToken(id, user.Name, encNewPwd, role)
	if err != nil {
		return nil, err
	}

	db := database.GetDBEngine()
	if err = db.Model(&user).Update(map[string]interface{}{"password": encNewPwd, "token": token}).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Get user by email
func GetUserByEmail(email string) (*models.User, error) {
	db := database.GetDBEngine()
	user := new(models.User)
	if err := db.Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Reset user secret by id
func ResetSecretWithId(id int64, secret string, expire time.Duration) (*models.User, error) {
	db := database.GetDBEngine()
	user := new(models.User)
	if err := db.First(user, id).Error; err != nil {
		return nil, err
	}

	secretExpireAt := time.Now().Add(expire)
	if err := db.Model(user).Update(map[string]interface{}{"reset_secret": secret, "secret_expired_at": secretExpireAt}).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Has email by email(true: has, false: not found)
func HasEmail(email string) (bool, error) {
	db := database.GetDBEngine()
	user := new(models.User)
	err := db.Where("email = ?", email).First(user).Error
	if err != nil {
		if err.Error() == "record not found" {
			return false, nil // not found email
		}
		return false, err
	}
	//fmt.Printf("user id:%v, email: %v\n", user.Id, user.Email)
	return true, nil
}
