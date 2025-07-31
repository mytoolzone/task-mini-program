package entity

import "time"

type Insurance struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	StartAt   time.Time `gorm:"column:start_at" json:"start_at"`
	EndAt     time.Time `gorm:"column:end_at" json:"end_at"`
	UserId    int       `gorm:"column:user_id" json:"user_id"`
}

func TableName() string {
	return "insurances"
}
