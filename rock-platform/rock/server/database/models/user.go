package models

import "time"

type User struct {
	Id       int64  `json:"id" gorm:"primary_key"`
	Name     string `json:"name" gorm:"not null;unique_index;type:varchar(50)"`
	Password string `json:"password" gorm:"not null"`
	Email    string `json:"email" gorm:"not null"`
	Salt     string `json:"salt" gorm:"not null"`
	//GitlabToken     string     `json:"gitlab_token"`
	//DroneToken      string     `json:"drone_token"`
	Token string `json:"token" gorm:"type:varchar(1024)"`
	//ResetSecret     string     `json:"reset_secret"`
	//SecretExpiredAt *time.Time `json:"secret_expired_at"`
	LoginRetryCount int64      `json:"login_retry_count"`
	LoginBlockUntil *time.Time `json:"login_block_until"`
	RoleId          *Role      `json:"role_id"`
	Common
}
