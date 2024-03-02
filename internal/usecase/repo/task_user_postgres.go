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
	if role != entity.UserTaskRoleRecorder && role != entity.UserTaskRoleLeader && role != entity.UserTaskRoleMember {
		return app_code.New(app_code.ErrorAuditParamInValid, "role invalid")
	}

	var userTask = entity.UserTask{}
	if err := t.Db.WithContext(ctx).Where("task_id = ? and user_id = ?", taskID, userID).First(&userTask).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	userTask.Role = role
	userTask.TaskID = taskID
	userTask.UserID = userID

	if err := t.Db.WithContext(ctx).Where("task_id = ? and user_id = ?", taskID, userID).Save(userTask).Error; err != nil {
		return err
	}

	return nil
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
		UserTask.Status = entity.UserTaskStatusNotApply
		return entity.UserTask{}, gorm.ErrRecordNotFound
	}
	return UserTask, err
}

// GetUserJoinTaskList 获取某个人参与的任务列表
func (t UserTaskRepo) GetUserJoinTaskList(ctx context.Context, userID int, status string, lastID int) ([]entity.UserTask, error) {
	var tasks []entity.UserTask
	query := t.Db.WithContext(ctx).Where("user_id = ?", userID)
	if lastID > 0 {
		query = query.Where("id < ?", lastID)
	}
	if status != "" {
		query = query.Where("status = ?)", status)
	}
	err := query.Debug().Preload("Task").Find(&tasks).Error
	return tasks, err
}
