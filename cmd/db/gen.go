package db

import (
	"fmt"
	"github.com/mytoolzone/task-mini-program/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
	"log"
	"os"
)

var GenModel = &cobra.Command{
	Use:   "gen-model",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines and likely contains examples`,
	Run: func(cmd *cobra.Command, args []string) {

		// Configuration
		cfg, err := config.NewConfig()
		if err != nil {
			log.Fatalf("Config error: %s", err)
		}

		if viper.GetString("pg_addr") != "" {
			cfg.PG.URL = viper.GetString("pg_addr")
		}

		pwd, _ := os.Getwd()
		g := gen.NewGenerator(gen.Config{
			OutPath: pwd + "/internal/entity/query",                                     // output path
			Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
		})

		gormdb, err := gorm.Open(postgres.Open(cfg.PG.URL))
		if err != nil {
			fmt.Errorf("gorm.Open err:%s", err.Error())
		}
		g.UseDB(gormdb) // reuse your gorm db
		g.GenerateAllTable()
		// Generate the code
		g.Execute()

		fmt.Println("gen model success")
	},
}

func init() {
	GenModel.PersistentFlags().StringP("pg_addr", "a", "", "pg_addr")

	err := viper.BindPFlag("pg_addr", GenModel.PersistentFlags().Lookup("pg_addr"))
	if err != nil {
		return
	}
}
