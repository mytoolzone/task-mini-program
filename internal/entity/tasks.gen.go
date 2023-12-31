// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package entity

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

const TableNameTask = "tasks"

// Task mapped from table <tasks>
type Task struct {
	ID         int            `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" `
	Name       string         `gorm:"column:name;not null" json:"name"`
	CreateBy   int            `gorm:"column:create_by" json:"create_by"`
	CreatedAt  time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at" swaggerignore:"true"`
	FinishedAt time.Time      `gorm:"column:finished_at" json:"finished_at" `
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at" swaggerignore:"true"`
	Describe   string         `gorm:"column:describe" json:"describe"`
	Require    string         `gorm:"column:require" json:"require"`
	Location   string         `gorm:"column:location" json:"location"`
	Status     string         `gorm:"column:status" json:"status"`
}

// TableName Task's table name
func (*Task) TableName() string {
	return TableNameTask
}

type Point struct {
	lat float64 `json:"lat"`
	lng float64 `json:"lng"`
}

// 实现driver.Valuer接口
func (p *Point) Value() (driver.Value, error) {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "(%f, %f)", p.lat, p.lng)
	return buf.Bytes(), nil
}

func (p *Point) String() string {
	return fmt.Sprintf("(%v, %v)", p.lat, p.lng)
}

// 实现sql.scanner接口
func (p *Point) Scan(val interface{}) (err error) {
	if bb, ok := val.([]uint8); ok {
		tmp := bb[1 : len(bb)-1]
		coors := strings.Split(string(tmp[:]), ",")
		if p.lat, err = strconv.ParseFloat(coors[0], 64); err != nil {
			return err
		}
		if p.lng, err = strconv.ParseFloat(coors[1], 64); err != nil {
			return err
		}
	}
	return nil
}
