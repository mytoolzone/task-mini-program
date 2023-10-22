package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mytoolzone/task-mini-program/internal/usecase"
)

type taskRoutes struct {
	u usecase.Task
}

func (r taskRoutes) create(context *gin.Context) {

}

func (r taskRoutes) detail(context *gin.Context) {

}

func (r taskRoutes) list(context *gin.Context) {

}

func (r taskRoutes) audit(context *gin.Context) {

}

func (r taskRoutes) prepare(context *gin.Context) {

}

func (r taskRoutes) sign(context *gin.Context) {

}

func (r taskRoutes) start(context *gin.Context) {

}

func (r taskRoutes) pause(context *gin.Context) {

}

func (r taskRoutes) finish(context *gin.Context) {

}

func (r taskRoutes) cancel(context *gin.Context) {

}

func (r taskRoutes) runList(context *gin.Context) {

}

func (r taskRoutes) runLogList(context *gin.Context) {

}

func (r taskRoutes) userList(context *gin.Context) {

}

func (r taskRoutes) userTaskList(context *gin.Context) {

}

func newUserRoutes(handler *gin.RouterGroup, u usecase.Task) {
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
		h.POST("/audit", ur.audit)
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
		// 获取任务运行列表
		h.GET("/runList", ur.runList)
		// 获取任务运行日志列表
		h.GET("/runLogList", ur.runLogList)
		// 获取任务用户列表
		h.GET("/userList", ur.userList)
		// 获取某人任务列表
		h.GET("/userTaskList", ur.userTaskList)

	}
}
