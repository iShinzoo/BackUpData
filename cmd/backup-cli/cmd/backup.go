package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/iShinzoo/BackUpData/pkg/logger"
	pb "github.com/iShinzoo/BackUpData/proto"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Execute database backup",
	Run: func(cmd *cobra.Command, args []string) {

		logg, err := logger.New()
		if err != nil {
			panic(err)
		}

		defer logg.Sync()

		logg.Info("Starting backup Command")

		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			logg.Fatal("Failed to connect to daemon", zap.Error(err))
		}

		defer conn.Close()

		logg.Info("Connected to backup daemon")

		client := pb.NewBackupServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		logg.Info("Sending backup request", zap.String("database", "postgres-db"))

		resp, err := client.RunBackup(ctx, &pb.BackupRequest{
			Database: "postgres-db",
		})

		if err != nil {
			logg.Fatal("backup request failed", zap.Error(err))
		}

		logg.Info(
			"Backup Daemon response received",
			zap.String("Status", resp.Status),
		)

		fmt.Println("Daemon Response:", resp.Status)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
