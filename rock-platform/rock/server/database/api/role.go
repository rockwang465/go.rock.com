package api

import (
	"fmt"
	"go.rock.com/rock-platform/rock/server/database"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
)

type RoleBriefResp struct {
	Id          int64            `json:"id" binding:"required" example:"1"`
	Name        string           `json:"name" binding:"required" example:"admin_role"`
	Description string           `json:"description" binding:"required" example:"description for role"`
	CreatedAt   models.LocalTime `json:"created_at" binding:"required" example:"2018-10-09T14:57:23+08:00"`
	UpdatedAt   models.LocalTime `json:"updated_at" binding:"required" example:"2018-10-09T14:57:23+08:00"`
	Version     int              `json:"version" binding:"required" example:"1"`
}

func GetRoleByName(name string) (*models.Role, error) {
	role := new(models.Role)
	db := database.GetDBEngine()
	if err := db.Where("name = ?", name).Find(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func CreateRole(name, description string) error {
	db := database.GetDBEngine()
	role := &models.Role{
		Name:        name,
		Description: description,
	}

	if err := db.Create(role).Error; err != nil {
		return err
	}
	return nil
}

func GetBriefRoleByName(name string) (*models.Role, error) {
	db := database.GetDBEngine()
	role := new(models.Role)
	if err := db.Where("name = ?", name).First(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func GetRoles(pageNum, pageSize int64, filedName string) (*models.RolePagination, error) {
	db := database.GetDBEngine()
	query := "%" + filedName + "%"
	Roles := make([]*models.Role, 0)

	var count int64
	if err := db.Order("updated_at desc").Offset((pageNum-1)*pageSize).Where("name like ?", query).Find(&Roles).Count(&count).Error; err != nil {
		return nil, err
	}

	if err := db.Order("updated_at desc").Offset((pageNum-1)*pageSize).Where("name like ?", query).Limit(pageSize).Find(&Roles).Error; err != nil {
		return nil, err
	}

	pages := utils.CalcPages(count, pageSize)
	var rolePagination = &models.RolePagination{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    count,
		Pages:    pages,
		Items:    Roles,
	}
	return rolePagination, nil
}

func GetRoleById(id int64) (*models.Role, error) {
	db := database.GetDBEngine()
	role := new(models.Role)
	if err := db.First(role, id).Error; err != nil {
		if err.Error() == "record not found" {
			err = utils.NewRockError(404, 40400006, fmt.Sprintf("Role with id(%v) was not found", id))
			return nil, err
		}
		return nil, err
	}
	return role, nil
}

func DeleteRoleById(id int64) error {
	db := database.GetDBEngine()
	role, err := GetRoleById(id)
	if err != nil {
		return err
	}

	if err := db.Delete(role).Error; err != nil {
		return err
	}
	return nil
}

func UpdateRole(id int64, desc string) (*models.Role, error) {
	db := database.GetDBEngine()
	role, err := GetRoleById(id)
	if err != nil {
		return nil, err
	}
	if err := db.Model(role).Update(map[string]interface{}{"description": desc}).Error; err != nil {
		return nil, err
	}
	return role, nil
}

// get all user by role id, no query field
func GetRoleUsers(roleId, pageNum, pageSize int64) (*models.UserPagination, error) {
	db := database.GetDBEngine()
	Users := make([]*models.User, 0)

	var count int64
	if err := db.Order("updated_at desc").
		Offset((pageNum-1)*pageSize).
		Where("role_id = ?", roleId).
		Find(&Users).
		Count(&count).Error; err != nil {
		return nil, err
	}

	if err := db.Order("updated_at desc").
		Offset((pageNum-1)*pageSize).
		Where("role_id = ?", roleId).
		Limit(pageSize).
		Find(&Users).Error; err != nil {
		return nil, err
	}

	userPagination := &models.UserPagination{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    count,
		Pages:    utils.CalcPages(count, pageSize),
		Items:    Users,
	}
	return userPagination, nil
}
