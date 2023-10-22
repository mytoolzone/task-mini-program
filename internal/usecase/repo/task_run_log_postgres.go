package repo

import (
	"context"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/pkg/postgres"
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

func (t TaskRunLogRepo) GetTaskRunLogList(ctx context.Context, taskID int) ([]entity.TaskRunLog, error) {
	var taskRunLogs []entity.TaskRunLog
	err := t.Db.WithContext(ctx).Where("task_id = ?", taskID).First(&taskRunLogs).Error
	return taskRunLogs, err
}
