package repo

import (
	"context"
	"errors"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/pkg/postgres"
	"gorm.io/gorm"
	"time"
)

type TaskRun struct {
	*postgres.Postgres
}

func NewTaskRunRepo(pg *postgres.Postgres) *TaskRun {
	return &TaskRun{pg}
}

func (t TaskRun) AddTaskRun(ctx context.Context, taskID int) (entity.TaskRun, error) {
	var taskRun entity.TaskRun = entity.TaskRun{
		TaskID:  taskID,
		Status:  entity.TaskStatusTorun,
		StartAt: time.Now(),
	}

	err := t.Db.Create(&taskRun).Error
	if err != nil {
		return entity.TaskRun{}, err
	}
	return taskRun, nil
}

func (t TaskRun) GetPendingTaskRun(ctx context.Context, taskID int) (entity.TaskRun, error) {
	var taskRun entity.TaskRun
	err := t.Db.WithContext(ctx).Where("task_id = ? and status = ?",
		taskID, entity.TaskStatusTorun).First(&taskRun).Error
	if err != nil {
		return entity.TaskRun{}, err
	}

	return taskRun, nil
}

func (t TaskRun) GetRunningTaskRun(ctx context.Context, taskID int) (entity.TaskRun, error) {
	var taskRun entity.TaskRun
	err := t.Db.WithContext(ctx).Where("task_id = ? and status = ?", taskID, entity.TaskStatusRunning).First(&taskRun).Error
	if err != nil {
		return entity.TaskRun{}, err
	}

	return taskRun, nil
}

func (t TaskRun) StartTaskRun(ctx context.Context, taskID int) error {
	var taskRun entity.TaskRun
	err := t.Db.Debug().WithContext(ctx).Where("task_id = ? and status = ?",
		taskID, entity.TaskStatusTorun).First(&taskRun).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	taskRun.Status = entity.TaskStatusRunning
	return t.Db.WithContext(ctx).Debug().Where("task_id = ?", taskID).Updates(&taskRun).Error
}

func (t TaskRun) FinishTaskRun(ctx context.Context, taskID int) error {
	var taskRun entity.TaskRun
	err := t.Db.WithContext(ctx).Where("task_id = ? and status = ?", taskID, entity.TaskStatusRunning).First(&taskRun).Error
	if err != nil {
		return err
	}
	taskRun.Status = entity.TaskStatusFinished
	return t.Db.WithContext(ctx).Where("task_id = ?", taskID).Updates(&taskRun).Error
}

func (t TaskRun) CancelTaskRun(ctx context.Context, taskID int) error {
	var taskRun entity.TaskRun
	err := t.Db.WithContext(ctx).Where("task_id = ? and status = ?", taskID, entity.TaskStatusRunning).First(&taskRun).Error
	if err != nil {
		return err
	}
	taskRun.Status = entity.TaskStatusCanceled
	return t.Db.WithContext(ctx).Where("task_id = ?", taskID).Updates(&taskRun).Error
}

func (t TaskRun) GetTaskRunList(ctx context.Context, taskID int) ([]entity.TaskRun, error) {
	var taskRuns []entity.TaskRun
	err := t.Db.WithContext(ctx).Where("task_id = ?", taskID).Find(&taskRuns).Error
	return taskRuns, err
}
