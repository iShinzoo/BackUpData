package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/iShinzoo/BackUpData/internal/core"
	"github.com/iShinzoo/BackUpData/internal/core/worker"
	"github.com/iShinzoo/BackUpData/internal/db/postgres"
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

		pgExecutor := postgres.Executor{}

		handler := func(ctx context.Context, job core.BackupJob) core.BackupResult {
			return core.BackupHandler(ctx, job, &pgExecutor)
		}

		go pool.Run(ctx, jobs, results, handler)

		jobs <- core.BackupJob{
			Name: "db1",
			URL:  "postgres://backup:backup@localhost:5432/testdb",
		}

		jobs <- core.BackupJob{
			Name: "db2",
			URL:  "postgres://backup:backup@localhost:5432/testdb",
		}

		jobs <- core.BackupJob{
			Name: "db3",
			URL:  "postgres://backup:backup@localhost:5432/testdb",
		}

		close(jobs)

		for r := range results {

			if r.Error != nil {
				fmt.Println("Backup FAILED:", r.Name, r.Error)
				continue
			}

			fmt.Println("Backup finished:", r.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
