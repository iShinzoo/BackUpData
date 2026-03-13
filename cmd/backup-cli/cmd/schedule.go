package cmd

import (
	"fmt"

	"github.com/iShinzoo/BackUpData/internal/scheduler"
	"github.com/spf13/cobra"
)

var scheduleCmd = &cobra.Command{
	Use:   "Schedule",
	Short: "Run Scheduled backups",
	Run: func(cmd *cobra.Command, agrs []string) {

		s := scheduler.New()

		s.AddJob("0 */6 * * *", func() {
			fmt.Println("Scheduler running...")
		})

		s.Start()

		select {}
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
}
