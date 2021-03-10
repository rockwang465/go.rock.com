package api

import (
	"go.rock.com/rock-platform/rock/server/database"
	"go.rock.com/rock-platform/rock/server/database/models"
)

// update or create app_conf by app_id and project_env_id and app config
func UpdateOrCreateAppConfById(appId, projectEnvId int64, config string) (*models.AppConf, error) {
	_, err := GetAppById(appId)
	if err != nil {
		return nil, err
	}
	_, err = GetProjectEnvById(projectEnvId)
	if err != nil {
		return nil, err
	}

	db := database.GetDBEngine()
	appConf := new(models.AppConf)
	if err := db.Where("app_id = ? AND project_env_id = ?", appId, projectEnvId).First(appConf).Error; err != nil {
		if err.Error() != "record not found" {
			return nil, err
		}
	}

	if appConf.Id == 0 { // if no record, then create it
		appConf.AppId = appId
		appConf.ProjectEnvId = projectEnvId
		appConf.Config = config
		if err := db.Create(appConf).Error; err != nil {
			return nil, err
		}
	} else { // if has record, then update it
		if err := db.Model(appConf).Update(map[string]interface{}{"config": config}).Error; err != nil {
			return nil, err
		}
	}
	return appConf, nil
}

func GetAppConfByAppAndProjectEnvId(appId, projectEnvId int64) (*models.AppConf, error) {
	_, err := GetAppById(appId)
	if err != nil {
		return nil, err
	}
	_, err = GetProjectEnvById(projectEnvId)
	if err != nil {
		return nil, err
	}

	db := database.GetDBEngine()
	appConf := new(models.AppConf)
	if err := db.Where("app_id = ? AND project_env_id = ?", appId, projectEnvId).First(appConf).Error; err != nil {
		return nil, err
	}
	return appConf, nil
}

func DeleteAppConfByProjectAndAppId(appId, projectEnvId int64) error {
	_, err := GetAppById(appId)
	if err != nil {
		return err
	}
	_, err = GetProjectEnvById(projectEnvId)
	if err != nil {
		return err
	}

	db := database.GetDBEngine()
	appConf := models.AppConf{
		AppId:        appId,
		ProjectEnvId: projectEnvId,
	}
	if err := db.Delete(appConf).Error; err != nil {
		return err
	}
	return nil
}

//func GetAppConfById(id int64) (*models.AppConf, error) {
//	db := database.GetDBEngine()
//	appConf := new(models.AppConf)
//	if err := db.First(appConf, id).Error; err != nil {
//		if err.Error() != "record not found" {
//			return nil, err
//		}
//	}
//	return appConf, nil
//}
