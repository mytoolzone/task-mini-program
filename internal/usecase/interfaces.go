package usecase

import (
	"context"

	"github.com/mytoolzone/task-mini-program/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// Translation -.
	Translation interface {
		Translate(context.Context, entity.Translation) (entity.Translation, error)
		History(context.Context) ([]entity.Translation, error)
	}

	// TranslationRepo -.
	TranslationRepo interface {
		Store(context.Context, entity.Translation) error
		GetHistory(context.Context) ([]entity.Translation, error)
	}

	// TranslationWebAPI -.
	TranslationWebAPI interface {
		Translate(entity.Translation) (entity.Translation, error)
	}

	// User -.
	User interface {
		Login(context.Context, entity.User) (entity.User, error)
		Register(context.Context, entity.User) (entity.User, error)
		UpdateSetting(ctx context.Context, userID int, setting entity.UserSetting) error
		GetSettingByUserID(ctx context.Context, userID int) (entity.UserSetting, error)
	}

	// UserRepo -.
	UserRepo interface {
		Store(context.Context, entity.User) error
		GetByUserID(ctx context.Context, userID int) (entity.User, error)
		GetUserSettingByUserID(ctx context.Context, userID int) (entity.UserSetting, error)
		UpdateUserSetting(ctx context.Context, userID int, setting entity.UserSetting) error
	}

	// Task -.
	// 任务列表 -> 任务详情 -> 开始签到 -> 签到 -> 开始子任务执行 -> 上传子任务记录 -> 暂停任务执行 -> 子任务完成
	// 开始签到 -> 签到 -> 开始子任务执行 -> 上传子任务记录 -> 完成任务 -> 子任务完成
	// -> 完成任务
	Task interface {
		CreateTask(context.Context, entity.Task) error
		// GetByUserID 获取一个人参与的任务
		GetByUserID(ctx context.Context, userID int) ([]entity.Task, error)
		// GetByTaskID 获取一个任务的详情
		GetByTaskID(ctx context.Context, taskID int) (entity.Task, error)
		// GetTaskList 任务大厅获取任务列表
		GetTaskList(ctx context.Context, lastId int) ([]entity.Task, error)
		// GetTaskUsers 获取任务参与者列表
		GetTaskUsers(ctx context.Context, taskID int) ([]entity.TaskUser, error)
		// GetTaskRunList 获取某个任务的子任务列表
		GetTaskRunList(ctx context.Context, taskID int) ([]entity.TaskRun, error)
		// GetTaskRunLogList 获取某个任务的记录员上传的任务记录
		GetTaskRunLogList(ctx context.Context, taskID int) ([]entity.TaskRunLog, error)
		// PrepareTaskRun 准备开始子任务 返回子任务id 后续签到生产二维码使用 ,队长点击开始签到调用这个接口
		// 1. 如果是第一次调用,则创建子任务
		// 2. 如果子任务没有完成 返回未完成子任务
		PrepareTaskRun(ctx context.Context, taskID int) (int, error)
		// Sign 签到
		Sign(ctx context.Context, taskID, taskRunID, userID int) error
		// StartTaskRun 开始子任务执行,如果没有自动生产一个子任务记录
		StartTaskRun(ctx context.Context, taskID int) error
		// PauseTaskRun 暂停任务执行
		PauseTaskRun(ctx context.Context, taskID int) error
		// FinishTaskRun 完成子任务执行
		FinishTaskRun(ctx context.Context, taskID int) error
	}

	// TaskRepo -.
	TaskRepo interface {
		CreateTask(context.Context, *entity.Task) error
		GetByUserID(ctx context.Context, userID int) ([]entity.Task, error)
		GetByTaskID(ctx context.Context, taskID int) (entity.Task, error)
		GetTaskList(ctx context.Context, lastId int) ([]entity.Task, error)

		AuditFailTask(ctx context.Context, taskID int) error
		AuditSuccessTask(ctx context.Context, taskID int) error
		StartTask(ctx context.Context, taskID int) error
		PauseTask(ctx context.Context, taskID int) error
		FinishTask(ctx context.Context, taskID int) error
	}

	// TaskRunRepo -.
	// 子任务记录
	TaskRunRepo interface {
		GetPendingTaskRun(ctx context.Context, taskID int) (entity.TaskRun, error)
		GetRunningTaskRun(ctx context.Context, taskID int) (entity.TaskRun, error)
		AddTaskRun(ctx context.Context, taskID int) (entity.TaskRun, error)
		StartTaskRun(ctx context.Context, taskID int) error
		FinishTaskRun(ctx context.Context, taskID int) error
		GetTaskRunList(ctx context.Context, taskID int) ([]entity.TaskRun, error)
	}

	// TaskRunLogRepo -.
	// 记录员上报的任务记录
	TaskRunLogRepo interface {
		AddTaskRunLog(ctx context.Context, log *entity.TaskRunLog) error
		GetTaskRunLogList(ctx context.Context, taskID int) ([]entity.TaskRunLog, error)
	}

	// TaskRunUserRepo -.
	// 任务执行人记录 执行任务详细记录
	TaskRunUserRepo interface {
		// AddTaskRunUser 签到时候调用
		AddTaskRunUser(ctx context.Context, taskID, taskRunID, userID int) (entity.TaskRunUser, error)
		// StartTaskRun 开始子任务执行时候调用
		StartTaskRun(ctx context.Context, taskID, taskRunID int) error
		// FinishTaskRun 完成子任务执行时候调用
		FinishTaskRun(ctx context.Context, taskID, taskRunID int) error
		GetTaskRunUserList(ctx context.Context, taskID int) ([]entity.TaskRunUser, error)
	}

	// TaskUserRepo -.
	// 任务参与者
	TaskUserRepo interface {
		// AddTaskUser 用户报名参与任务
		AddTaskUser(ctx context.Context, taskID, userID int) (entity.TaskUser, error)
		// AuditTaskUser 审核任务参与者
		AuditTaskUser(ctx context.Context, taskID, userID int, status string) (entity.TaskUser, error)
		// GetTaskUserList 获取任务参与者列表
		GetTaskUserList(ctx context.Context, taskID int) ([]entity.TaskUser, error)
		// GetTaskUserByUserID 获取任务参与者状态
		GetTaskUserByUserID(ctx context.Context, taskID, userID int) (entity.TaskUser, error)
	}

	// Notice -.
	Notice interface {
		// AddNotice 添加通知
		AddNotice(ctx context.Context, notice entity.Notice) error
		// GetNoticeListByUser 获取某个用户未读通知列表
		GetNoticeListByUser(ctx context.Context, userID int) ([]entity.Notice, error)
		// GetNoticeByNoticeID 获取通知详情 - 更新阅读通知状态
		GetNoticeByNoticeID(ctx context.Context, noticeID int) (entity.Notice, error)
	}

	// NoticeRepo -.
	NoticeRepo interface {
		// AddNotice 添加通知
		AddNotice(ctx context.Context, notice entity.Notice) error
		// GetNoticeListByUser 获取某个用户未读通知列表
		GetNoticeListByUser(ctx context.Context, userID int) ([]entity.Notice, error)
		// GetNoticeByNoticeID 获取通知详情 - 更新阅读通知状态
		GetNoticeByNoticeID(ctx context.Context, noticeID int) (entity.Notice, error)
	}
)
