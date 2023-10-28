package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mytoolzone/task-mini-program/internal/app_code"
	"github.com/mytoolzone/task-mini-program/internal/controller/http/http_util"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/internal/usecase"
	"github.com/mytoolzone/task-mini-program/pkg/auth"
	"net/http"
)

type userRoutes struct {
	u usecase.User
	a auth.Auth
}

func newUserRoutes(handler *gin.RouterGroup, authH gin.HandlerFunc, a auth.Auth, u usecase.User) {
	ur := userRoutes{u, a}

	h := handler.Group("/user")
	{
		h.POST("/login", ur.login)
		h.POST("/miniProgramLogin", ur.miniProgramLogin)
		h.POST("/register", ur.register)
		h.POST("/updateSetting", ur.updateSetting, authH)
		h.GET("/getSetting", ur.getSetting, authH)
	}
}

type doMiniProgramLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

type doMiniProgramLoginResponse struct {
	Token    string `json:"token"`
	UserID   int    `json:"userId"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
}

// @Summary 小程序登录
// @Description 小程序登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param jsonBody body doMiniProgramLoginRequest true  "登录参数"
// @Success 200 {object} doLoginResponse
// @Failure 400 {object} http_util.Response
// @Failure 500 {object} http_util.Response
// @Router /user/miniProgramLogin [post]
func (ur userRoutes) miniProgramLogin(context *gin.Context) {
	//登录获取返回 jwt_token
	//1.获取参数
	//2.参数校验
	//3.根据用户名和密码查询用户信息
	//4.生成token
	//5.返回token
	var userReq doMiniProgramLoginRequest
	if err := context.ShouldBindJSON(&userReq); err != nil {
		http_util.Error(context, app_code.WithError(app_code.ErrorBadRequest, err))
		return
	}

	user, err := ur.u.MiniProgramLogin(context.Request.Context(), userReq.Code)
	if err != nil {
		http_util.Error(context, err)
		return
	}

	token, err := ur.a.GenerateToken(user.Username, user.ID)
	if err != nil {
		http_util.Error(context, err)
		return
	}

	context.JSON(http.StatusOK, doMiniProgramLoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Phone:    user.Phone,
	})
}

type doLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type doLoginResponse struct {
	Token    string `json:"token"`
	UserID   int    `json:"userId"`
	Username string `json:"username"`
}

// @Summary 登录
// @Description 登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param jsonBody body doLoginRequest true "登录参数"
// @Success 200 {object} doLoginResponse
// @Failure 400 {object} http_util.Response
// @Failure 500 {object} http_util.Response
// @Router /user/login [post]
func (ur userRoutes) login(context *gin.Context) {
	//登录获取返回 jwt_token
	//1.获取参数
	//2.参数校验
	//3.根据用户名和密码查询用户信息
	//4.生成token
	//5.返回token
	var userReq doLoginRequest
	if err := context.ShouldBindJSON(&userReq); err != nil {
		http_util.Error(context, app_code.WithError(app_code.ErrorBadRequest, err))
		return
	}

	user, err := ur.u.Login(context.Request.Context(), userReq.Username, userReq.Password)
	if err != nil {
		http_util.Error(context, err)
		return
	}

	token, err := ur.a.GenerateToken(user.Username, user.ID)
	if err != nil {
		http_util.Error(context, err)
		return
	}

	context.JSON(http.StatusOK, doLoginResponse{Token: token, UserID: user.ID, Username: user.Username})
}

type doRegisterRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone"`
	Setting  entity.UserSetting
}

type doRegisterResponse struct {
	UserID int `json:"userId"`
}

// @Summary 注册
// @Description 注册
// @Tags 用户
// @Accept json
// @Produce json
// @Param jsonBody body doRegisterRequest true "注册参数"
// @Success 200 {object} doRegisterRequest
// @Failure 400 {object} http_util.Response
// @Failure 500 {object} http_util.Response
// @Router /user/register [post]
func (ur userRoutes) register(context *gin.Context) {
	// 注册
	// 1.获取参数
	// 2.参数校验
	// 3.根据用户名查询用户信息
	// 4.保存用户信息
	// 5.返回用户ID
	var userReq doRegisterRequest
	if err := context.ShouldBindJSON(&userReq); err != nil {
		http_util.Error(context, app_code.WithError(app_code.ErrorBadRequest, err))
		return
	}

	ext, _ := json.Marshal(userReq.Setting)
	user := entity.User{
		Username: userReq.UserName,
		Password: userReq.Password,
		Phone:    userReq.Phone,
		Ext:      string(ext),
	}

	user, err := ur.u.Register(context.Request.Context(), user)
	if err != nil {
		http_util.Error(context, err)
		return
	}

	context.JSON(http.StatusOK, doRegisterResponse{UserID: user.ID})
}

type doUpdateUserSettingRequest struct {
	UserSetting entity.UserSetting `json:"userSetting" binding:"required"`
	UserID      int                `json:"userID" binding:"required"`
}

// @Summary 更新用户设置
// @Description 更新用户设置
// @Tags 用户
// @Accept json
// @Produce json
// @Param jsonBody body entity.UserSetting true "用户设置"
// @Success 200 {object} http_util.Response
// @Failure 400 {object} http_util.Response
// @Failure 500 {object} http_util.Response
// @Router /user/updateSetting [post]
func (ur userRoutes) updateSetting(c *gin.Context) {
	var request entity.UserSetting
	if err := c.ShouldBindJSON(&request); err != nil {
		http_util.Error(c, app_code.WithError(app_code.ErrorBadRequest, err))
		return
	}
	userID := http_util.GetUserID(c)

	err := ur.u.UpdateSetting(c.Request.Context(), userID, request)
	if err != nil {
		http_util.Error(c, app_code.WithError(app_code.ErrorUpdateUserSetting, err))
		return
	}
	http_util.Success(c, nil)
}

// @Summary 获取用户设置
// @Description 获取用户设置
// @Tags 用户
// @Accept json
// @Produce json
// @Param userID query int true "用户ID"
// @Success 200 {object} http_util.Response{data=entity.UserSetting}
// @Failure 400 {object} http_util.Response
// @Failure 500 {object} http_util.Response
// @Router /user/getSetting [get]
func (ur userRoutes) getSetting(c *gin.Context) {
	userID := http_util.GetUserID(c)
	setting, err := ur.u.GetSettingByUserID(c.Request.Context(), userID)
	if err != nil {
		http_util.Error(c, app_code.WithError(app_code.ErrorGetUserSetting, err))
		return
	}
	c.JSON(http.StatusOK, setting)
}
