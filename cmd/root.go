/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/gw123/glog"
	"github.com/mytoolzone/task-mini-program/cmd/db-cmd"
	"github.com/mytoolzone/task-mini-program/cmd/server-cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "minitask",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// 详细日志
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose")
	err := viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	if err != nil {
		glog.WithErr(err).Fatal("viper.BindPFlag verbose")
		return
	}

	RootCmd.AddCommand(server_cmd.ServerCmd)
	RootCmd.AddCommand(db_cmd.GenModel)
}
