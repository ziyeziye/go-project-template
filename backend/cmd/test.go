package cmd

import (
	"go-project-template/cmd/job/testjob"
	"go-project-template/config"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

func testCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "test",
		Long:  `test`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("test")
			spew.Dump(config.Get())
			testjob.NewJob().Run()
		},
	}
}
