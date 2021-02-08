package api

import (
	"fmt"
	"go.rock.com/rock-platform/rock/server/database"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
)

func CreateEnv(namespace, desc string, ClusterId int64) (*models.Env, error) {
	db := database.GetDBEngine()

	// must have cluster id in mysql
	_, err := GetClusterById(ClusterId)
	if err != nil {
		return nil, err
	}

	if err := hasNotEnvWithSameClusterAndNamespace(namespace, ClusterId); err != nil {
		return nil, err
	}

	env := &models.Env{
		Description: desc,
		Namespace:   namespace,
		ClusterId:   ClusterId,
	}
	if err := db.Create(env).Error; err != nil {
		return nil, err
	}
	return env, nil
}

func GetEnvs(pageNum, pageSize int64) (*models.EnvPagination, error) {
	db := database.GetDBEngine()
	Envs := make([]*models.Env, 0)

	var count int64
	if err := db.Order("updated_at desc").
		Offset((pageNum - 1) * pageSize).
		Find(&Envs).
		Count(&count).Error; err != nil {
		return nil, err
	}

	if err := db.Order("updated_at desc").
		Offset((pageNum - 1) * pageSize).
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

func GetEnvById(id int64) (*models.Env, error) {
	db := database.GetDBEngine()
	env := new(models.Env)
	if err := db.First(env, id).Error; err != nil {
		if err.Error() == "record not found" {
			err = utils.NewRockError(404, 40400008, fmt.Sprintf("Env with id(%v) is not found", id))
			return nil, err
		}
		return nil, err
	}
	return env, nil
}

func DeleteEnvById(id int64) error {
	db := database.GetDBEngine()
	env, err := GetEnvById(id)
	if err != nil {
		return err
	}

	if err := db.Delete(env).Error; err != nil {
		return err
	}
	return nil
}

func UpdateEnv(id int64, desc string) (*models.Env, error) {
	db := database.GetDBEngine()
	env, err := GetEnvById(id)
	if err != nil {
		return nil, err
	}

	if err := db.Model(env).Update(map[string]interface{}{"description": desc}).Error; err != nil {
		return nil, err
	}
	return env, nil
}

// ensure not namespace + cluster_id data in mysql
func hasNotEnvWithSameClusterAndNamespace(namespace string, ClusterId int64) error {
	db := database.GetDBEngine()
	env := new(models.Env)
	if err := db.Where("namespace = ? AND cluster_id = ?", namespace, ClusterId).First(env).Error; err != nil {
		if err.Error() == "record not found" {
			return nil
		}
		return err
	}
	err := utils.NewRockError(400, 40000020, fmt.Sprintf("Cluster with name(%v) cluster_id(%v) already exists", namespace, ClusterId))
	return err
}
