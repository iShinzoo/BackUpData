package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var scheduleCmd = &cobra.Command{
	Use:   "Schedule",
	Short: "Run Scheduled backups",
	Run: func(cmd *cobra.Command, agrs []string) {

		fmt.Println("Scheduler running...")
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
}
