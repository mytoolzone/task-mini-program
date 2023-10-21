package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mytoolzone/task-mini-program/internal/app_error"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/internal/usecase"
	"net/http"
)

type userRoutes struct {
	u usecase.User
}

func newUserRoutes(handler *gin.RouterGroup, u usecase.User) {
	ur := userRoutes{u}

	h := handler.Group("/user")
	{
		h.POST("/login", ur.login)
		h.POST("/register", ur.register)
		h.POST("/updateSetting", ur.updateSetting)
		h.GET("/getSetting", ur.getSetting)
	}
}

func (ur userRoutes) login(context *gin.Context) {

}

func (ur userRoutes) register(context *gin.Context) {

}

type doUpdateUserSettingRequest struct {
	UserSetting entity.UserSetting `json:"userSetting" binding:"required"`
	UserID      int                `json:"userID" binding:"required"`
}

func (ur userRoutes) updateSetting(c *gin.Context) {
	var request doUpdateUserSettingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, app_error.WithError(app_error.ErrorBadRequest, err))
		return
	}

	err := ur.u.UpdateSetting(c.Request.Context(), request.UserID, request.UserSetting)
	if err != nil {
		errorResponse(c, app_error.WithError(app_error.ErrorUpdateUserSetting, err))
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{})
}

func (ur userRoutes) getSetting(c *gin.Context) {
	userID := c.GetInt("userID")
	setting, err := ur.u.GetSettingByUserID(c.Request.Context(), userID)
	if err != nil {
		errorResponse(c, app_error.WithError(app_error.ErrorGetUserSetting, err))
		return
	}
	c.JSON(http.StatusOK, setting)
}
