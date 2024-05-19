package repo

import (
	"context"
	"time"

	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/pkg/postgres"
	"gorm.io/gorm"
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
	err := t.Db.WithContext(ctx).Where("task_id = ? and task_run_id = ?", taskID, taskRunID).Find(&taskRunUsers).Error
	return taskRunUsers, err
}

func (t TaskRunUserRepo) GetUserTaskSummary(ctx context.Context, userID int, startTime, endTime string, taskID int) (entity.UserTaskSummary, error) {

	// 获取用户参与完成的任务总数
	// err := t.Db.WithContext(ctx).Model(&entity.UserTask{}).Where("user_id = ?", userID).Count(&userTaskSummary.TotalTask).Error
	// if err != nil {
	// 	return entity.UserTaskSummary{}, err
	// }
	var userTaskSummary entity.UserTaskSummary
	var query, query1, query2 *gorm.DB
	var err error
	if startTime != "" && endTime != "" {
		query = t.Db.WithContext(ctx).Debug().Model(&entity.TaskRunUser{}).
			Where("user_id = ? and status = ? and start_at >= ? and finished_at <= ?", userID, entity.TaskStatusFinished, startTime, endTime)
	} else {
		// 先查完成任务总数是否大于0
		query = t.Db.WithContext(ctx).Debug().Model(&entity.TaskRunUser{}).
			Where("user_id = ? and status = ?", userID, entity.TaskStatusFinished)
	}
	if taskID > 0 {
		query = query.Where("task_id = ?", taskID)
	}
	err = query.Select("count(DISTINCT task_id)").Scan(&userTaskSummary.TotalTask).Error
	if err != nil {
		return entity.UserTaskSummary{}, err
	}
	if userTaskSummary.TotalTask > 0 {
		if startTime != "" && endTime != "" {
			// 查询所有任务的总时长
			query1 = t.Db.WithContext(ctx).Debug().Model(&entity.TaskRunUser{}).
				Where("user_id = ? and status = ? and start_at >= ? and finished_at <= ?", userID, entity.TaskStatusFinished, startTime, endTime)
			if taskID > 0 {
				query1 = query1.Where("task_id = ?", taskID)
			}
			err = query1.
				Select("sum(duration) as total_duration").Scan(&userTaskSummary.TotalDuration).Error
			if err != nil {
				return entity.UserTaskSummary{}, err
			}
			// 分组查询每个任务的总时长
			query2 = t.Db.WithContext(ctx).Debug().Model(&entity.TaskRunUser{}).
				Joins("INNER JOIN tasks ON tasks.id = task_run_users.task_id").
				Where("task_run_users.user_id = ? and task_run_users.status = ? and task_run_users.start_at >= ? and task_run_users.finished_at <= ?", userID, entity.TaskStatusFinished, startTime, endTime)
			if taskID > 0 {
				query2 = query2.Where("task_run_users.task_id = ?", taskID)
			}
			err = query2.
				Group("task_run_users.task_id").
				Select("sum(duration) as total_duration,task_id").Scan(&userTaskSummary.UserTaskSummaryList).Error

			if err != nil {
				return entity.UserTaskSummary{}, err
			}
		} else {
			// 查总任务时长
			query1 = t.Db.WithContext(ctx).Debug().Model(&entity.TaskRunUser{}).
				Where("user_id = ? and status = ?", userID, entity.TaskStatusFinished)
			if taskID > 0 {
				query1 = query1.Where("task_id = ?", taskID)
			}
			err = query1.
				Select("sum(duration) as total_duration").Scan(&userTaskSummary.TotalDuration).Error
			if err != nil {
				return entity.UserTaskSummary{}, err
			}
			// 分组查询每个任务的总时长
			query2 = t.Db.WithContext(ctx).Debug().Model(&entity.TaskRunUser{}).
				Joins("INNER JOIN tasks ON tasks.id = task_run_users.task_id").
				Where("task_run_users.user_id = ? and task_run_users.status = ?", userID, entity.TaskStatusFinished)
			if taskID > 0 {
				query2 = query2.Where("task_run_users.task_id = ?", taskID)
			}
			err = query2.
				Group("task_run_users.task_id").
				Select("sum(duration) as total_duration,task_id").Scan(&userTaskSummary.UserTaskSummaryList).Error

			if err != nil {
				return entity.UserTaskSummary{}, err
			}
		}
	}
	return userTaskSummary, nil
}
