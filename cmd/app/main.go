package main

import (
	"log"

	"github.com/mytoolzone/task-mini-program/config"
	"github.com/mytoolzone/task-mini-program/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
