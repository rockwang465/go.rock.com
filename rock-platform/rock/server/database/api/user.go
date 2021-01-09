package api

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"go.rock.com/rock-platform/rock/server/database"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
)

//type UserBriefResp struct {
//	Id        int64            `json:"id" example:"1"`
//	Name      string           `json:"name" example:"admin_user"`
//	Email     string           `json:"email" example:"admin_user@sensetime.com"`
//	CreatedAt models.LocalTime `json:"created_at" example:"2020-12-20 15:15:22"`
//	UpdatedAt models.LocalTime `json:"updated_at" example:"2020-12-20 15:15:22"`
//	//Version   int              `json:"version" binding:"required" example:"1"`
//}

//type UserDetailResp struct { // user信息是这样的:{"id":2,"name":"rock1", ..., "role":{"id":0,...}}, 所以要做嵌套结构体
//	//UserBriefResp
//	Id        int64            `json:"id" example:"1"`
//	Name      string           `json:"name" example:"admin_user"`
//	Email     string           `json:"email" example:"admin_user@sensetime.com"`
//	CreatedAt models.LocalTime `json:"created_at" example:"2020-12-20 15:15:22"`
//	UpdatedAt models.LocalTime `json:"updated_at" example:"2020-12-20 15:15:22"`
//	//Role      *RoleBriefResp   `json:"role"` // xorm支持基于 UserBriefResp 查 RoleBriefResp 关联信息，gorm不支持；顺义建议gorm用join功能查 RoleBriefResp 信息，这样性能比xorm更好。
//}

type ResultResp struct {
	Id        int64            `json:"id" example:"1"`
	Name      string           `json:"name" example:"admin_user"`
	Email     string           `json:"email" example:"admin_user@sensetime.com"`
	CreatedAt models.LocalTime `json:"created_at" example:"2020-12-20 15:15:22"`
	UpdatedAt models.LocalTime `json:"updated_at" example:"2020-12-20 15:15:22"`
	RoleId    int64            `json:"role_id" example:"1"`
}

func hasUserById(db *database.DBEngine, userId int64) (*models.User, error) {
	var user = &models.User{}
	db.First(user, userId)
	if user.Id == 0 {
		err := utils.NewRockError(404, 40400001, fmt.Sprintf("User with id %v was not found", userId))
		return nil, err
	}
	return user, nil
}

func hasUserByName(db *database.DBEngine, username string) (*models.User, error) {
	var user = &models.User{}
	db.Where("username = ?", username).First(user)
	if user.Id == 0 {
		err := utils.NewRockError(404, 40400002, fmt.Sprintf("User with username %v was not found", username))
		return nil, err
	}
	return user, nil
}

func CreateUser(username, password, email string, roleId int64) (*models.User, error) {
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
		RoleId:   roleId,
		//Token:           "???",
		//LoginRetryCount: 0,
		//LoginBlockUntil: nil,
	}
	fmt.Println("User model info:", User.Name, User.Password, User.Email, User.RoleId)

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

	//// 这里应该是创建用户，不应该是select查询
	//resp := new(ResultResp)
	//if err := db.Raw("SELECT a.id, a.name, a.email, a.created_at, a.updated_at from user a LEFT JOIN role b ON a.role_id = b.id").Scan(resp).Error; err != nil {
	//	return nil, err
	//}
	//fmt.Printf("resp UserDetailResp------------>")
	//fmt.Printf("%#v\n", resp)
	//fmt.Println(resp.Id, resp.Name, resp.Email, resp.CreatedAt, resp.UpdatedAt,resp.RoleId)
	////fmt.Println(resp.Role)

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
