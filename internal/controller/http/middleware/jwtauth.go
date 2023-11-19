package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gw123/glog"
	"github.com/mytoolzone/task-mini-program/internal/app_code"
	"github.com/mytoolzone/task-mini-program/internal/controller/http/http_util"
	"github.com/mytoolzone/task-mini-program/pkg/auth"
	"strings"
	"time"
)

// JWT 自定义中间件
func JWT(authH auth.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			code  app_code.CodeType = app_code.Success
			claim *auth.Claims
			err   error
		)

		token := c.GetHeader("Authorization")
		if token == "" {
			code = app_code.ErrorTokenNotSet
		} else {
			if strings.HasPrefix(token, "Bearer ") {
				token = token[7:]
			}
			// 解析token
			claim, err = authH.ParseToken(token)
			if err != nil {
				glog.Errorf("parse token failed, err: %v, token: %v", err, token)
				code = app_code.ErrorAuthFailed
			} else if time.Now().Unix() > claim.ExpiresAt {
				code = app_code.ErrorTokenTimeout
			}
		}

		if code != app_code.Success {
			http_util.Error(c, app_code.New(code, "auth failed"))
			c.Abort()
			return
		}

		glog.Infof("userId -------------- %v", claim.UserID)
		// 将解析出来的用户id放入上下文
		http_util.SetUserID(c, claim.UserID)
		http_util.SetUserName(c, claim.Username)

		c.Next()
	}
}
