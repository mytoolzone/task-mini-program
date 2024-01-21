package v1

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gw123/glog"
	"github.com/mytoolzone/task-mini-program/internal/app_code"
	"github.com/mytoolzone/task-mini-program/internal/controller/http/http_util"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/internal/usecase"
)

type taskRoutes struct {
	task   usecase.Task
	notice usecase.Notice
}

func newTaskRoutes(handler *gin.RouterGroup, auth gin.HandlerFunc, role gin.HandlerFunc, u usecase.Task, n usecase.Notice) {
	ur := taskRoutes{u, n}

	h := handler.Group("/task", auth, role)
	{
		// 创建任务
		h.POST("/create", ur.create)
		// 获取任务详情
		h.GET("/detail", ur.detail)
		// 获取任务列表
		h.GET("/list", ur.list)
		// 审核任务
		h.POST("/auditTask", ur.auditTask)
		// 分配参与任务人角色
		h.POST("/assignRole", ur.assignRole)
		// 报名任务
		h.POST("/apply", ur.apply)
		// 获取任务报名用户列表
		h.GET("/applyUsers", ur.applyUserList)
		// 审核报名
		h.POST("/auditApplyTask", ur.auditApplyTask)
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
		// 获取任务用户列表 全任务周期
		h.GET("/users", ur.userList)
		// 获取某人创建的任务列表
		h.GET("/userTasks", ur.userTaskList)
		// 获取某个人参加的任务
		h.GET("/userJoinTask", ur.userJoinTask)
		// 获取某个用户的统计数据
		h.GET("/userSummary", ur.userSummary)
		// 获取任务参加人列表获取某个任务，已经审核通过的人列表用在分配角色环节
		h.GET("/approvedUsers", ur.approvedUsers)
	}
}

// Task mapped from table <tasks>
type Task struct {
	Name     string `gorm:"column:name;not null" json:"name"`
	Describe string `gorm:"column:describe" json:"describe"`
	Require  string `gorm:"column:require" json:"require"`
	Location string `gorm:"column:location" json:"location"`
}

// @Summary     Create task
// @Description 创建任务
// @ID          create-task
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param Authorization header string true "jwt_token"
// @Param       jsonBody body Task true "创建任务"
// @Success     200 {object} http_util.Response{data=entity.Task}
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/create [post]
func (r taskRoutes) create(ctx *gin.Context) {
	var task entity.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		http_util.Error(ctx, app_code.WithError(app_code.ErrorBadRequest, err))
		return
	}

	task.CreateBy = http_util.GetUserID(ctx)
	task.Status = entity.TaskStatusNew
	err := r.task.CreateTask(ctx.Request.Context(), &task)
	if err != nil {
		http_util.Error(ctx, app_code.WithError(app_code.ErrorCreateTask, err))
		return
	}
	http_util.Success(ctx, task)
}

// @Summary     Show task detail
// @Description 获取任务详情
// @ID          detail
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param Authorization header string true "jwt_token"
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
// @Param Authorization header string true "jwt_token"
// @Param       lastID query int false "lastID"
// @Param       status query string false "status"  Enums(torun, audit_fail,join,new,running,paused,finished,canceled,deleted)
// @Success     200 {object} http_util.Response{data=[]entity.Task}
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/list [get]
func (r taskRoutes) list(ctx *gin.Context) {
	lastIDStr, _ := ctx.GetQuery("lastID")
	lastID, _ := strconv.Atoi(lastIDStr)
	if lastID < 0 {
		lastID = 0
	}

	statusStr, _ := ctx.GetQuery("status")
	keyword, _ := ctx.GetQuery("keyword")

	list, err := r.task.GetTaskList(ctx.Request.Context(), lastID, keyword, statusStr)
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
// @Param Authorization header string true "jwt_token"
// @Param       taskID query int true "taskID"
// @Param       auditStatus query string true "auditStatus" Enums(rejected, approved)
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

	task, err := r.task.AuditTask(ctx.Request.Context(), taskID, auditStatus)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}

	// 发送通知
	go func() {
		var content string
		if auditStatus == entity.StatusAuditReject {
			content = "您的任务'" + task.Name + "'已被管理员拒绝，请重新提交任务"
		} else {
			content = "您的任务'" + task.Name + "'已被管理员通过，请及时开启任务"
		}

		err := r.notice.AddNotice(ctx.Request.Context(), entity.Notice{
			Title:   "任务审核",
			Content: content,
			Status:  entity.NotifyStatusUnread,
			UserID:  task.CreateBy,
		})
		if err != nil {
			glog.WithErr(err).Error("add notice field")
		}
	}()

	http_util.Success(ctx, nil)
}

