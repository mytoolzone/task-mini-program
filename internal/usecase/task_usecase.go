package usecase

import (
	"context"
	"errors"
	"github.com/mytoolzone/task-mini-program/internal/app_code"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"gorm.io/gorm"
)

type TaskUseCase struct {
	t   TaskRepo
	tr  TaskRunRepo
	tru TaskRunUserRepo
	trl TaskRunLogRepo
	tu  UserTaskRepo
}

func NewTaskUseCase(t TaskRepo, tr TaskRunRepo, tru TaskRunUserRepo, trl TaskRunLogRepo, tu UserTaskRepo) *TaskUseCase {
	return &TaskUseCase{
		t:   t,
		tr:  tr,
		tru: tru,
		trl: trl,
		tu:  tu,
	}
}

func (t TaskUseCase) CreateTask(ctx context.Context, task entity.Task) error {
	return t.t.CreateTask(ctx, &task)
}

func (t TaskUseCase) GetTaskDetail(ctx context.Context, taskID int) (entity.Task, error) {
	return t.t.GetByTaskID(ctx, taskID)
}

func (t TaskUseCase) GetByUserID(ctx context.Context, userID int) ([]entity.Task, error) {
	return t.t.GetByUserID(ctx, userID)
}

func (t TaskUseCase) GetByTaskID(ctx context.Context, taskID int) (entity.Task, error) {
	return t.GetByTaskID(ctx, taskID)
}

func (t TaskUseCase) GetTaskList(ctx context.Context, lastId int) ([]entity.Task, error) {
	return t.GetTaskList(ctx, lastId)
}

func (t TaskUseCase) GetUserTasks(ctx context.Context, taskID int) ([]entity.UserTask, error) {
	return t.tu.GetUserTaskList(ctx, taskID)
}

func (t TaskUseCase) GetTaskRunList(ctx context.Context, taskID int) ([]entity.TaskRun, error) {
	return t.tr.GetTaskRunList(ctx, taskID)
}

func (t TaskUseCase) GetTaskRunLogList(ctx context.Context, taskID int) ([]entity.TaskRunLog, error) {
	return t.trl.GetTaskRunLogList(ctx, taskID)
}

func (t TaskUseCase) UploadRunLog(ctx context.Context, runLog entity.TaskRunLog) error {
	return t.trl.AddTaskRunLog(ctx, &runLog)
}

func (t TaskUseCase) PrepareTaskRun(ctx context.Context, taskID int) (int, error) {
	// 0. 判断任务是否已经开始
	run, err := t.tr.GetPendingTaskRun(ctx, taskID)
	if err != nil {
		return 0, err
	}
	if run.ID != 0 {
		return run.ID, nil
	}
	// 1. 如果是第一次调用,则创建子任务
	taskRun, err := t.tr.AddTaskRun(ctx, taskID)
	if err != nil {
		return 0, err
	}
	return taskRun.ID, nil
}

// Sign 签到 生产子任务id , 扫描时候使用
func (t TaskUseCase) Sign(ctx context.Context, taskID, taskRunID, userID int) error {
	_, err := t.tru.AddTaskRunUser(ctx, taskID, taskRunID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (t TaskUseCase) StartTaskRun(ctx context.Context, taskID int) error {
	run, err := t.tr.GetPendingTaskRun(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app_code.New(app_code.ErrorTaskRunNotFound, "没有准备开始的任务,任务可能已经开始或者未签到完成")
		}
		return err
	}

	err = t.t.StartTask(ctx, taskID)
	if err != nil {
		return err
	}

	if err := t.tr.StartTaskRun(ctx, taskID); err != nil {
		return err
	}

	return t.tru.StartTaskRun(ctx, taskID, run.ID)
}

func (t TaskUseCase) PauseTaskRun(ctx context.Context, taskID int) error {
	run, err := t.tr.GetRunningTaskRun(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app_code.New(app_code.ErrorTaskRunNotFound, "没有正在执行中任务")
		}
		return err
	}

	err = t.t.PauseTask(ctx, taskID)
	if err != nil {
		return err
	}

	if err := t.tr.FinishTaskRun(ctx, taskID); err != nil {
		return err
	}

	return t.tru.FinishTaskRun(ctx, taskID, run.ID)
}

func (t TaskUseCase) FinishTaskRun(ctx context.Context, taskID int) error {
	task, err := t.t.GetByTaskID(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app_code.New(app_code.ErrorTaskRunNotFound, "任务不存在")
		}
		return err
	}

	if task.Status == entity.TaskStatusRunning {
		run, err := t.tr.GetRunningTaskRun(ctx, taskID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return app_code.New(app_code.ErrorTaskRunNotFound, "没有正在执行中任务")
			}
			return err
		}

		if err := t.tr.FinishTaskRun(ctx, taskID); err != nil {
			return err
		}
		err = t.tru.FinishTaskRun(ctx, taskID, run.ID)
		if err != nil {
			return err
		}
	}

	err = t.t.FinishTask(ctx, taskID)
	if err != nil {
		return err
	}

	return nil
}

func (t TaskUseCase) CancelTaskRun(ctx context.Context, taskID int) error {
	task, err := t.t.GetByTaskID(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app_code.New(app_code.ErrorTaskRunNotFound, "任务不存在")
		}
		return err
	}

	if task.Status == entity.TaskStatusRunning {
		run, err := t.tr.GetRunningTaskRun(ctx, taskID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return app_code.New(app_code.ErrorTaskRunNotFound, "没有正在执行中任务")
			}
			return err
		}

		if err := t.tr.CancelTaskRun(ctx, taskID); err != nil {
			return err
		}
		err = t.tru.CancelTaskRun(ctx, taskID, run.ID)
		if err != nil {
			return err
		}
	}

	err = t.t.CancelTask(ctx, taskID)
	if err != nil {
		return err
	}

	return nil
}

func (t TaskUseCase) AuditTask(ctx context.Context, taskID int, status string) error {
	var err error
	switch status {
	case entity.TaskStatusAuditFail:
		err = t.t.AuditFailTask(ctx, taskID)

	case entity.TaskStatusAuditPass:
		err = t.t.AuditSuccessTask(ctx, taskID)

	}
	if err != nil {
		return err
	}
	return nil
}

func (t TaskUseCase) JoinTask(ctx context.Context, taskID, userID int) error {
	_, err := t.tu.AddUserTask(ctx, taskID, userID)
	return err
}

func (t TaskUseCase) AuditUserTask(ctx context.Context, taskID, userID int, status string) error {
	var err error
	switch status {
	case entity.UserTaskStatusAuditFail:
		_, err = t.tu.AuditUserTask(ctx, taskID, userID, status)

	case entity.UserTaskStatusAuditPass:
		_, err = t.tu.AuditUserTask(ctx, taskID, userID, status)

	}
	if err != nil {
		return err
	}
	return nil
}

func (t TaskUseCase) GetUserTaskSummary(ctx context.Context, userID int) (entity.UserTaskSummary, error) {
	return t.tru.GetUserTaskSummary(ctx, userID)
}
