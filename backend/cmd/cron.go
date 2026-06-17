package cmd

import (
	"go-project-template/cmd/job/testjob"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

// cronCmd 定时任务
func cronCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cron",
		Short: "定时任务",
		Run: func(cmd *cobra.Command, args []string) {
			c := cron.New()

			c.AddJob("@every 6s", cron.NewChain(
				cron.Recover(cron.DefaultLogger), cron.SkipIfStillRunning(cron.DefaultLogger),
			).Then(testjob.NewJob()))

			c.Start()
		},
	}
}
