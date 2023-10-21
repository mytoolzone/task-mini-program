// Package server-cmd configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/mytoolzone/task-mini-program/config"
	v1 "github.com/mytoolzone/task-mini-program/internal/controller/http/v1"
	"github.com/mytoolzone/task-mini-program/internal/usecase"
	"github.com/mytoolzone/task-mini-program/internal/usecase/repo"
	"github.com/mytoolzone/task-mini-program/internal/usecase/webapi"
	"github.com/mytoolzone/task-mini-program/pkg/httpserver"
	"github.com/mytoolzone/task-mini-program/pkg/logger"
	"github.com/mytoolzone/task-mini-program/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("server - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use case
	translationUseCase := usecase.New(
		repo.New(pg),
		webapi.New(),
	)

	userUseCase := usecase.NewUserUseCase(repo.NewUserRepo(pg))

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, translationUseCase, userUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("server - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("server - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("server - Run - httpServer.Shutdown: %w", err))
	}
}
