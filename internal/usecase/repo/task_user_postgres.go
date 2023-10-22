package repo

import (
	"context"
	"github.com/mytoolzone/task-mini-program/internal/app_error"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/pkg/postgres"
)

type TaskUserRepo struct {
	*postgres.Postgres
}

func NewTaskUserRepo(pg *postgres.Postgres) *TaskUserRepo {
	return &TaskUserRepo{pg}
}

func (t TaskUserRepo) AddTaskUser(ctx context.Context, taskID, userID int) (entity.TaskUser, error) {
	var taskUser entity.TaskUser = entity.TaskUser{
		TaskID: taskID,
		UserID: userID,
		Status: entity.UserTaskStatusApply,
	}
	if err := t.Db.WithContext(ctx).Create(&taskUser).Error; err != nil {
		return entity.TaskUser{}, err
	}
	return taskUser, nil
}

func (t TaskUserRepo) AuditTaskUser(ctx context.Context, taskID, userID int, status string) (entity.TaskUser, error) {
	if status != entity.UserTaskStatusAuditFail && status != entity.UserTaskStatusAuditPass {
		return entity.TaskUser{}, app_error.New(app_error.ErrorAuditParamInValid, "status invalid")
	}

	var taskUser entity.TaskUser = entity.TaskUser{
		TaskID: taskID,
		UserID: userID,
		Status: status,
	}

	if err := t.Db.WithContext(ctx).Where("task_id = ? and user_id = ?", taskID, userID).Updates(&taskUser).Error; err != nil {
		return entity.TaskUser{}, err
	}
	return taskUser, nil
}

func (t TaskUserRepo) GetTaskUserList(ctx context.Context, taskID int) ([]entity.TaskUser, error) {
	var taskUsers []entity.TaskUser
	err := t.Db.WithContext(ctx).Where("task_id = ?", taskID).Find(&taskUsers).Error
	return taskUsers, err
}

// GetTaskUserByUserID 获取某个人参与的任务的状态
func (t TaskUserRepo) GetTaskUserByUserID(ctx context.Context, taskID, userID int) (entity.TaskUser, error) {
	var taskUser entity.TaskUser
	err := t.Db.WithContext(ctx).Where("task_id = ? and user_id = ?", taskID, userID).First(&taskUser).Error
	return taskUser, err
}
