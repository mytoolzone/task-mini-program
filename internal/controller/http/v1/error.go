package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mytoolzone/task-mini-program/internal/app_error"
	"net/http"
)

type response struct {
	Error string `json:"error" example:"message"`
}

func errorResponse(c *gin.Context, err error) {
	var appErr app_error.AppError
	ok := errors.As(err, &appErr)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{err.Error()})
		return
	}

	statusCode := http.StatusInternalServerError
	switch appErr.Code {
	case app_error.ErrorBadRequest:
		statusCode = http.StatusBadRequest
	case app_error.ErrorNotFound:
	case app_error.ErrorUserNotFound:
	case app_error.ErrorTaskNotFound:
		statusCode = http.StatusNotFound
	case app_error.ErrorTaskExist:
	case app_error.ErrorUserExist:
		statusCode = http.StatusConflict
	}

	c.AbortWithStatusJSON(statusCode, response{appErr.Message})
	return
}
