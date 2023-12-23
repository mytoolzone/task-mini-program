package entity

import (
	"time"
)

type File struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name      string    `gorm:"column:name;not null" json:"name"`
	Size      int64     `gorm:"column:size;not null" json:"size"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	DeletedAt time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	Path      string    `gorm:"column:path" json:"path"`
}
