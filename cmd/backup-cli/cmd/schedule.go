package cmd

import (
	"github.com/iShinzoo/BackUpData/internal/scheduler"
	"github.com/iShinzoo/BackUpData/pkg/logger"
	"github.com/spf13/cobra"
)

var scheduleCmd = &cobra.Command{
	Use: "schedule",
	Run: func(cmd *cobra.Command, agrs []string) {

		logg, err := logger.New()
		if err != nil {
			panic(err)
		}

		s := scheduler.New()

		s.AddJob("*/10 * * * * *", func() {

			logg.Info("Scheduler running...")

		})

		s.Start()

		select {}
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
}
