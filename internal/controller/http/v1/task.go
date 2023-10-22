package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mytoolzone/task-mini-program/internal/app_error"
	"github.com/mytoolzone/task-mini-program/internal/controller/http/http_util"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/internal/usecase"
	"strconv"
)

type taskRoutes struct {
	task usecase.Task
}

func newTaskRoutes(handler *gin.RouterGroup, u usecase.Task) {
	ur := taskRoutes{u}

	h := handler.Group("/task")
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
		h.POST("/prepare", ur.prepare)
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

// 创建任务
func (r taskRoutes) create(ctx *gin.Context) {
	var request entity.Task
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errorResponse(ctx, app_error.WithError(app_error.ErrorBadRequest, err))
		return
	}

	request.CreateBy = ctx.GetInt("userID")
	err := r.task.CreateTask(ctx.Request.Context(), request)
	if err != nil {
		errorResponse(ctx, app_error.WithError(app_error.ErrorCreateTask, err))
		return
	}
	http_util.Success(ctx, nil)
}

// 获取任务详情
func (r taskRoutes) detail(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		errorResponse(ctx, app_error.New(app_error.ErrorBadRequest, "taskID is required"))
		return
	}

	detail, err := r.task.GetTaskDetail(ctx.Request.Context(), taskID)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	http_util.Success(ctx, detail)
}

// 获取任务列表
func (r taskRoutes) list(ctx *gin.Context) {
	lastIdStr, _ := ctx.GetQuery("lastId")
	lastId, _ := strconv.Atoi(lastIdStr)
	if lastId < 0 {
		lastId = 0
	}
	list, err := r.task.GetTaskList(ctx.Request.Context(), lastId)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	http_util.Success(ctx, list)
}

func (r taskRoutes) auditTask(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		errorResponse(ctx, app_error.New(app_error.ErrorBadRequest, "taskID is required"))
		return
	}
	auditStatus, _ := ctx.GetQuery("auditStatus")
	if auditStatus == "" {
		errorResponse(ctx, app_error.New(app_error.ErrorBadRequest, "auditStatus is required"))
		return
	}

	err := r.task.AuditTask(ctx.Request.Context(), taskID, auditStatus)
	if err != nil {
		errorResponse(ctx, err)
		return
	}
	http_util.Success(ctx, nil)
}

func (r taskRoutes) apply(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		errorResponse(ctx, app_error.New(app_error.ErrorBadRequest, "taskID is required"))
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

func (r taskRoutes) auditUserTask(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		errorResponse(ctx, app_error.New(app_error.ErrorBadRequest, "taskID is required"))
		return
	}
	userIDStr, _ := ctx.GetQuery("userID")
	userID, _ := strconv.Atoi(userIDStr)
	if userID <= 0 {
		errorResponse(ctx, app_error.New(app_error.ErrorBadRequest, "userID is required"))
		return
	}
	auditStatus, _ := ctx.GetQuery("auditStatus")
	if auditStatus == "" {
		errorResponse(ctx, app_error.New(app_error.ErrorBadRequest, "auditStatus is required"))
		return
	}

	err := r.task.AuditUserTask(ctx.Request.Context(), taskID, userID, auditStatus)
	if err != nil {
		return
	}

	http_util.Success(ctx, nil)
	return
}

func (r taskRoutes) prepare(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		errorResponse(ctx, app_error.New(app_error.ErrorBadRequest, "taskID is required"))
		return
	}

	run, err := r.task.PrepareTaskRun(ctx.Request.Context(), taskID)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	http_util.Success(ctx, run)
}

func (r taskRoutes) sign(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		errorResponse(ctx, app_error.New(app_error.ErrorBadRequest, "taskID is required"))
		return
	}

	userID := ctx.GetInt("userID")

	taskRunIDStr, _ := ctx.GetQuery("taskRunID")
	taskRunID, _ := strconv.Atoi(taskRunIDStr)
	if taskRunID <= 0 {
		errorResponse(ctx, app_error.New(app_error.ErrorBadRequest, "taskRunID is required"))
		return
	}

	err := r.task.Sign(ctx.Request.Context(), taskID, taskRunID, userID)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	http_util.Success(ctx, nil)
}

func (r taskRoutes) start(ctx *gin.Context) {
	taskRunIDStr, _ := ctx.GetQuery("taskRunID")
	taskRunID, _ := strconv.Atoi(taskRunIDStr)
	if taskRunID <= 0 {
		errorResponse(ctx, app_error.New(app_error.ErrorBadRequest, "taskRunID is required"))
		return
	}

	err := r.task.StartTaskRun(ctx.Request.Context(), taskRunID)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	http_util.Success(ctx, nil)
}

func (r taskRoutes) pause(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		errorResponse(ctx, app_error.New(app_error.ErrorBadRequest, "taskID is required"))
		return
	}

	err := r.task.PauseTaskRun(ctx.Request.Context(), taskID)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	http_util.Success(ctx, nil)
}

func (r taskRoutes) finish(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		errorResponse(ctx, app_error.New(app_error.ErrorBadRequest, "taskID is required"))
		return
	}

	err := r.task.FinishTaskRun(ctx.Request.Context(), taskID)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	http_util.Success(ctx, nil)

}

func (r taskRoutes) cancel(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		errorResponse(ctx, app_error.New(app_error.ErrorBadRequest, "taskID is required"))
		return
	}

	err := r.task.CancelTaskRun(ctx.Request.Context(), taskID)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	http_util.Success(ctx, nil)
}

func (r taskRoutes) runLogList(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		errorResponse(ctx, app_error.New(app_error.ErrorBadRequest, "taskID is required"))
		return
	}

	list, err := r.task.GetTaskRunLogList(ctx.Request.Context(), taskID)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	http_util.Success(ctx, list)
}

func (r taskRoutes) userList(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		errorResponse(ctx, app_error.New(app_error.ErrorBadRequest, "taskID is required"))
		return
	}
	userTasks, err := r.task.GetUserTasks(ctx.Request.Context(), taskID)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	http_util.Success(ctx, userTasks)
}

func (r taskRoutes) userTaskList(ctx *gin.Context) {
	userID := ctx.GetInt("userID")
	tasks, err := r.task.GetByUserID(ctx.Request.Context(), userID)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	http_util.Success(ctx, tasks)
}

func (r taskRoutes) uploadRunLog(ctx *gin.Context) {
	userID := ctx.GetInt("userID")
	var request entity.TaskRunLog
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errorResponse(ctx, app_error.WithError(app_error.ErrorBadRequest, err))
		return
	}

	request.UserID = userID
	err := r.task.UploadRunLog(ctx.Request.Context(), request)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	http_util.Success(ctx, nil)
}

func (r taskRoutes) userSummary(ctx *gin.Context) {
	userID := ctx.GetInt("userID")

	summary, err := r.task.GetUserTaskSummary(ctx.Request.Context(), userID)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	http_util.Success(ctx, summary)
}
