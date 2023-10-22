// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package entity

import (
	"time"
)

const TableNameTaskRunLog = "task_run_logs"

// TaskRunLog mapped from table <task_run_logs>
type TaskRunLog struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	TaskID    int       `gorm:"column:task_id" json:"task_id"`
	TaskRunID int       `gorm:"column:task_run_id" json:"task_run_id"`
	UserID    int       `gorm:"column:user_id" json:"user_id"`
	Content   string    `gorm:"column:content" json:"content"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	Images    string    `gorm:"column:images" json:"images"`
	Videos    string    `gorm:"column:videos" json:"videos"`
}

// TableName TaskRunLog's table name
func (*TaskRunLog) TableName() string {
	return TableNameTaskRunLog
}
