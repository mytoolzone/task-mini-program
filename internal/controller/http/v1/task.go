package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mytoolzone/task-mini-program/internal/app_code"
	"github.com/mytoolzone/task-mini-program/internal/controller/http/http_util"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/internal/usecase"
	"strconv"
)

type taskRoutes struct {
	task usecase.Task
}

func newTaskRoutes(handler *gin.RouterGroup, auth gin.HandlerFunc, u usecase.Task) {
	ur := taskRoutes{u}

	h := handler.Group("/task", auth)
	{
		// 创建任务
		h.POST("/create", ur.create)
		// 获取任务详情
		h.GET("/detail", ur.detail)
		// 获取任务列表
		h.GET("/list", ur.list)
		// 审核任务
		h.POST("/audiTask", ur.auditTask)
		// 报名任务
		h.POST("/apply", ur.apply)
		// 审核报名
		h.POST("/auditUserTask", ur.auditUserTask)
		// 获取签到二维码
		h.GET("/prepare", ur.prepare)
		// 签到
		h.POST("/sign", ur.sign)
		// 开始任务
		h.POST("/start", ur.start)
		// 暂停任务
		h.POST("/pause", ur.pause)
		// 完成任务
		h.POST("/finish", ur.finish)
		// 取消任务
		h.POST("/cancel", ur.cancel)
		// 上报任务运行日志列表
		h.POST("/uploadRunLog", ur.uploadRunLog)
		// 获取任务运行日志列表
		h.GET("/runLogs", ur.runLogList)
		// 获取任务用户列表
		h.GET("/users", ur.userList)
		// 获取某人任务列表
		h.GET("/userTasks", ur.userTaskList)
		// 获取某个用户的统计数据
		h.GET("/userSummary", ur.userSummary)
	}
}

// @Summary     Create task
// @Description 创建任务
// @ID          create-task
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Success     200 {object} http_util.Response
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/create [post]
func (r taskRoutes) create(ctx *gin.Context) {
	var request entity.Task
	if err := ctx.ShouldBindJSON(&request); err != nil {
		http_util.Error(ctx, app_code.WithError(app_code.ErrorBadRequest, err))
		return
	}

	request.CreateBy = ctx.GetInt("userID")
	err := r.task.CreateTask(ctx.Request.Context(), request)
	if err != nil {
		http_util.Error(ctx, app_code.WithError(app_code.ErrorCreateTask, err))
		return
	}
	http_util.Success(ctx, nil)
}

// @Summary     Show task detail
// @Description 获取任务详情
// @ID          detail
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param       taskID query int true "taskID"
// @Success     200 {object} http_util.Response{data=entity.Task}
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/detail [get]
func (r taskRoutes) detail(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "taskID is required"))
		return
	}

	detail, err := r.task.GetTaskDetail(ctx.Request.Context(), taskID)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}

	http_util.Success(ctx, detail)
}

// 获取任务列表
// @Summary     List tasks
// @Description 获取任务列表
// @ID          list
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param       lastId query int false "lastId"
// @Success     200 {object} http_util.Response{data=[]entity.Task}
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/list [get]
func (r taskRoutes) list(ctx *gin.Context) {
	lastIdStr, _ := ctx.GetQuery("lastId")
	lastId, _ := strconv.Atoi(lastIdStr)
	if lastId < 0 {
		lastId = 0
	}
	list, err := r.task.GetTaskList(ctx.Request.Context(), lastId)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}
	http_util.Success(ctx, list)
}

// @Summary     Audit task
// @Description 管理员审核任务
// @ID          audit-task
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param       taskID query int true "taskID"
// @Param       auditStatus query string true "auditStatus"
// @Success     200 {object} http_util.Response
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/auditTask [post]
func (r taskRoutes) auditTask(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "taskID is required"))
		return
	}
	auditStatus, _ := ctx.GetQuery("auditStatus")
	if auditStatus == "" {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "auditStatus is required"))
		return
	}

	err := r.task.AuditTask(ctx.Request.Context(), taskID, auditStatus)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}
	http_util.Success(ctx, nil)
}

// @Summary     Apply task
// @Description 用户申请参加任务
// @ID          apply-task
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param       taskID query int true "taskID"
// @Success     200 {object} http_util.Response
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/apply [post]
func (r taskRoutes) apply(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "taskID is required"))
		return
	}
	userID := ctx.GetInt("userID")

	err := r.task.JoinTask(ctx.Request.Context(), taskID, userID)
	if err != nil {
		return
	}

	http_util.Success(ctx, nil)
	return
}

// @Summary     Audit user task
// @Description 管理员审核用户参加任务
// @ID          audit-user-task
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param       taskID query int true "taskID"
// @Param       userID query int true "userID"
// @Param       auditStatus query string true "auditStatus"
// @Success     200 {object} http_util.Response
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/auditUserTask [post]
func (r taskRoutes) auditUserTask(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "taskID is required"))
		return
	}
	userIDStr, _ := ctx.GetQuery("userID")
	userID, _ := strconv.Atoi(userIDStr)
	if userID <= 0 {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "userID is required"))
		return
	}
	auditStatus, _ := ctx.GetQuery("auditStatus")
	if auditStatus == "" {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "auditStatus is required"))
		return
	}

	err := r.task.AuditUserTask(ctx.Request.Context(), taskID, userID, auditStatus)
	if err != nil {
		return
	}

	http_util.Success(ctx, nil)
	return
}

