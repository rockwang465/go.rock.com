package models

type Common struct {
	CreatedAt LocalTime `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt LocalTime `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	//DeletedAt   LocalTime     `json:"deleted_at"`
	Version int `json:"version" gorm:"not null"`
}
