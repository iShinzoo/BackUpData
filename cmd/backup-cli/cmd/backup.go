package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/iShinzoo/BackUpData/internal/core"
	"github.com/iShinzoo/BackUpData/internal/core/worker"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Execute database backup",
	Run: func(cmd *cobra.Command, args []string) {

		ctx, cancel := signal.NotifyContext(
			context.Background(),
			os.Interrupt,
			syscall.SIGTERM,
		)

		defer cancel()

		jobs := make(chan core.BackupJob, 10)
		results := make(chan core.BackupResult)

		pool := worker.WorkerPool{
			Workers: 3,
		}

		go pool.Run(ctx, jobs, results, core.BackupHandler)

		jobs <- core.BackupJob{Name: "db1"}
		jobs <- core.BackupJob{Name: "db2"}
		jobs <- core.BackupJob{Name: "db3"}

		close(jobs)

		for r := range results {
			fmt.Println("Backup finished:", r.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
