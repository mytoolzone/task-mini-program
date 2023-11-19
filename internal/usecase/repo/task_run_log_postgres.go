package repo

import (
	"context"
	"errors"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/pkg/postgres"
	"gorm.io/gorm"
)

type TaskRunLogRepo struct {
	*postgres.Postgres
}

func NewTaskRunLogRepo(pg *postgres.Postgres) *TaskRunLogRepo {
	return &TaskRunLogRepo{pg}
}

func (t TaskRunLogRepo) AddTaskRunLog(ctx context.Context, log *entity.TaskRunLog) error {
	if err := t.Db.WithContext(ctx).Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func (t TaskRunLogRepo) GetTaskRunLogList(ctx context.Context, taskID, lastID int) ([]entity.TaskRunLog, error) {
	var taskRunLogs []entity.TaskRunLog
	query := t.Db.WithContext(ctx).Where("task_id = ?", taskID)

	if lastID > 0 {
		query = query.Where("id <?", lastID)
	}

	err := query.Limit(50).Order("id desc").Find(&taskRunLogs).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return taskRunLogs, nil
}
