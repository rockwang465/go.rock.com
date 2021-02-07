package models

import "time"

type Common struct {
	CreatedAt LocalTime  `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt LocalTime  `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"deleted_at"` // can not use LocalTime
	//DeletedAt LocalTime `json:"deleted_at" gorm:"type:timestamp;default:null"`  // error
	//DeletedAt LocalTime `json:"deleted_at" gorm:"type:timestamp null"`  // error
	Version int `json:"version" gorm:"not null"`
}

type ConfCtx struct {
	UserId     int64  `json:"user_id"`
	Username   string `json:"username"`
	Role       string `json:"role"`
	DroneToken string `json:"drone_token"`
}
