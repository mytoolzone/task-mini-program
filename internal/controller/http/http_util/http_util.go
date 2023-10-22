package http_util

import "github.com/gin-gonic/gin"

const (
	ctxUserIDKey = "ctxUserID"
)

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}

func Error(ctx *gin.Context, code string, message string) {
	ctx.JSON(200, gin.H{
		"code": code,
		"msg":  message,
	})
}

// GetUserID 从 ctx 中获取用户 id
func GetUserID(ctx *gin.Context) int {
	return ctx.GetInt("userID")
}

// SetUserID 将用户 id 设置到 ctx 中
func SetUserID(ctx *gin.Context, userID int) {
	ctx.Set(ctxUserIDKey, userID)
}
