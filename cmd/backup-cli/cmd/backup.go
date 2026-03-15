package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/iShinzoo/BackUpData/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Execute database backup",
	Run: func(cmd *cobra.Command, args []string) {

		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}

		defer conn.Close()

		client := pb.NewBackupServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		resp, err := client.RunBackup(ctx, &pb.BackupRequest{
			Database: "postgres-db",
		})

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Daemon Response:", resp.Status)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
