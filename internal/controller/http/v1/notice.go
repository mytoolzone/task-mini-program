package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mytoolzone/task-mini-program/internal/app_error"
	"github.com/mytoolzone/task-mini-program/internal/controller/http/http_util"
	"github.com/mytoolzone/task-mini-program/internal/usecase"
	"strconv"
)

type noticeRoutes struct {
	n usecase.Notice
}

func newNoticeRoutes(handler *gin.RouterGroup, u usecase.Notice) {
	ur := noticeRoutes{u}

	h := handler.Group("/notice")
	{
		// 获取用户通知列表
		h.GET("/list", ur.list)
		// 获取用户通知详情
		h.GET("/detail", ur.detail)
	}
}

func (r noticeRoutes) list(ctx *gin.Context) {
	userID := ctx.GetInt("userID")
	notices, err := r.n.GetNoticeListByUser(ctx, userID)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	http_util.Success(ctx, notices)
}

func (r noticeRoutes) detail(ctx *gin.Context) {
	noticeIDStr, _ := ctx.GetQuery("noticeID")
	noticeID, _ := strconv.Atoi(noticeIDStr)
	if noticeID <= 0 {
		errorResponse(ctx, app_error.New(app_error.ErrorBadRequest, "not set params noticeID"))
		return
	}

	notices, err := r.n.GetNoticeByNoticeID(ctx, noticeID)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	userID := ctx.GetInt("userID")

	if notices.UserID != userID {
		errorResponse(ctx, app_error.New(app_error.ErrorBadRequest, "not your notice"))
		return
	}

	http_util.Success(ctx, notices)
}
