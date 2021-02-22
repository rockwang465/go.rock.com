package api

import (
	"fmt"
	"go.rock.com/rock-platform/rock/server/database"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
)

func CreateApp(name, fullName, owner, desc, gitlabAddr string, projectId, gitlabProjectId, droneRepoId int64) (*models.App, error) {
	if err := hasNotAppWithSameNameAndProject(name, projectId); err != nil {
		return nil, err
	}

	db := database.GetDBEngine()
	app := &models.App{
		Name:            name,
		FullName:        fullName,
		Owner:           owner,
		Description:     desc,
		GitlabAddress:   gitlabAddr,
		ProjectId:       projectId,
		DroneRepoId:     droneRepoId,
		GitlabProjectId: gitlabProjectId,
	}
	if err := db.Create(app).Error; err != nil {
		return nil, err
	}
	return app, nil
}

// get all apps
func GetApps(pageNum, pageSize int64, queryField string, projectId int64) (*models.AppPagination, error) {
	db := database.GetDBEngine()
	_, err := GetProjectById(projectId)
	if err != nil {
		return nil, err
	}

	query := "%" + queryField + "%"
	Apps := make([]*models.App, 0)

	var count int64
	if err := db.Order("updated_at desc").
		Offset((pageNum-1)*pageSize).
		Where("project_id = ? AND name like ?", projectId, query).
		Find(&Apps).
		Count(&count).Error; err != nil {
		return nil, err
	}

	if err := db.Order("updated_at desc").
		Offset((pageNum-1)*pageSize).
		Where("project_id = ? AND name like ?", projectId, query).
		Limit(pageSize).
		Find(&Apps).Error; err != nil {
		return nil, err
	}

	appPagination := &models.AppPagination{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    count,
		Pages:    utils.CalcPages(count, pageSize),
		Items:    Apps,
	}
	return appPagination, nil
}

func GetAppById(id int64) (*models.App, error) {
	db := database.GetDBEngine()
	app := new(models.App)
	if err := db.First(app, id).Error; err != nil {
		if err.Error() == "record not found" {
			err = utils.NewRockError(404, 40400003, fmt.Sprintf("App with id(%v) is not found", id))
			return nil, err
		}
		return nil, err
	}
	return app, nil
}

func DeleteAppById(id int64) error {
	db := database.GetDBEngine()
	app, err := GetAppById(id)
	if err != nil {
		return err
	}

	if err := db.Delete(app).Error; err != nil {
		return err
	}
	return nil
}

func UpdateApp(id int64, desc string) (*models.App, error) {
	db := database.GetDBEngine()
	app, err := GetAppById(id)
	if err != nil {
		return nil, err
	}
	if err := db.Model(app).Update(map[string]interface{}{"description": desc}).Error; err != nil {
		return nil, err
	}
	return app, nil
}

// get all apps
func GetAppsList() ([]*models.App, error) {
	db := database.GetDBEngine()
	apps := make([]*models.App, 0)
	if err := db.Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

// update an app gitlabAddress, fullName, owner, gitlabProjectId, droneRepoId
func UpdateAppGitlabAddressById(appId int64, fullName, owner, gitlabAddress string, gitlabProjectId, droneRepoId int64) (*models.App, error) {
	app, err := GetAppById(appId)
	if err != nil {
		return nil, err
	}

	db := database.GetDBEngine()
	if err := db.Model(app).Update(map[string]interface{}{"full_name": fullName, "owner": owner, "gitlab_project_id": gitlabProjectId, "gitlab_address": gitlabAddress, "drone_repo_id": droneRepoId}).Error; err != nil {
		return nil, err
	}
	return app, nil
}

// ensure not same name app in same projectId
func hasNotAppWithSameNameAndProject(name string, projectId int64) error {
	db := database.GetDBEngine()
	// ensure has project id
	_, err := GetProjectById(projectId)
	if err != nil {
		return err
	}

	// ensure same projectId not same name app
	app := new(models.App)
	if err := db.Where("project_id = ? AND name = ?", projectId, name).First(&app).Error; err != nil {
		if err.Error() == "record not found" {
			return nil
		}
		return err
	}
	err = utils.NewRockError(400, 40000017, fmt.Sprintf("App with name(%v) in project_id(%v) already exists", name, projectId))
	return err
}
