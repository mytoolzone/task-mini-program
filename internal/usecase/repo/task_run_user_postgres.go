package repo

import (
	"context"
	"strconv"
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
		Duration:   0,
	}
	// 更新task_run_user表中 finished_at 为当前时间，同时计算有效工时分钟数
	var taskRunUser entity.TaskRunUser
	if err := t.Db.Debug().WithContext(ctx).Model(&entity.TaskRunUser{}).Where("task_id = ? and task_run_id = ?", taskID, taskRunID).First(&taskRunUser).Error; err != nil {
		return err
	}
	tru.Duration = int(tru.FinishedAt.Sub(taskRunUser.StartAt).Minutes())
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

func (t TaskRunUserRepo) GetUserTaskSummary(ctx context.Context, userID int, startTime, endTime string, taskID int, taskName, is_group_user string, page, page_size string) (entity.UserTaskSummary, error) {

	var userTaskSummary entity.UserTaskSummary
	var query, query1, query2 *gorm.DB
	var err error
	// 默认当天时间
	if startTime == "" {
		startTime = time.Now().Format("2006-01-02") + " 00:00:00"
	}
	if endTime == "" {
		endTime = time.Now().Format("2006-01-02") + " 23:59:59"
	}
	// 任务标题，用户ID，开始时间，结束时间
	// 查工时总计，任务列表，每个任务的总工时，开始结束时间
	query = t.Db.WithContext(ctx).Debug().Model(&entity.TaskRunUser{})
	if taskName != "" { //任务标题
		query = query.Joins("INNER JOIN tasks ON tasks.id = task_run_users.task_id").Where("tasks.name like ?", "%"+taskName+"%")
	}
	if startTime != "" && endTime != "" { //任务开始结束时间
		query = query.Where("task_run_users.status = ? and task_run_users.start_at >= ? and task_run_users.finished_at <= ?", entity.TaskStatusFinished, startTime, endTime)
	} else {
		// 先查完成任务总数是否大于0
		query = query.
			Where("status = ?", entity.TaskStatusFinished)
	}
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if taskID > 0 {
		query = query.Where("task_id = ?", taskID)
	}
	err = query.Select("count(DISTINCT task_id)").Scan(&userTaskSummary.TotalTask).Error
	if err != nil {
		return entity.UserTaskSummary{}, err
	}
	if userTaskSummary.TotalTask > 0 {
		query1 = t.Db.WithContext(ctx).Debug().Model(&entity.TaskRunUser{})
		if taskName != "" { //任务标题
			query1 = query1.Joins("INNER JOIN tasks ON tasks.id = task_run_users.task_id").Where("tasks.name like ?", "%"+taskName+"%")
		}
		query2 = t.Db.WithContext(ctx).Debug().Model(&entity.TaskRunUser{}).Joins("INNER JOIN tasks ON tasks.id = task_run_users.task_id")
		if taskName != "" {
			query2 = query2.Where("tasks.name like ?", "%"+taskName+"%")
		}
		if startTime != "" && endTime != "" {
			// 查询所有任务的总时长
			query1 = query1.
				Where("task_run_users.status = ? and task_run_users.start_at >= ? and task_run_users.finished_at <= ?", entity.TaskStatusFinished, startTime, endTime)
			if taskID > 0 {
				query1 = query1.Where("task_id = ?", taskID)
			}
			if userID > 0 {
				query1 = query1.Where("user_id = ?", userID)
			}
			err = query1.
				Select("sum(duration) as total_duration").Scan(&userTaskSummary.TotalDuration).Error
			if err != nil {
				return entity.UserTaskSummary{}, err
			}
			// 分组查询每个任务的总时长
			query2 = query2.Where("task_run_users.status = ? and task_run_users.start_at >= ? and task_run_users.finished_at <= ?", entity.TaskStatusFinished, startTime, endTime)
			if taskID > 0 {
				query2 = query2.Where("task_run_users.task_id = ?", taskID)
			}
			if userID > 0 { //用户ID的时候默认按照用户ID分组
				query2 = query2.Where("task_run_users.user_id = ?", userID).
					Group("task_run_users.task_id,user_id").
					Select("sum(duration) as total_duration,task_id,user_id")
			} else {
				if is_group_user == "1" {
					query2 = query2.
						Group("task_run_users.task_id,user_id").
						Select("sum(duration) as total_duration,task_id,user_id")
				} else {
					query2 = query2.
						Group("task_run_users.task_id").
						Select("sum(duration) as total_duration,task_id")
				}
			}
			// id降序，分页每页20条,默认查第一页
			if page != "" && page_size != "" {
				// 先将page和page_size强制转int
				page, _ := strconv.Atoi(page)
				page_size, _ := strconv.Atoi(page_size)
				if page > 0 && page_size > 0 {
					query2.Order("task_run_users.task_id desc").Limit(page_size).Offset((page - 1) * page_size)
				} else {
					query2.Order("task_run_users.task_id desc").Limit(20)
				}
			} else {
				query2.Order("task_run_users.task_id desc").Limit(20)
			}
			err = query2.Scan(&userTaskSummary.UserTaskSummaryList).Error

			if err != nil {
				return entity.UserTaskSummary{}, err
			}
		} else {
			// 查总任务时长
			query1 = query1.
				Where("status = ?", userID, entity.TaskStatusFinished)
			if userID > 0 {
				query1 = query1.Where("user_id = ?", userID)
			}
			if taskID > 0 {
				query1 = query1.Where("task_id = ?", taskID)
			}
			err = query1.
				Select("sum(duration) as total_duration").Scan(&userTaskSummary.TotalDuration).Error
			if err != nil {
				return entity.UserTaskSummary{}, err
			}
			// 分组查询每个任务的总时长
			query2 = query2.
				Where("task_run_users.user_id = ? and task_run_users.status = ?", userID, entity.TaskStatusFinished)
			if taskID > 0 {
				query2 = query2.Where("task_run_users.task_id = ?", taskID)
			}
			if userID > 0 {
				query2 = query2.Where("task_run_users.user_id = ?", userID).
					Group("task_run_users.task_id,user_id").
					Select("sum(duration) as total_duration,task_id,user_id")
			} else {
				if is_group_user == "1" {
					query2 = query2.
						Group("task_run_users.task_id,user_id").
						Select("sum(duration) as total_duration,task_id,user_id")
				} else {
					query2 = query2.
						Group("task_run_users.task_id").
						Select("sum(duration) as total_duration,task_id")
				}
			}
			err = query2.Scan(&userTaskSummary.UserTaskSummaryList).Error

			if err != nil {
				return entity.UserTaskSummary{}, err
			}
		}
	}
	return userTaskSummary, nil
}

