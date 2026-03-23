package cmd

import (
	"context"
	"time"

	"github.com/iShinzoo/BackUpData/internal/scheduler"
	"github.com/iShinzoo/BackUpData/pkg/logger"
	pb "github.com/iShinzoo/BackUpData/proto"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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

			logg.Info("Scheduler backup triggered...")

			conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

			if err != nil {
				logg.Error("Failed to connect to daemon", zap.Error(err))
				return
			}

			defer conn.Close()

			client := pb.NewBackupServiceClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			resp, err := client.RunBackup(ctx, &pb.BackupRequest{
				Database: "postgres-db",
			})

			if err != nil {
				logg.Error("backup failed", zap.Error(err))
				return
			}

			logg.Info("backup completed via scheduler", zap.String("status", resp.Status))
		})

		s.Start()

		select {}
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
}
