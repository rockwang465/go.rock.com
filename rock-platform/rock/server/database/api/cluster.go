package api

import (
	"fmt"
	"go.rock.com/rock-platform/rock/server/database"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
)

func CreateCluster(name, description, config string) (*models.Cluster, error) {
	db := database.GetDBEngine()
	if err := hasNotClusterWithSameName(name); err != nil {
		return nil, err
	}

	cluster := models.Cluster{
		Name:        name,
		Description: description,
		Config:      config,
	}
	if err := db.Create(&cluster).Error; err != nil {
		return nil, err
	}
	return &cluster, nil
}

func GetClusters(pageNum, pageSize int64, queryField string) (*models.ClusterPagination, error) {
	db := database.GetDBEngine()
	query := "%" + queryField + "%"
	Clusters := make([]*models.Cluster, 0)

	var count int64
	if err := db.Order("updated_at desc").
		Offset((pageNum-1)*pageSize).
		Where("name like ?", query).
		Find(&Clusters).
		Count(&count).Error; err != nil {
		return nil, err
	}

	if err := db.Order("updated_at desc").
		Offset((pageNum-1)*pageSize).
		Where("name like ?", query).
		Limit(pageSize).
		Find(&Clusters).Error; err != nil {
		return nil, err
	}

	clusterPagination := &models.ClusterPagination{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    count,
		Pages:    utils.CalcPages(count, pageSize),
		Items:    Clusters,
	}
	return clusterPagination, nil
}

func GetClusterById(id int64) (*models.Cluster, error) {
	db := database.GetDBEngine()
	cluster := new(models.Cluster)
	if err := db.First(cluster, id).Error; err != nil {
		if err.Error() == "record not found" {
			err := utils.NewRockError(404, 40400007, fmt.Sprintf("Cluster with id(%v) was not found", id))
			return nil, err
		}
		return nil, err
	}
	return cluster, nil
}

func DeleteClusterById(id int64) error {
	db := database.GetDBEngine()
	cluster, err := GetClusterById(id)
	if err != nil {
		return err
	}

	if err := db.Delete(cluster, id).Error; err != nil {
		return err
	}
	return nil
}

func UpdateCluster(id int64, desc, config string) (*models.Cluster, error) {
	db := database.GetDBEngine()
	cluster, err := GetClusterById(id)
	if err != nil {
		return nil, err
	}

	if err := db.Model(cluster).Update(map[string]interface{}{"description": desc, "config": config}).Error; err != nil {
		return nil, err
	}
	return cluster, nil
}

func GetClusterEnvsById(clusterId, pageNum, pageSize int64) (*models.EnvPagination, error) {
	db := database.GetDBEngine()
	Envs := make([]*models.Env, 0)

	_, err := GetClusterById(clusterId)
	if err != nil {
		return nil, err
	}

	var count int64
	if err := db.Order("updated_at desc").
		Offset((pageNum-1)*pageSize).
		Where("cluster_id = ?", clusterId).
		Find(&Envs).
		Count(&count).Error; err != nil {
		return nil, err
	}

	if err := db.Order("updated_at desc").
		Offset((pageNum-1)*pageSize).
		Where("cluster_id = ?", clusterId).
		Limit(pageSize).
		Find(&Envs).Error; err != nil {
		return nil, err
	}

	envPagination := &models.EnvPagination{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    count,
		Pages:    utils.CalcPages(count, pageSize),
		Items:    Envs,
	}
	return envPagination, nil
}

func hasNotClusterWithSameName(name string) error {
	db := database.GetDBEngine()
	cluster := new(models.Cluster)
	if err := db.Where("name = ?", name).First(cluster).Error; err != nil {
		if err.Error() == "record not found" {
			return nil
		}
		return err
	}
	err := utils.NewRockError(400, 40000018, fmt.Sprintf("Cluster with name(%v) already exists", name))
	return err
}
