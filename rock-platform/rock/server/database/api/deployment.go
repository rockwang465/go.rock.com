package api

import (
	"fmt"
	"go.rock.com/rock-platform/rock/server/database"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
)

// insert the current deployment info into the database
func CreateDeployment(appId, envId int64, chartName, chartVersion, description, namespace string) (*models.Deployment, error) {
	deployment := models.Deployment{
		Description:  description,
		ChartName:    chartName,
		ChartVersion: chartVersion,
		AppId:        appId,
		EnvId:        envId,
	}

	deployment.Name = utils.GenerateChartName(chartName, namespace)

	// 由于 app_id env_id 在调用此函数前已经做了检测了，所以这里就不再做检测了。
	db := database.GetDBEngine()
	if err := db.Create(&deployment).Error; err != nil {
		return nil, err
	}
	return &deployment, nil
}

func GetDeployments(pageNum, pageSize int64) (*models.DeploymentPagination, error) {
	db := database.GetDBEngine()
	Deployments := make([]*models.Deployment, 0)

	var count int64
	if err := db.Order("updated_at desc").
		Offset((pageNum - 1) * pageSize).
		Find(&Deployments).
		Count(&count).Error; err != nil {
		return nil, err
	}

	if err := db.Order("updated_at desc").
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&Deployments).Error; err != nil {
		return nil, err
	}

	deploymentPagination := &models.DeploymentPagination{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    count,
		Pages:    utils.CalcPages(count, pageSize),
		Items:    Deployments,
	}
	return deploymentPagination, nil
}

// get deployment by deployment id
func GetDeploymentById(id int64) (*models.Deployment, error) {
	deployment := new(models.Deployment)
	db := database.GetDBEngine()
	if err := db.First(deployment, id).Error; err != nil {
		if err.Error() == "record not found" {
			e := utils.NewRockError(404, 40400011, fmt.Sprintf("Deployment with id(%v) was not found", id))
			return nil, e
		}
		return nil, err
	}
	return deployment, nil
}

// delete deployment by deployment id
func DeleteDeploymentById(id int64) error {
	deployment, err := GetDeploymentById(id)
	if err != nil {
		return err
	}

	db := database.GetDBEngine()
	if err := db.Delete(deployment, id).Error; err != nil {
		return err
	}
	return nil
}

// update the deployment by id app_id env_id chart_name chart_version
func UpdateDeploymentById(id, appId, envId int64, chartName, chartVersion, description string) (*models.Deployment, error) {
	deployment, err := GetDeploymentById(id)
	if err != nil {
		return nil, err
	}
	_, err = GetAppById(appId)
	if err != nil {
		return nil, err
	}

	_, err = GetEnvById(envId)
	if err != nil {
		return nil, err
	}

	db := database.GetDBEngine()
	if err := db.Model(deployment).Update(map[string]interface{}{"app_id": appId, "env_id": envId, "chart_name": chartName, "chart_version": chartVersion, "description": description}).Error; err != nil {
		return nil, err
	}

	return deployment, nil
}

// get all deployment by instance_name(equal deployment_name)
func GetDeploymentsByName(name string, pageNum, pageSize int64) (*models.DeploymentPagination, error) {
	db := database.GetDBEngine()
	Deployments := make([]*models.Deployment, 0)

	var count int64
	if err := db.Order("updated_at desc").
		Offset((pageNum-1)*pageSize).
		Where("name = ?", name).
		Find(&Deployments).
		Count(&count).Error; err != nil {
		return nil, err
	}

	if err := db.Order("updated_at desc").
		Offset((pageNum-1)*pageSize).
		Where("name = ?", name).
		Limit(pageSize).
		Find(&Deployments).Error; err != nil {
		return nil, err
	}

	deploymentPagination := &models.DeploymentPagination{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    count,
		Pages:    utils.CalcPages(count, pageSize),
		Items:    Deployments,
	}
	return deploymentPagination, nil
}
