package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gw123/glog"
	"github.com/mytoolzone/task-mini-program/internal/app_code"
	"github.com/mytoolzone/task-mini-program/internal/controller/http/http_util"
	"github.com/mytoolzone/task-mini-program/internal/entity"
	"github.com/mytoolzone/task-mini-program/internal/usecase"
)

var adminPaths = []string{
	"/v1/task/auditTask",
	"/v1/task/assignRole",
	"/v1/task/auditUserTask",
	"/v1/task/auditApplyTask",
}

// CheckRole 检查用户角色是否有权限
func CheckRole(userCase usecase.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		userID := http_util.GetUserID(c)
		var role = entity.UserRoleMember

		roleModel, err := userCase.GetUserRole(c, userID)
		if err != nil {
			glog.WithErr(err).Error("获取用户角色失败")
		} else {
			role = roleModel.Role
		}

		http_util.SetUserRole(c, role)
		// 如果path需要管理员权限检查用户是否有管理员权限
		if stringInSlice(path, adminPaths) {
			if role == entity.UserRoleAdmin {
				c.Next()
			} else {
				http_util.Error(c, app_code.New(app_code.ErrorForbidden, "没有权限"))
			}
		} else {
			c.Next()
		}
	}
}

func stringInSlice(path string, paths []string) bool {
	for _, v := range paths {
		if v == path {
			return true
		}
	}
	return false
}
