package usecase

import (
	"context"
	"errors"

	"github.com/gw123/glog"
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

func (t TaskUseCase) CreateTask(ctx context.Context, task *entity.Task) error {
	return t.t.CreateTask(ctx, task)
}

func (t TaskUseCase) GetTaskDetail(ctx context.Context, taskID int) (entity.Task, error) {
	return t.t.GetByTaskID(ctx, taskID)
}

func (t TaskUseCase) GetByUserID(ctx context.Context, userID int, status string, lastID int) ([]entity.Task, error) {
	return t.t.GetByUserID(ctx, userID, status, lastID)
}

func (t TaskUseCase) GetByTaskID(ctx context.Context, taskID int) (entity.Task, error) {
	return t.t.GetByTaskID(ctx, taskID)
}

func (t TaskUseCase) GetTaskList(ctx context.Context, lastID int, keyword, status string) ([]entity.Task, error) {
	return t.t.GetTaskList(ctx, lastID, keyword, status)
}

func (t TaskUseCase) GetTaskUsers(ctx context.Context, taskID int, status string) ([]entity.UserTask, error) {
	taskUserList, err := t.tu.GetTaskUserList(ctx, taskID, status)
	if err != nil {
		return nil, err
	}
	return taskUserList, err
}

func (t TaskUseCase) GetApprovedTaskUsers(ctx context.Context, taskID int) ([]entity.UserTask, error) {
	taskUserList, err := t.tu.GetTaskUserList(ctx, taskID, entity.UserTaskStatusAuditPass)
	if err != nil {
		return nil, err
	}

	tr, err := t.tr.GetTaskLatestRun(ctx, taskID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	for index, _ := range taskUserList {
		taskUserList[index].Status = entity.TaskStatusNotSign
	}

	if err == gorm.ErrRecordNotFound {
		return taskUserList, nil
	}

	taskRunUserList, err := t.tru.GetTaskRunUserList(ctx, taskID, tr.ID)
	if err != nil {
		return nil, err
	}

	for index, _ := range taskUserList {
		for _, tru := range taskRunUserList {
			if taskUserList[index].UserID == tru.UserID {
				taskUserList[index].Status = entity.TaskStatusSign
			}
		}
	}

	return taskUserList, nil
}

func (t TaskUseCase) GetUserTaskRole(ctx context.Context, taskID, userID int) (entity.UserTask, error) {
	return t.tu.GetUserTaskByUserID(ctx, taskID, userID)
}

func (t TaskUseCase) GetTaskRunList(ctx context.Context, taskID int) ([]entity.TaskRun, error) {
	return t.tr.GetTaskRunList(ctx, taskID)
}

func (t TaskUseCase) GetTaskRunLogList(ctx context.Context, taskID, lastID int) ([]entity.TaskRunLog, error) {
	return t.trl.GetTaskRunLogList(ctx, taskID, lastID)
}

func (t TaskUseCase) UploadRunLog(ctx context.Context, runLog entity.TaskRunLog) error {
	task, err := t.GetTaskDetail(ctx, runLog.TaskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app_code.New(app_code.ErrorNotFound, "任务不存在")
		}
		return err
	}

	switch task.Status {
	case entity.TaskStatusFinished:
		return app_code.New(app_code.ErrorBadRequest, "任务已经结束")
	case entity.TaskStatusNew:
		return app_code.New(app_code.ErrorBadRequest, "任务审核未完成")
	//case entity.TaskStatusPaused:
	//	return app_code.New(app_code.ErrorBadRequest, "任务暂停")
	case entity.StatusAuditReject:
		return app_code.New(app_code.ErrorBadRequest, "任务审核未通过")
	case entity.TaskStatusTorun:
		return app_code.New(app_code.ErrorRepeat, "任务未开始执行")
	}

	return t.trl.AddTaskRunLog(ctx, &runLog)
}

// PrepareTaskRun 准备签到接口,返回签到需要子任务id,通过子任务id生产二维码
func (t TaskUseCase) PrepareTaskRun(ctx context.Context, taskID int) (int, error) {
	// 1. 判断任务是否合法
	task, err := t.GetTaskDetail(ctx, taskID)
	if err != nil {
		return 0, err
	}

	switch task.Status {
	case entity.TaskStatusFinished:
		return 0, app_code.New(app_code.ErrorBadRequest, "任务已经结束")
	case entity.TaskStatusNew:
		return 0, app_code.New(app_code.ErrorBadRequest, "任务审核未完成")
	case entity.StatusAuditReject:
		return 0, app_code.New(app_code.ErrorBadRequest, "任务审核未通过")
	case entity.TaskStatusRunning:
		return 0, app_code.New(app_code.ErrorRepeat, "任务已经开始执行")
	}

	// 2. 判断task_run是否已经存在
	run, err := t.tr.GetPendingTaskRun(ctx, taskID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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

	{
		// 队长自动签到
		leader, ok, err := t.tu.GetTaskLeader(ctx, taskID)
		if !ok {
			return 0, errors.New("获取任务队长失败")
		}

		_, err = t.tru.AddTaskRunUser(ctx, taskID, taskRun.ID, leader.UserID)
		if err != nil {
			return 0, err
		}
	}

	return taskRun.ID, nil
}

// Sign 签到 参数为子任务id , 扫描时候使用
func (t TaskUseCase) Sign(ctx context.Context, taskID, taskRunID, userID int) error {
	_, err := t.tru.AddTaskRunUser(ctx, taskID, taskRunID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (t TaskUseCase) StartTaskRun(ctx context.Context, taskID int) error {
	task, err := t.GetTaskDetail(ctx, taskID)
	if err != nil {
		return err
	}

	switch task.Status {
	case entity.TaskStatusFinished:
		return app_code.New(app_code.ErrorBadRequest, "任务已经结束")
	case entity.TaskStatusNew:
		return app_code.New(app_code.ErrorBadRequest, "任务审核未完成")
	//case entity.TaskStatusPaused:
	//	return app_code.New(app_code.ErrorBadRequest, "任务暂停")
	case entity.StatusAuditReject:
		return app_code.New(app_code.ErrorBadRequest, "任务审核未通过")
	case entity.TaskStatusRunning:
		return app_code.New(app_code.ErrorRepeat, "任务已经开始执行")
	}

	run, err := t.tr.GetPendingTaskRun(ctx, taskID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return app_code.New(app_code.ErrorTaskRunNotFound, "任务暂停,请扫码签到开启任务")
	}

	glog.Infof("pending task run %v", run)
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

func (t TaskUseCase) AuditTask(ctx context.Context, taskID int, status string) (*entity.Task, error) {
	var err error
	var task *entity.Task
	switch status {
	case entity.StatusAuditReject:
		task, err = t.t.AuditFailTask(ctx, taskID)

	case entity.StatusAuditApproved:
		task, err = t.t.AuditSuccessTask(ctx, taskID)
	default:
		err = errors.New("arg status not found")
	}

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (t TaskUseCase) JoinTask(ctx context.Context, taskID, userID int) error {
	userTask, err := t.tu.GetUserTaskByUserID(ctx, taskID, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if userTask.ID != 0 {
		return app_code.New(app_code.ErrorRepeat, "已经发起参加该任务申请 ,当前审核状态:"+userTask.Status)
	}

	_, err = t.tu.AddUserTask(ctx, taskID, userID)
	return err
}

func (t TaskUseCase) AuditUserTask(ctx context.Context, taskID, userID int, status string) error {
	var err error
	switch status {
	case entity.UserTaskStatusAuditFail:
		_, err = t.tu.AuditUserTask(ctx, taskID, userID, status)

	case entity.UserTaskStatusAuditPass:
		_, err = t.tu.AuditUserTask(ctx, taskID, userID, status)

	default:
		err = app_code.New(app_code.ErrorBadRequest, "arg status not valid")
	}

	if err != nil {
		return err
	}
	return nil
}

func (t TaskUseCase) GetUserTaskSummary(ctx context.Context, userID int) (entity.UserTaskSummary, error) {
	return t.tru.GetUserTaskSummary(ctx, userID)
}

func (t TaskUseCase) GetUserJoinTaskList(ctx context.Context, userID int, status string, lastID int) ([]entity.UserTask, error) {
	return t.tu.GetUserJoinTaskList(ctx, userID, status, lastID)
}

func (t TaskUseCase) AssignRole(ctx context.Context, taskID, userID int, role string) error {
	return t.tu.AssignRole(ctx, taskID, userID, role)
}
