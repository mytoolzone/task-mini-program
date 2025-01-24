package repo

import (
	"context"
	"errors"

	"github.com/mytoolzone/task-mini-program/internal/app_code"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/pkg/postgres"
	"gorm.io/gorm"
)

type UserTaskRepo struct {
	*postgres.Postgres
}

func NewUserTaskRepo(pg *postgres.Postgres) *UserTaskRepo {
	return &UserTaskRepo{pg}
}

func (t UserTaskRepo) AddUserTask(ctx context.Context, taskID, userID int) (entity.UserTask, error) {
	//
	var UserTask entity.UserTask = entity.UserTask{
		TaskID: taskID,
		UserID: userID,
		Status: entity.UserTaskStatusApply,
		Role:   entity.UserRoleMember,
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

// AssignRole 分配角色
func (t UserTaskRepo) AssignRole(ctx context.Context, taskID, userID int, role string) error {
	if role != entity.UserTaskRoleRecorder && role != entity.UserTaskRoleLeader && role != entity.UserTaskRoleMember && role != entity.UserTaskRoleSafer {
		return app_code.New(app_code.ErrorBadRequest, "role invalid")
	}

	var userTask = entity.UserTask{}
	if err := t.Db.WithContext(ctx).Debug().Where("task_id = ? and user_id = ?", taskID, userID).First(&userTask).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app_code.New(app_code.ErrorUserTaskNotFound, "user task not found")
		}
		return err
	}
	// 需要前端校验
	/* if userTask.Role == entity.UserTaskRoleLeader { //队长
		return app_code.New(app_code.ErrorForbidden, "指挥官不允许变更角色")
	} */
	userTask.Role = role
	userTask.TaskID = taskID
	userTask.UserID = userID

	if role == entity.UserTaskRoleLeader { //队长
		leader, ok, err := t.GetTaskLeader(ctx, taskID)
		if err != nil {
			return err
		}
		if ok { //存在更新为队员
			leader.Role = entity.UserRoleMember
			if err := t.Db.WithContext(ctx).Save(&leader).Error; err != nil {
				return err
			}
		}
	}

	if role == entity.UserTaskRoleRecorder { //记录员
		recorder, ok, err := t.GetTaskRecorder(ctx, taskID)
		if err != nil {
			return err
		}
		if ok {
			recorder.Role = entity.UserRoleMember
			if err := t.Db.WithContext(ctx).Save(&recorder).Error; err != nil {
				return err
			}

		}
	}

	if role == entity.UserTaskRoleSafer { //安全员
		safer, ok, err := t.GetTaskSafer(ctx, taskID)
		if err != nil {
			return err
		}
		if ok {
			safer.Role = entity.UserRoleMember
			if err := t.Db.WithContext(ctx).Save(&safer).Error; err != nil {
				return err
			}

		}
	}

	if err := t.Db.WithContext(ctx).Where("task_id = ? and user_id = ?", taskID, userID).Save(userTask).Error; err != nil {
		return err
	}

	return nil
}

// GetTaskLeader 获取任务中的队长
func (t UserTaskRepo) GetTaskLeader(ctx context.Context, taskID int) (entity.UserTask, bool, error) {
	var userTask entity.UserTask
	err := t.Db.WithContext(ctx).Where("task_id = ? and role = ?", taskID, entity.UserTaskRoleLeader).First(&userTask).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.UserTask{}, false, nil
	}
	if err != nil {
		return entity.UserTask{}, false, err
	}
	return userTask, true, nil
}

// GetTaskRecorder 获取任务中的记录员
func (t UserTaskRepo) GetTaskRecorder(ctx context.Context, taskID int) (entity.UserTask, bool, error) {
	var userTask entity.UserTask
	err := t.Db.WithContext(ctx).Where("task_id = ? and role = ?", taskID, entity.UserTaskRoleRecorder).First(&userTask).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.UserTask{}, false, nil
	}
	if err != nil {
		return entity.UserTask{}, false, err
	}
	return userTask, true, nil
}

// GetTaskSafer 获取任务中的安全员
func (t UserTaskRepo) GetTaskSafer(ctx context.Context, taskID int) (entity.UserTask, bool, error) {
	var userTask entity.UserTask
	err := t.Db.WithContext(ctx).Where("task_id = ? and role = ?", taskID, entity.UserTaskRoleSafer).First(&userTask).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.UserTask{}, false, nil
	}
	if err != nil {
		return entity.UserTask{}, false, err
	}
	return userTask, true, nil
}

func (t UserTaskRepo) GetTaskUserList(ctx context.Context, taskID int, status string) ([]entity.UserTask, error) {
	var UserTasks []entity.UserTask
	query := t.Db.WithContext(ctx).Where("task_id = ?", taskID)
	if status != "" {
		query = query.Where("status =?", status)
	}
	err := query.Debug().Preload("User").Find(&UserTasks).Error
	return UserTasks, err
}

// GetUserTaskByUserID 获取某个人参与的任务的状态
func (t UserTaskRepo) GetUserTaskByUserID(ctx context.Context, taskID, userID int) (entity.UserTask, error) {
	var UserTask entity.UserTask
	err := t.Db.WithContext(ctx).Where("task_id = ? and user_id = ?", taskID, userID).First(&UserTask).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.UserTask{}, nil
	}
	if err := t.Db.WithContext(ctx).Debug().Where("id = ?", userID).First(&UserTask.User).Error; err != nil {
		return entity.UserTask{}, err
	}
	return UserTask, err
}

// GetUserJoinTaskList 获取某个人参与的任务列表
func (t UserTaskRepo) GetUserJoinTaskList(ctx context.Context, userID int, status string, lastID int) (*entity.UserTaskMap, error) {
	var tasks []entity.UserTask
	query := t.Db.WithContext(ctx).Model(&entity.UserTask{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	var count int64
	err := query.Debug().Count(&count).Error
	if err != nil {
		return &entity.UserTaskMap{}, err
	}
	if lastID > 0 {
		query = query.Where("id < ?", lastID)
	}
	// 根据lastID分页
	query = query.Order("id DESC").Limit(50)
	err = query.Debug().Preload("Task").Find(&tasks).Error
	mapRes := entity.UserTaskMap{
		"count":    count,
		"TaskList": tasks,
	}
	return &mapRes, err
}
