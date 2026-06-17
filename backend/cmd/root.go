package cmd

import (
	"log"
	"time"

	"go-project-template/app"
	"go-project-template/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// 配置文件
	envFile string

	rootCmd = &cobra.Command{
		Use:   "info",
		Short: "info",
		Long:  `info`,
		Run: func(cmd *cobra.Command, args []string) {
			// 定时任务
			go cronCmd().Run(cmd, args)
			app.RunServe()

			select {}
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	cobra.OnInitialize(afterInit)
	rootCmd.PersistentFlags().StringVarP(&envFile, "env", "e", "", "env file (default .env)")
	// rootCmd.AddCommand(cronCmd())
	rootCmd.AddCommand(testCmd())
}

func afterInit() {
	if envFile == "" {
		envFile = ".env"
	}

	viper.SetConfigFile(envFile)
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	loc, _ := time.LoadLocation("UTC")
	time.Local = loc

	if _, err := config.Load(); err != nil {
		log.Fatalf("load config: %v", err)
	}

	// 迁移
	// err := db.Engine().AutoMigrate(
	// 	model.User{},
	// )
	// if err != nil {
	// 	panic(err)
	// }
}
