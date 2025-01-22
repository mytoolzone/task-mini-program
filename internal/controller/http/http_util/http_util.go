package http_util

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gw123/glog"
	"github.com/mytoolzone/task-mini-program/internal/app_code"
	"github.com/mytoolzone/task-mini-program/internal/entity"
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
	// fmt.Printf("appErr: %#v", appErr.Code)
	if !ok {
		glog.Errorf("http request error %+v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, Response{
			Error: err.Error(),
			Code:  app_code.ErrorServerError,
		})
		return
	}

	statusCode := http.StatusInternalServerError
	// 匹配code转成httpcode编码
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
	case app_code.ErrorTokenNotSet:
		statusCode = http.StatusBadRequest
	case app_code.ErrorRepeat:
		statusCode = http.StatusOK
	case app_code.ErrorAuthFailed, app_code.ErrorTokenTimeout:
		statusCode = http.StatusOK
	}

	glog.Errorf("http request code [%+v] - err %+v", appErr.Code, appErr.Message)
	c.AbortWithStatusJSON(statusCode, Response{
		Code:  appErr.Code,
		Error: appErr.Message,
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

// 根据用户角色判断userID是否合法
func CheckUserID(ctx *gin.Context) (int, error) {
	var userID int
	userRole := GetUserRole(ctx) //获取用户角色
	i, err := strconv.Atoi(ctx.Query("userID"))
	if err == nil {
		userID = i
	}

	if userRole == entity.UserRoleAdmin { //管理员可以搜索指定人员数据

	} else {
		selfuserID := GetUserID(ctx) //非管理员查自己的数据
		if selfuserID != userID {
			return 0, errors.New("非法操作")
		}
		return userID, nil
	}
	return userID, nil
}
