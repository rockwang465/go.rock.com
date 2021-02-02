package api

import (
	"fmt"
	"go.rock.com/rock-platform/rock/server/database"
	"go.rock.com/rock-platform/rock/server/database/models"
	"go.rock.com/rock-platform/rock/server/utils"
)

func CreateProject(name, description string) (*models.Project, error) {
	db := database.GetDBEngine()
	if err := HasNotProjectByName(name); err != nil {
		return nil, err
	}
	project := &models.Project{
		Name:        name,
		Description: description,
	}

	if err := db.Create(project).Error; err != nil {
		return nil, err
	}
	return project, nil
}

func GetProjects(pageNum, pageSize int64, queryField string) (*models.ProjectPagination, error) {
	db := database.GetDBEngine()
	query := "%" + queryField + "%"
	Projects := make([]*models.Project, 0)

	var count int64
	if err := db.Order("updated_at desc").
		Offset((pageNum-1)*pageSize).
		Where("name like ?", query).
		Find(&Projects).
		Count(&count).Error; err != nil {
		return nil, err
	}

	if err := db.Order("updated_at desc").
		Offset((pageNum-1)*pageSize).
		Where("name like ?", query).
		Limit(pageSize).
		Find(&Projects).Error; err != nil {
		return nil, err
	}

	projectPagination := &models.ProjectPagination{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    count,
		Pages:    utils.CalcPages(count, pageSize),
		Items:    Projects,
	}
	return projectPagination, nil
}

// ensure project not exists
func HasNotProjectByName(name string) error {
	db := database.GetDBEngine()
	project := new(models.Project)

	if err := db.Where("name = ?", name).First(project).Error; err != nil {
		if err.Error() == "record not found" {
			return nil
		}
		return err
	}

	err := utils.NewRockError(400, 40000015, fmt.Sprintf("Project with name(%v) is alerady exist", name))
	return err
}

func GetProjectById(id int64) (*models.Project, error) {
	db := database.GetDBEngine()
	project := new(models.Project)
	if err := db.First(project, id).Error; err != nil {
		if err.Error() == "record not found" {
			err := utils.NewRockError(404, 40400004, fmt.Sprintf("Project with id(%v) is not found", id))
			return nil, err
		}
		return nil, err
	}
	return project, nil
}

func DeleteProjectById(id int64) error {
	db := database.GetDBEngine()
	project, err := GetProjectById(id)
	if err != nil {
		return err
	}

	if err := db.Delete(project, id).Error; err != nil {
		return err
	}
	return nil
}

func UpdateProject(id int64, desc string) (*models.Project, error) {
	db := database.GetDBEngine()
	project, err := GetProjectById(id)
	if err != nil {
		return nil, err
	}

	if err := db.Model(project).Update(map[string]interface{}{"description": desc}).Error; err != nil {
		return nil, err
	}
	return project, nil
}
