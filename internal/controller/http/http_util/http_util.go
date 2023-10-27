package http_util

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mytoolzone/task-mini-program/internal/app_code"
	"net/http"
)

type Response struct {
	Code  app_code.CodeType `json:"code" example:"success"`
	Error string            `json:"error" example:"error message"`
	Data  interface{}       `json:""`
}

const (
	ctxUserIDKey = "ctxUserID"
	ctxUserName  = "ctxUserName"
)

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, Response{
		Code: app_code.CodeType("success"),
		Data: data,
	})
}

func Error(c *gin.Context, err error) {
	var appErr app_code.AppError
	ok := errors.As(err, &appErr)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, Response{
			Error: err.Error(),
			Code:  app_code.ErrorServerError,
		})
		return
	}

	statusCode := http.StatusInternalServerError
	switch appErr.Code {
	case app_code.ErrorBadRequest:
		statusCode = http.StatusBadRequest
	case app_code.ErrorNotFound:
	case app_code.ErrorUserNotFound:
	case app_code.ErrorTaskNotFound:
		statusCode = http.StatusNotFound
	case app_code.ErrorTaskExist:
	case app_code.ErrorUserExist:
		statusCode = http.StatusConflict
	}

	c.AbortWithStatusJSON(statusCode, Response{
		Error: appErr.Message,
		Code:  appErr.Code,
	})
	return
}

// GetUserID 从 ctx 中获取用户 id
func GetUserID(ctx *gin.Context) int {
	return ctx.GetInt(ctxUserIDKey)
}

// SetUserID 将用户 id 设置到 ctx 中
func SetUserID(ctx *gin.Context, userID int) {
	ctx.Set(ctxUserIDKey, userID)
}

// SetUserName 将用户名设置到 ctx 中
func SetUserName(ctx *gin.Context, userName string) {
	ctx.Set(ctxUserName, userName)
}

// GetUserName 从 ctx 中获取用户名
func GetUserName(ctx *gin.Context) string {
	return ctx.GetString(ctxUserName)
}
