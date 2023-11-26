package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mytoolzone/task-mini-program/internal/app_code"
	"github.com/mytoolzone/task-mini-program/internal/controller/http/http_util"
	"github.com/mytoolzone/task-mini-program/internal/usecase"
	"strconv"
)

type noticeRoutes struct {
	n usecase.Notice
}

func newNoticeRoutes(handler *gin.RouterGroup, auth gin.HandlerFunc, roleH gin.HandlerFunc, u usecase.Notice) {
	ur := noticeRoutes{u}

	h := handler.Group("/notice", auth)
	{
		// 获取用户通知列表
		h.GET("/list", ur.list)
		// 获取用户通知详情
		h.GET("/detail", ur.detail)
	}
}

// @Summary 获取用户通知列表
// @Description 获取用户通知列表
// @Tags 通知
// @Accept json
// @Produce json
// @Param Authorization header string true "jwt_token"
// @Success 200 {object} http_util.Response{data=[]entity.Notice}
// @Failure 400 {object} http_util.Response
// @Failure 500 {object} http_util.Response
// @Router /notice/list [get]
func (r noticeRoutes) list(ctx *gin.Context) {
	userID := http_util.GetUserID(ctx)
	notices, err := r.n.GetNoticeListByUser(ctx, userID)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}

	http_util.Success(ctx, notices)
}

// @Summary 获取用户通知详情
// @Description 获取用户通知详情
// @Tags 通知
// @Accept json
// @Produce json
// @Param Authorization header string true "jwt_token"
// @Param noticeID query int true "通知ID"
// @Success 200 {object} http_util.Response{data=entity.Notice}
// @Failure 400 {object} http_util.Response
// @Failure 500 {object} http_util.Response
// @Router /notice/detail [get]
func (r noticeRoutes) detail(ctx *gin.Context) {
	noticeIDStr, _ := ctx.GetQuery("noticeID")
	noticeID, _ := strconv.Atoi(noticeIDStr)
	if noticeID <= 0 {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "not set params noticeID"))
		return
	}

	notices, err := r.n.GetNoticeByNoticeID(ctx, noticeID)
	if err != nil {
		http_util.Error(ctx, err)
		return
	}

	userID := http_util.GetUserID(ctx)
	if notices.UserID != userID {
		http_util.Error(ctx, app_code.New(app_code.ErrorBadRequest, "not your notice"))
		return
	}

	err = r.n.SetNoticeRead(ctx, noticeID)
	if err != nil {
		return
	}
	http_util.Success(ctx, notices)
}
