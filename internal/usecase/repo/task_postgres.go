package repo

import (
	"context"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/pkg/postgres"
)

type TaskRepo struct {
	*postgres.Postgres
}

func (t *TaskRepo) AuditFailTask(ctx context.Context, taskID int) error {
	var task entity.Task
	err := t.Db.WithContext(ctx).Where("id = ?", taskID).First(&task).Error
	if err != nil {
		return err
	}
	task.Status = entity.TaskStatusAuditFail
	return t.Db.WithContext(ctx).Where("id = ?", taskID).Updates(&task).Error
}

func (t *TaskRepo) AuditSuccessTask(ctx context.Context, taskID int) error {
	var task entity.Task
	err := t.Db.WithContext(ctx).Where("id = ?", taskID).First(&task).Error
	if err != nil {
		return err
	}
	// 任务审核通过待执行
	task.Status = entity.TaskStatusPending
	return t.Db.WithContext(ctx).Where("id = ?", taskID).Updates(&task).Error
}

func (t *TaskRepo) CreateTask(ctx context.Context, task *entity.Task) error {
	task.Status = entity.TaskStatusPublish
	return t.Db.Create(task).Error
}

func (t *TaskRepo) GetByUserID(ctx context.Context, userID int) ([]entity.Task, error) {
	var tasks []entity.Task
	err := t.Db.WithContext(ctx).Where("user_id = ?", userID).Find(&tasks).Error
	return tasks, err
}

func (t *TaskRepo) GetByTaskID(ctx context.Context, taskID int) (entity.Task, error) {
	var task entity.Task
	err := t.Db.WithContext(ctx).Where("id = ?", taskID).First(&task).Error
	return task, err
}

// GetTaskList 瀑布流方式获取任务列表
func (t *TaskRepo) GetTaskList(ctx context.Context, lastId int) ([]entity.Task, error) {
	var tasks []entity.Task
	err := t.Db.WithContext(ctx).Where("id < ?", lastId).Find(&tasks).Error
	return tasks, err
}

func (t *TaskRepo) StartTask(ctx context.Context, taskID int) error {
	var task entity.Task
	err := t.Db.WithContext(ctx).Where("id = ?", taskID).First(&task).Error
	if err != nil {
		return err
	}
	task.Status = entity.TaskStatusRunning
	return t.Db.WithContext(ctx).Where("id = ?", taskID).Updates(&task).Error
}

func (t *TaskRepo) PauseTask(ctx context.Context, taskID int) error {
	var task entity.Task
	err := t.Db.WithContext(ctx).Where("id = ?", taskID).First(&task).Error
	if err != nil {
		return err
	}
	task.Status = entity.TaskStatusPaused
	return t.Db.WithContext(ctx).Where("id = ?", taskID).Updates(&task).Error
}

func (t *TaskRepo) FinishTask(ctx context.Context, taskID int) error {
	var task entity.Task
	err := t.Db.WithContext(ctx).Where("id = ?", taskID).First(&task).Error
	if err != nil {
		return err
	}
	task.Status = entity.TaskStatusFinished
	return t.Db.WithContext(ctx).Where("id = ?", taskID).Updates(&task).Error
}

func (t *TaskRepo) CancelTask(ctx context.Context, taskID int) error {
	var task entity.Task
	err := t.Db.WithContext(ctx).Where("id = ?", taskID).First(&task).Error
	if err != nil {
		return err
	}
	task.Status = entity.TaskStatusCanceled
	return t.Db.WithContext(ctx).Where("id = ?", taskID).Updates(&task).Error
}

func NewTaskRepo(pg *postgres.Postgres) *TaskRepo {
	return &TaskRepo{pg}
}