// 查询任务工时详情
func (t TaskRunUserRepo) GetUserTaskSummaryDetail(ctx context.Context, userID int, startTime, endTime string, taskID int, taskName, is_group_user string) (entity.UserTaskSummary, error) {
	var userTaskSummary entity.UserTaskSummary
	var query, query1, query2 *gorm.DB
	var err error
	// 任务标题，用户ID，开始时间，结束时间
	// 查工时总计，任务列表，每个任务的总工时，开始结束时间
	query = t.Db.WithContext(ctx).Debug().Model(&entity.TaskRunUser{})
	// if taskName != "" { //任务标题
	// 	query = query.Joins("INNER JOIN tasks ON tasks.id = task_run_users.task_id").Where("tasks.name like ?", "%"+taskName+"%")
	// }
	if startTime != "" && endTime != "" { //任务开始结束时间
		query = query.Where("task_run_users.status = ? and task_run_users.start_at >= ? and task_run_users.finished_at <= ?", entity.TaskStatusFinished, startTime, endTime)
	} else {
		// 先查完成任务总数是否大于0
		query = query.
			Where("status = ?", entity.TaskStatusFinished)
	}
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if taskID > 0 {
		query = query.Where("task_id = ?", taskID)
	}
	err = query.Select("count(DISTINCT task_id)").Scan(&userTaskSummary.TotalTask).Error
	if err != nil {
		return entity.UserTaskSummary{}, err
	}
	query1 = t.Db.WithContext(ctx).Debug().Model(&entity.TaskRunUser{})
	query2 = t.Db.WithContext(ctx).Debug().Model(&entity.TaskRunUser{}).Joins("INNER JOIN tasks ON tasks.id = task_run_users.task_id")
	if userTaskSummary.TotalTask > 0 {
		if taskName != "" { //任务标题
			query1 = query1.Joins("INNER JOIN tasks ON tasks.id = task_run_users.task_id").Where("tasks.name like ?", "%"+taskName+"%")
		}

		if taskName != "" {
			query2 = query2.Where("tasks.name like ?", "%"+taskName+"%")
		}
		if startTime != "" && endTime != "" {
			// 查询所有任务的总时长
			query1 = query1.
				Where("task_run_users.status = ? and task_run_users.start_at >= ? and task_run_users.finished_at <= ?", entity.TaskStatusFinished, startTime, endTime)
			if taskID > 0 {
				query1 = query1.Where("task_id = ?", taskID)
			}
			if userID > 0 {
				query1 = query1.Where("user_id = ?", userID)
			}
			err = query1.
				Select("sum(duration) as total_duration").Scan(&userTaskSummary.TotalDuration).Error
			if err != nil {
				return entity.UserTaskSummary{}, err
			}
			// 分组查询每个任务的总时长
			query2 = query2.Where("task_run_users.status = ? and task_run_users.start_at >= ? and task_run_users.finished_at <= ?", entity.TaskStatusFinished, startTime, endTime)
			if taskID > 0 {
				query2 = query2.Where("task_run_users.task_id = ?", taskID)
			}
			if userID > 0 { //用户ID的时候默认按照用户ID分组
				query2 = query2.Where("task_run_users.user_id = ?", userID).
					Group("task_run_users.task_id,user_id").
					Select("sum(duration) as total_duration,task_id,user_id")
			} else {
				if is_group_user == "1" {
					query2 = query2.
						Group("task_run_users.task_id,user_id").
						Select("sum(duration) as total_duration,task_id,user_id")
				} else {
					query2 = query2.
						Group("task_run_users.task_id").
						Select("sum(duration) as total_duration,task_id")
				}
			}
			err = query2.Scan(&userTaskSummary.UserTaskSummaryList).Error

			if err != nil {
				return entity.UserTaskSummary{}, err
			}
		} else {
			// 查总任务时长
			query1 = query1.
				Where("status = ?", entity.TaskStatusFinished)
			if userID > 0 {
				query1 = query1.Where("user_id = ?", userID)
			}
			if taskID > 0 {
				query1 = query1.Where("task_id = ?", taskID)
			}
			err = query1.
				Select("sum(duration) as total_duration").Scan(&userTaskSummary.TotalDuration).Error
			if err != nil {
				return entity.UserTaskSummary{}, err
			}
			// 分组查询每个任务的总时长
			query2 = query2.
				Where("task_run_users.status = ?", entity.TaskStatusFinished)
			if taskID > 0 {
				query2 = query2.Where("task_run_users.task_id = ?", taskID)
			}
			// if userID > 0 {
			// 	query2 = query2.Where("task_run_users.user_id = ?", userID).
			// 		Group("task_run_users.task_id,user_id").
			// 		Select("sum(duration) as total_duration,task_id,user_id")
			// }
			err = query2.
				Group("task_run_users.task_id,user_id").
				Select("sum(duration) as total_duration,task_id,user_id").
				Scan(&userTaskSummary.UserTaskSummaryList).Error
			if err != nil {
				return entity.UserTaskSummary{}, err
			}
		}
	}
	return userTaskSummary, nil
}
