package repo

import (
	"context"
	"github.com/mytoolzone/task-mini-program/internal/app_code"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/pkg/postgres"
)

type UserTaskRepo struct {
	*postgres.Postgres
}

func NewUserTaskRepo(pg *postgres.Postgres) *UserTaskRepo {
	return &UserTaskRepo{pg}
}

func (t UserTaskRepo) AddUserTask(ctx context.Context, taskID, userID int) (entity.UserTask, error) {
	var UserTask entity.UserTask = entity.UserTask{
		TaskID: taskID,
		UserID: userID,
		Status: entity.UserTaskStatusApply,
	}
	if err := t.Db.WithContext(ctx).Create(&UserTask).Error; err != nil {
		return entity.UserTask{}, err
	}
	return UserTask, nil
}

func (t UserTaskRepo) AuditUserTask(ctx context.Context, taskID, userID int, status string) (entity.UserTask, error) {
	if status != entity.UserTaskStatusAuditFail && status != entity.UserTaskStatusAuditPass {
		return entity.UserTask{}, app_code.New(app_code.ErrorAuditParamInValid, "status invalid")
	}

	var UserTask entity.UserTask = entity.UserTask{
		TaskID: taskID,
		UserID: userID,
		Status: status,
	}

	if err := t.Db.WithContext(ctx).Where("task_id = ? and user_id = ?", taskID, userID).Updates(&UserTask).Error; err != nil {
		return entity.UserTask{}, err
	}
	return UserTask, nil
}

func (t UserTaskRepo) GetUserTaskList(ctx context.Context, taskID int) ([]entity.UserTask, error) {
	var UserTasks []entity.UserTask
	err := t.Db.WithContext(ctx).Where("task_id = ?", taskID).Find(&UserTasks).Error
	return UserTasks, err
}

// GetUserTaskByUserID 获取某个人参与的任务的状态
func (t UserTaskRepo) GetUserTaskByUserID(ctx context.Context, taskID, userID int) (entity.UserTask, error) {
	var UserTask entity.UserTask
	err := t.Db.WithContext(ctx).Where("task_id = ? and user_id = ?", taskID, userID).First(&UserTask).Error
	return UserTask, err
}