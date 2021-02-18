package api

import (
	"fmt"
	"go.rock.com/rock-platform/rock/server/database"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
)

func CreateProjectEnvByProjectId(projectId, envId int64, projectEnvName, desc string) (*models.ProjectEnv, error) {
	// ensure has project_id
	_, err := GetProjectById(projectId)
	if err != nil {
		return nil, err
	}

	// ensure has env_id
	_, err = GetEnvById(envId)
	if err != nil {
		panic(err)
	}

	db := database.GetDBEngine()
	projectEnv := &models.ProjectEnv{
		Name:        projectEnvName,
		Description: desc,
		EnvId:       envId,
		ProjectId:   projectId,
	}

	// ensure has not same project_env
	if err := hasNotSameProjectEnv(projectEnv); err != nil {
		return nil, err
	}

	// create project_env
	if err := db.Create(projectEnv).Error; err != nil {
		return nil, err
	}
	return projectEnv, nil
}

func GetProjectEnvs(projectId, pageNum, pageSize int64, queryField string) (*models.ProjectEnvPagination, error) {
	_, err := GetProjectById(projectId)
	if err != nil {
		return nil, err
	}

	db := database.GetDBEngine()
	query := "%" + queryField + "%"
	ProjectEnvs := make([]*models.BriefProjectEnv, 0)

	var count int64
	if err := db.Table("project_env").
		Order("updated_at desc").
		Offset((pageNum-1)*pageSize).
		Where("project_id = ? AND name like ?", projectId, query).
		Find(&ProjectEnvs).
		Count(&count).Error; err != nil {
		return nil, err
	}

	if err := db.Table("project_env").
		Order("updated_at desc").
		Offset((pageNum-1)*pageSize).
		Where("project_id = ? AND name like ?", projectId, query).
		Limit(pageSize).
		Find(&ProjectEnvs).Error; err != nil {
		return nil, err
	}

	projectEvnPagination := &models.ProjectEnvPagination{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    count,
		Pages:    utils.CalcPages(count, pageSize),
		Items:    ProjectEnvs,
	}
	return projectEvnPagination, nil
}

func GetProjectEnvById(id int64) (*models.ProjectEnv, error) {
	db := database.GetDBEngine()
	projectEnv := new(models.ProjectEnv)
	if err := db.First(projectEnv, id).Error; err != nil {
		if err.Error() == "record not found" {
			err := utils.NewRockError(404, 40400009, fmt.Sprintf("ProjectEnv with id(%v) was not found", id))
			return nil, err
		}
		return nil, err
	}
	return projectEnv, nil
}

func hasNotSameProjectEnv(projectEnv *models.ProjectEnv) error {
	fmt.Printf("%#v\n", projectEnv)
	db := database.GetDBEngine()
	pe := new(models.ProjectEnv)
	if err := db.Where(projectEnv).First(pe).Error; err != nil {
		if err.Error() == "record not found" {
			return nil
		}
		return err
	}
	err := utils.NewRockError(400, 40000022, fmt.Sprintf("ProjectEnv with name(%v) project_id(%v) env_id(%v) already exists", projectEnv.Name, projectEnv.ProjectId, projectEnv.EnvId))
	return err
}
