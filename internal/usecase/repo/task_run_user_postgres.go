package repo

import (
	"context"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/pkg/postgres"
	"time"
)

type TaskRunUserRepo struct {
	*postgres.Postgres
}

func NewTaskRunUserRepo(pg *postgres.Postgres) *TaskRunUserRepo {
	return &TaskRunUserRepo{pg}
}

func (t TaskRunUserRepo) AddTaskRunUser(ctx context.Context, taskID, taskRunID, userID int) (entity.TaskRunUser, error) {
	// 判断是否存在避免重复添加
	var existTaskRunUser = entity.TaskRunUser{}
	t.Db.WithContext(ctx).Where("task_run_id = ? and user_id = ?", taskRunID, userID).First(&existTaskRunUser)
	if existTaskRunUser.ID > 0 {
		return existTaskRunUser, nil
	}

	var taskRunUser entity.TaskRunUser = entity.TaskRunUser{
		TaskID:    taskID,
		TaskRunID: taskRunID,
		UserID:    userID,
		Status:    entity.TaskStatusSign,
		CreatedAt: time.Now(),
	}
	if err := t.Db.WithContext(ctx).Create(&taskRunUser).Error; err != nil {
		return entity.TaskRunUser{}, err
	}
	return taskRunUser, nil
}

func (t TaskRunUserRepo) StartTaskRun(ctx context.Context, taskID, taskRunID int) error {
	tru := entity.TaskRunUser{
		Status:  entity.TaskStatusRunning,
		StartAt: time.Now(),
	}
	// 更新task_run_user表中 startAt 为当前时间
	return t.Db.WithContext(ctx).Model(&entity.TaskRunUser{}).Where("task_id = ? and task_run_id = ?", taskID, taskRunID).
		Updates(tru).Error
}

func (t TaskRunUserRepo) FinishTaskRun(ctx context.Context, taskID, taskRunID int) error {
	tru := entity.TaskRunUser{
		Status:     entity.TaskStatusFinished,
		FinishedAt: time.Now(),
	}
	// 更新task_run_user表中 finished_at 为当前时间
	return t.Db.Debug().WithContext(ctx).Where("task_id = ? and task_run_id =? ", taskID, taskRunID).
		Updates(tru).Error
}

func (t TaskRunUserRepo) CancelTaskRun(ctx context.Context, taskID, taskRunID int) error {
	// 更新task_run_user表中 finished_at 为当前时间
	return t.Db.WithContext(ctx).Where("task_id = ? and task_run_id ", taskID, taskRunID).Updates(map[string]string{
		"finished_at": "now()",
		"status":      entity.TaskStatusCanceled,
	}).Error
}

func (t TaskRunUserRepo) GetTaskRunUserList(ctx context.Context, taskID int, taskRunID int) ([]entity.TaskRunUser, error) {
	var taskRunUsers []entity.TaskRunUser
	err := t.Db.WithContext(ctx).Where("task_id = ? and task_run_id ", taskID, taskRunID).Find(&taskRunUsers).Error
	return taskRunUsers, err
}

func (t TaskRunUserRepo) GetUserTaskSummary(ctx context.Context, userID int) (entity.UserTaskSummary, error) {
	var userTaskSummary entity.UserTaskSummary
	// 获取有用户参与的任务总数
	err := t.Db.WithContext(ctx).Model(&entity.UserTask{}).Where("user_id = ?", userID).Count(&userTaskSummary.TotalTask).Error
	if err != nil {
		return entity.UserTaskSummary{}, err
	}
	// 统计用户参与任务总时长
	err = t.Db.WithContext(ctx).Model(&entity.TaskRunUser{}).
		Where("user_id = ? and status = ?", userID, entity.TaskStatusFinished).
		Select("sum(duration)").Scan(&userTaskSummary.TotalDuration).Error

	if err != nil {
		return entity.UserTaskSummary{}, err
	}

	return userTaskSummary, nil
}