// @Summary     Assign role
// @Description 管理员分配任务角色
// @ID          assign-role
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param Authorization header string true "jwt_token"
// @Param       taskID query int true "taskID"
// @Param       userID query int true "userID"
// @Param       role query string true "role" Enums(leader, member, recorder)
// @Success     200 {object} http_util.Response
// @Failure     400 {object} http_util.Response
// @Router      /task/assignRole [post]
func (r taskRoutes) assignRole(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "taskID is required"))
		return
	}
	role, _ := ctx.GetQuery("role")
	if role == "" {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "role is required"))
		return
	}
	userIDStr, _ := ctx.GetQuery("userID")
	userID, _ := strconv.Atoi(userIDStr)
	if userID <= 0 {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "userID is required"))
		return
	}

	err := r.task.AssignRole(ctx.Request.Context(), taskID, userID, role)
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
// @Param Authorization header string true "jwt_token"
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
	userID := http_util.GetUserID(ctx)

	err := r.task.JoinTask(ctx.Request.Context(), taskID, userID)
	if err != nil {
		http_util.Error(ctx, err)
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
// @Param Authorization header string true "jwt_token"
// @Param       taskID query int true "taskID"
// @Param       userID query int true "userID"
// @Param       auditStatus query string true "auditStatus" Enums(rejected, approved)
// @Success     200 {object} http_util.Response
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/auditApplyTask [post]
func (r taskRoutes) auditApplyTask(ctx *gin.Context) {
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
		http_util.Error(ctx, err)
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
// @Param Authorization header string true "jwt_token"
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
// @Param Authorization header string true "jwt_token"
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

	userID := http_util.GetUserID(ctx)
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
// @Param Authorization header string true "jwt_token"
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
// @Param Authorization header string true "jwt_token"
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
// @Param Authorization header string true "jwt_token"
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
// @Param Authorization header string true "jwt_token"
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
// @Param Authorization header string true "jwt_token"
// @Param       taskID query int true "taskID"
// @Param       lastID query int false "lastID"
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
	lastIDStr, _ := ctx.GetQuery("lastID")
	lastID, _ := strconv.Atoi(lastIDStr)

	list, err := r.task.GetTaskRunLogList(ctx.Request.Context(), taskID, lastID)
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
// @Param Authorization header string true "jwt_token"
// @Param       taskID query int true "taskID"
// @Param       status query string false "status"
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

	status := ctx.Query("status")

	userTasks, err := r.task.GetUserTasks(ctx.Request.Context(), taskID, status)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}

	http_util.Success(ctx, userTasks)
}

// @Summary     Apply Task User list
// @Description 获取任务报名用户列表
// @ID          apply-user-list
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param Authorization header string true "jwt_token"
// @Param       taskID query int true "taskID"
// @Success     200 {object} http_util.Response{data=[]entity.UserTask}
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/users [get]
func (r taskRoutes) applyUserList(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "taskID is required"))
		return
	}

	userTasks, err := r.task.GetUserTasks(ctx.Request.Context(), taskID, entity.UserTaskStatusApply)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}

	http_util.Success(ctx, userTasks)
}

