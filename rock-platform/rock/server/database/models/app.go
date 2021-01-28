package models

type App struct {
	Id              int64  `json:"id" gorm:"unique_index;not null"`
	Name            string `json:"name" gorm:"not null"`
	FullName        string `json:"full_name" gorm:"not null"`
	Owner           string `json:"owner" gorm:"not null"`
	Description     string `json:"description"`
	GitlabAddress   string `json:"gitlab_address" gorm:"not null"`
	ProjectId       int64  `json:"project_id" gorm:"not null"`
	DroneRepoId     int64  `json:"drone_repo_id" gorm:"not null"`
	GitlabProjectId int64  `json:"gitlab_project_id" gorm:"not null"`
	Common
}
