package repo

import (
	"context"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/pkg/postgres"
)

type TaskRunUserRepo struct {
	*postgres.Postgres
}

func NewTaskRunUserRepo(pg *postgres.Postgres) *TaskRunUserRepo {
	return &TaskRunUserRepo{pg}
}

func (t TaskRunUserRepo) AddTaskRunUser(ctx context.Context, taskID, taskRunID, userID int) (entity.TaskRunUser, error) {
	var taskRunUser entity.TaskRunUser = entity.TaskRunUser{
		TaskID:    taskID,
		TaskRunID: taskRunID,
		UserID:    userID,
	}
	if err := t.Db.Create(&taskRunUser).Error; err != nil {
		return entity.TaskRunUser{}, err
	}
	return taskRunUser, nil
}

func (t TaskRunUserRepo) StartTaskRun(ctx context.Context, taskID, taskRunID int) error {
	// 更新task_run_user表中 startAt 为当前时间
	return t.Db.WithContext(ctx).Where("task_id = ? and task_run_id ", taskID, taskRunID).Updates(map[string]string{
		"start_at": "now()",
	}).Error
}

func (t TaskRunUserRepo) FinishTaskRun(ctx context.Context, taskID, taskRunID int) error {
	// 更新task_run_user表中 finished_at 为当前时间
	return t.Db.WithContext(ctx).Where("task_id = ? and task_run_id ", taskID, taskRunID).Updates(map[string]string{
		"finished_at": "now()",
	}).Error
}

func (t TaskRunUserRepo) GetTaskRunUserList(ctx context.Context, taskID int) ([]entity.TaskRunUser, error) {
	var taskRunUsers []entity.TaskRunUser
	err := t.Db.WithContext(ctx).Where("task_id = ?", taskID).Find(&taskRunUsers).Error
	return taskRunUsers, err
}