// @Summary     User task list
// @Description 获取某人创建的任务列表
// @ID          user-task-list
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param Authorization header string true "jwt_token"
// @Param       userID query int true "userID"
// @Param       lastID query int false "lastID"
// @Param       status query string false "status" enum(new, running, pause, finish, cancel)
// @Success     200 {object} http_util.Response{data=[]entity.Task}
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/userTasks [get]
func (r taskRoutes) userTaskList(ctx *gin.Context) {
	userID := http_util.GetUserID(ctx)
	lastIdStr, _ := ctx.GetQuery("lastID")
	lastId, _ := strconv.Atoi(lastIdStr)
	status := ctx.Query("status")
	tasks, err := r.task.GetByUserID(ctx.Request.Context(), userID, status, lastId)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}

	http_util.Success(ctx, tasks)
}

// @Summary     User task list
// @Description 获取某人参加的任务列表
// @ID          user-join-tasks
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param Authorization header string true "jwt_token"
// @Param       userID query int true "userID"
// @Param       lastID query int false "lastID"
// @Param       status query string false "status" enum(apply,rejected,approved,running,finish,cancel)
// @Success     200 {object} http_util.Response{data=[]entity.Task}
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/userJoinTask [get]
func (r taskRoutes) userJoinTask(ctx *gin.Context) {
	userID := http_util.GetUserID(ctx)
	lastIdStr, _ := ctx.GetQuery("lastID")
	lastId, _ := strconv.Atoi(lastIdStr)
	status := ctx.Query("status")
	tasks, err := r.task.GetUserJoinTaskList(ctx.Request.Context(), userID, status, lastId)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}
	http_util.Success(ctx, tasks)
}

type UploadRunLogRequest struct {
	TaskID  int      `json:"taskID"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
	Videos  []string `json:"videos"`
}

// @Summary     Upload run log
// @Description 上报任务运行日志列表
// @ID          upload-run-log
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param Authorization header string true "jwt_token"
// @Param       jsonBody body UploadRunLogRequest true "上报任务运行日志"
// @Success     200 {object} http_util.Response
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/uploadRunLog [post]
func (r taskRoutes) uploadRunLog(ctx *gin.Context) {
	var request UploadRunLogRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		http_util.Error(ctx, app_code.WithError(app_code.ErrorBadRequest, err))
		return
	}

	var runLog = entity.TaskRunLog{
		TaskID:  request.TaskID,
		Content: request.Content,
		UserID:  http_util.GetUserID(ctx),
		Images:  strings.Join(request.Images, ","),
		Videos:  strings.Join(request.Videos, ","),
	}

	err := r.task.UploadRunLog(ctx.Request.Context(), runLog)
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
// @Param Authorization header string true "jwt_token"
// @Param       userID query int true "userID"
// @Success     200 {object} http_util.Response{data=entity.UserTaskSummary}
func (r taskRoutes) userSummary(ctx *gin.Context) {
	userID := http_util.GetUserID(ctx)

	summary, err := r.task.GetUserTaskSummary(ctx.Request.Context(), userID)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}

	http_util.Success(ctx, summary)
}

// @Summary     approvedUsers list
// @Description 获取任务参加人列表获取某个任务，已经审核通过的人列表,该接口用在分配角色环节
// @ID          approvedUsers-list
// @Tags  	    task
// @Accept      json
// @Produce     json
// @Param Authorization header string true "jwt_token"
// @Param       taskID query int true "taskID"
// @Success     200 {object} http_util.Response{data=[]entity.UserTask}
// @Failure     400 {object} http_util.Response
// @Failure     500 {object} http_util.Response
// @Router      /task/approvedUsers [get]
func (r taskRoutes) approvedUsers(ctx *gin.Context) {
	taskIDStr, _ := ctx.GetQuery("taskID")
	taskID, _ := strconv.Atoi(taskIDStr)
	if taskID <= 0 {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "taskID is required"))
		return
	}

	userTasks, err := r.task.GetUserTasks(ctx.Request.Context(), taskID, entity.UserTaskStatusAuditPass)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}
	http_util.Success(ctx, userTasks)
}
