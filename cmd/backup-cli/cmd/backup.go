package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Execute database backup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting backup process...")
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
