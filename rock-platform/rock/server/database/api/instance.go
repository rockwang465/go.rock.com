package api

import (
	"fmt"
	"go.rock.com/rock-platform/rock/server/database"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
)

// create or update deployment chart instance info
func CreateOrUpdateInstance(chartName, chartVersion, clusterName, namespace, projectName string, deployment *models.Deployment) (*models.Instance, error) {
	db := database.GetDBEngine()
	instance := &models.Instance{}

	if err := db.Where("cluster_name = ? AND env_namespace = ? AND project_name = ? AND name = ?", clusterName, namespace, projectName, deployment.Name).First(instance).Error; err != nil {
		if err.Error() != "record not found" {
			return nil, err
		}
	}

	if instance.Id != 0 { // 不为0，表示存在记录，则更新
		updateIns := map[string]interface{}{"last_deployment": deployment.Id, "chart_name": chartName, "chart_version": chartVersion}
		if err := db.Model(instance).Update(updateIns).Error; err != nil {
			return nil, err
		}
	} else { // 为0，表示不存在记录，则创建
		instance := &models.Instance{
			ClusterName:    clusterName,
			EnvNamespace:   namespace,
			ProjectName:    projectName,     // example: sensenebula-guard-std 、idea-aurora
			Name:           deployment.Name, // example: senseguard-watchlist-management-default
			ChartName:      chartName,       // example: senseguard-watchlist-management
			ChartVersion:   chartVersion,
			LastDeployment: deployment.Id,
			AppId:          deployment.AppId,
			EnvId:          deployment.EnvId,
		}
		if err := db.Create(instance).Error; err != nil {
			return nil, err
		}
	}

	return instance, nil
}

// get instances by app_id
func GetAppInstances(appId, pageNum, pageSize int64) (*models.InstancePagination, error) {
	db := database.GetDBEngine()
	Instances := make([]*models.Instance, 0)

	var count int64
	if err := db.Order("updated_at desc").
		Offset((pageNum-1)*pageSize).
		Where("app_id = ?", appId).
		Find(&Instances).
		Count(&count).Error; err != nil {
		return nil, err
	}

	if err := db.Order("updated_at desc").
		Offset((pageNum-1)*pageSize).
		Where("app_id = ?", appId).
		Limit(pageSize).
		Find(&Instances).Error; err != nil {
		return nil, err
	}

	instancePagination := &models.InstancePagination{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    count,
		Pages:    utils.CalcPages(count, pageSize),
		Items:    Instances,
	}
	return instancePagination, nil
}

// get instances by pageNum, pageSize, appId(app_id), queryField, cluster(cluster_name), project(project_name)
func GetInstances(pageNum, pageSize int64, queryField, cluster, project string, appId int64) (*models.InstancePagination, error) {
	query := "%" + queryField + "%"

	db := database.GetDBEngine()
	Instances := make([]*models.Instance, 0)
	//generate a query DB object
	queryDb := db.Model(&Instances).Order("updated_at desc").
		Offset((pageNum-1)*pageSize).
		Where("name like ?", query)

	if cluster != "" {
		queryDb = queryDb.Where("cluster_name = ?", cluster)
	}
	if project != "" {
		queryDb = queryDb.Where("project_name = ?", project)
	}
	if appId != 0 {
		queryDb = queryDb.Where("app_id = ?", appId)
	}

	var count int64
	if err := queryDb.Limit(pageSize).
		Find(&Instances).
		Count(&count).Error; err != nil {
		return nil, err
	}

	if err := queryDb.Limit(pageSize).
		Find(&Instances).Error; err != nil {
	}

	instancePagination := &models.InstancePagination{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    count,
		Pages:    utils.CalcPages(count, pageSize),
		Items:    Instances,
	}
	return instancePagination, nil
}

// get instance by instance id
func GetInstanceById(id int64) (*models.Instance, error) {
	db := database.GetDBEngine()
	instance := new(models.Instance)
	if err := db.First(instance, id).Error; err != nil {
		if err.Error() == "record not found" {
			e := utils.NewRockError(404, 40400012, fmt.Sprintf("Instance with id(%v) was not found", id))
			return nil, e
		}
		return nil, err
	}
	return instance, nil
}

// delete instance by instance id
func DeleteInstanceById(id int64) error {
	instance, err := GetInstanceById(id)
	if err != nil {
		return err
	}

	db := database.GetDBEngine()
	if err := db.Delete(instance).Error; err != nil {
		return err
	}
	return nil
}
