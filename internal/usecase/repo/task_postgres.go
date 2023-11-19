package repo

import (
	"context"
	"github.com/mytoolzone/task-mini-program/internal/app_code"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/pkg/postgres"
	"strings"
)

type TaskRepo struct {
	*postgres.Postgres
}

func (t *TaskRepo) AuditFailTask(ctx context.Context, taskID int) (*entity.Task, error) {
	var task entity.Task
	err := t.Db.WithContext(ctx).Where("id = ?", taskID).First(&task).Error
	if err != nil {
		return nil, err
	}

	if task.Status != entity.TaskStatusNew {
		return nil, app_code.New(app_code.ErrorRepeat, "任务已经审核完成")
	}

	// 任务审核失败
	task.Status = entity.TaskStatusAuditFail
	err = t.Db.WithContext(ctx).Where("id = ?", taskID).Updates(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (t *TaskRepo) AuditSuccessTask(ctx context.Context, taskID int) (*entity.Task, error) {
	var task entity.Task
	err := t.Db.WithContext(ctx).Where("id = ?", taskID).First(&task).Error
	if err != nil {
		return nil, err
	}

	if task.Status != entity.TaskStatusNew {
		return nil, app_code.New(app_code.ErrorRepeat, "任务已经审核完成")
	}

	// 任务审核通过待执行
	task.Status = entity.TaskStatusTorun
	err = t.Db.WithContext(ctx).Where("id = ?", taskID).Updates(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (t *TaskRepo) CreateTask(ctx context.Context, task *entity.Task) error {
	return t.Db.Create(task).Error
}

func (t *TaskRepo) GetByUserID(ctx context.Context, userID int) ([]entity.Task, error) {
	var tasks []entity.Task
	err := t.Db.WithContext(ctx).Where("create_by = ?", userID).Find(&tasks).Error
	return tasks, err
}

func (t *TaskRepo) GetByTaskID(ctx context.Context, taskID int) (entity.Task, error) {
	var task entity.Task
	err := t.Db.WithContext(ctx).Where("id = ?", taskID).First(&task).Error
	return task, err
}

// GetTaskList 瀑布流方式获取任务列表
func (t *TaskRepo) GetTaskList(ctx context.Context, lastId int, keyword, status string) ([]entity.Task, error) {
	var tasks []entity.Task
	query := t.Db.WithContext(ctx)
	if lastId > 0 {
		query = query.Where("id < ?", lastId)
	}

	if status != "" {
		statusArr := strings.Split(status, ",")
		query = query.Where("status in(?)", statusArr)
	}

	if keyword != "" {
		query = query.Where("name like ?", "%"+keyword+"%")
	}

	err := query.Order("id desc").Find(&tasks).Error
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
