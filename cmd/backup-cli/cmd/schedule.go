package cmd

import (
	"fmt"

	"github.com/iShinzoo/BackUpData/internal/scheduler"
	"github.com/spf13/cobra"
)

var scheduleCmd = &cobra.Command{
	Use: "schedule",
	Run: func(cmd *cobra.Command, agrs []string) {

		s := scheduler.New()

		s.AddJob("*/10 * * * * *", func() {
			fmt.Println("Scheduler running...")
		})

		s.Start()

		select {}
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
}
