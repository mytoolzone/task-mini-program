package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mytoolzone/task-mini-program/internal/app_code"
	"net/http"
)

type response struct {
	Error string `json:"error" example:"message"`
}

func errorResponse(c *gin.Context, err error) {
	var appErr app_code.AppError
	ok := errors.As(err, &appErr)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, response{err.Error()})
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

	c.AbortWithStatusJSON(statusCode, response{appErr.Message})
	return
}
