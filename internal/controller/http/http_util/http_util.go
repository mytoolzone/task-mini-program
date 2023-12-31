package http_util

import (
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gw123/glog"
	"github.com/mytoolzone/task-mini-program/internal/app_code"
)

type Response struct {
	Code  app_code.CodeType `json:"code" example:"success"`
	Error string            `json:"error" example:"error message"`
	Data  interface{}       `json:""`
}

const (
	ctxUserIDKey = "ctxUserID"
	ctxUserName  = "ctxUserName"
	ctxUserRole  = "ctxUserRole"
)

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, Response{
		Code: app_code.CodeType("success"),
		Data: data,
	})
}

func Error(c *gin.Context, err error) {
	appErr, ok := err.(*app_code.AppError)
	if !ok {
		glog.Errorf("http request error %+v", err)
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
	case app_code.ErrorRepeat:
		statusCode = http.StatusOK
	}

	glog.Errorf("http request code [%+v] - err %+v", appErr.Code, appErr.Message)
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

// SetUserRole 将用户角色设置到 ctx 中
func SetUserRole(ctx *gin.Context, role string) {
	ctx.Set(ctxUserRole, role)
}

// GetUserRole 从 ctx 中获取用户角色
func GetUserRole(ctx *gin.Context) string {
	return ctx.GetString(ctxUserRole)
}

// SetUserName 将用户名设置到 ctx 中
func SetUserName(ctx *gin.Context, userName string) {
	ctx.Set(ctxUserName, userName)
}

// GetUserName 从 ctx 中获取用户名
func GetUserName(ctx *gin.Context) string {
	return ctx.GetString(ctxUserName)
}

// IsImage 检测文件是否是图片
func IsImage(file *multipart.FileHeader) bool {
	// 获取文件后缀
	tmpArr := strings.Split(file.Filename, ".")
	ext := tmpArr[len(tmpArr)-1]
	switch ext {
	case "jpg":
		return true
	case "jpeg":
		return true
	case "png":
		return true
	case "gif":
		return true
	case "bmp":
		return true
	case "webp":
		return true
	}

	switch file.Header.Get("Content-Type") {
	case "image/jpeg":
		return true
	case "image/png":
		return true
	case "image/gif":
		return true
	case "image/bmp":
		return true
	case "image/webp":
		return true
	}

	return false
}
