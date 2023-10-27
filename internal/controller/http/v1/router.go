// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/mytoolzone/task-mini-program/internal/controller/http/middleware"
	"github.com/mytoolzone/task-mini-program/pkg/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	_ "github.com/mytoolzone/task-mini-program/docs"
	"github.com/mytoolzone/task-mini-program/internal/usecase"
	"github.com/mytoolzone/task-mini-program/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title       Task Mini Program
// @description Task Mini Program API
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(handler *gin.Engine,
	l logger.Interface,
	u usecase.User,
	tk usecase.Task,
	n usecase.Notice,
	auth auth.Auth,
) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	hl := handler.Group("/v1")
	jwt := middleware.JWT(auth)
	{
		newTaskRoutes(hl, jwt, tk)
		newUserRoutes(hl, jwt, auth, u)
		newNoticeRoutes(hl, jwt, n)
	}
}
