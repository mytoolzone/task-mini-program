package server_cmd

import (
	"github.com/mytoolzone/task-mini-program/config"
	"github.com/mytoolzone/task-mini-program/internal/app"
	"github.com/spf13/cobra"
	"log"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Configuration
		cfg, err := config.NewConfig()
		if err != nil {
			log.Fatalf("Config error: %s", err)
		}

		// Run
		app.Run(cfg)
	},
}