// @Summary     Prepare task
// @Description 获取签到二维码
// @ID          prepare-task
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param       taskID query int true "taskID"
// @Success     200 {object} http_util.Response{data=entity.TaskRun}
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/prepare [get]
func (r taskRoutes) prepare(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "taskID is required"))
		return
	}

	run, err := r.task.PrepareTaskRun(ctx.Request.Context(), taskID)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}

	http_util.Success(ctx, run)
}

// @Summary     Sign task
// @Description 签到
// @ID          sign-task
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param       taskID query int true "taskID"
// @Param       taskRunID query int true "taskRunID"
// @Success     200 {object} http_util.Response
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/sign [post]
func (r taskRoutes) sign(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "taskID is required"))
		return
	}

	taskRunIDStr, _ := ctx.GetQuery("taskRunID")
	taskRunID, _ := strconv.Atoi(taskRunIDStr)
	if taskRunID <= 0 {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "taskRunID is required"))
		return
	}

	userID := ctx.GetInt("userID")
	err := r.task.Sign(ctx.Request.Context(), taskID, taskRunID, userID)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}

	http_util.Success(ctx, nil)
}

// @Summary     Start task
// @Description 开始任务
// @ID          start-task
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param       taskRunID query int true "taskID"
// @Success     200 {object} http_util.Response
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/start [post]
func (r taskRoutes) start(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "taskRunID is required"))
		return
	}

	err := r.task.StartTaskRun(ctx.Request.Context(), taskID)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}

	http_util.Success(ctx, nil)
}

// @Summary     Pause task
// @Description 暂停任务
// @ID          pause-task
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param       taskID query int true "taskID"
// @Success     200 {object} http_util.Response
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/pause [post]
func (r taskRoutes) pause(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "taskID is required"))
		return
	}

	err := r.task.PauseTaskRun(ctx.Request.Context(), taskID)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}

	http_util.Success(ctx, nil)
}

// @Summary     Finish task
// @Description 完成任务
// @ID          finish-task
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param       taskID query int true "taskID"
// @Success     200 {object} http_util.Response
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/finish [post]
func (r taskRoutes) finish(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "taskID is required"))
		return
	}

	err := r.task.FinishTaskRun(ctx.Request.Context(), taskID)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}

	http_util.Success(ctx, nil)
}

// @Summary     Cancel task
// @Description 取消任务
// @ID          cancel-task
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param       taskID query int true "taskID"
// @Success     200 {object} http_util.Response
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/cancel [post]
func (r taskRoutes) cancel(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "taskID is required"))
		return
	}

	err := r.task.CancelTaskRun(ctx.Request.Context(), taskID)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}

	http_util.Success(ctx, nil)
}

// @Summary     Run log list
// @Description 获取任务运行日志列表
// @ID          run-log-list
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param       taskID query int true "taskID"
// @Param       lastId query int false "lastId"
// @Success     200 {object} http_util.Response{data=[]entity.TaskRunLog}
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/runLogs [get]
func (r taskRoutes) runLogList(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "taskID is required"))
		return
	}

	list, err := r.task.GetTaskRunLogList(ctx.Request.Context(), taskID)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}

	http_util.Success(ctx, list)
}

// @Summary     User list
// @Description 获取任务用户列表
// @ID          user-list
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param       taskID query int true "taskID"
// @Success     200 {object} http_util.Response{data=[]entity.UserTask}
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/users [get]
func (r taskRoutes) userList(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "taskID is required"))
		return
	}
	userTasks, err := r.task.GetUserTasks(ctx.Request.Context(), taskID)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}

	http_util.Success(ctx, userTasks)
}

// @Summary     User task list
// @Description 获取某人任务列表
// @ID          user-task-list
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param       userID query int true "userID"
// @Success     200 {object} http_util.Response{data=[]entity.Task}
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/userTasks [get]
func (r taskRoutes) userTaskList(ctx *gin.Context) {
	userID := ctx.GetInt("userID")
	tasks, err := r.task.GetByUserID(ctx.Request.Context(), userID)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}

	http_util.Success(ctx, tasks)
}

// @Summary     Upload run log
// @Description 上报任务运行日志列表
// @ID          upload-run-log
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param       jsonBody body entity.TaskRunLog true "上报任务运行日志"
// @Success     200 {object} http_util.Response
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/uploadRunLog [post]
func (r taskRoutes) uploadRunLog(ctx *gin.Context) {
	userID := ctx.GetInt("userID")
	var request entity.TaskRunLog
	if err := ctx.ShouldBindJSON(&request); err != nil {
		http_util.Error(ctx, app_code.WithError(app_code.ErrorBadRequest, err))
		return
	}

	request.UserID = userID
	err := r.task.UploadRunLog(ctx.Request.Context(), request)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}

	http_util.Success(ctx, nil)
}

// @Summary     User summary
// @Description 获取某个用户的统计数据
// @ID          user-summary
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param       userID query int true "userID"
// @Success     200 {object} http_util.Response{data=entity.UserTaskSummary}
func (r taskRoutes) userSummary(ctx *gin.Context) {
	userID := ctx.GetInt("userID")

	summary, err := r.task.GetUserTaskSummary(ctx.Request.Context(), userID)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}

	http_util.Success(ctx, summary)
}
