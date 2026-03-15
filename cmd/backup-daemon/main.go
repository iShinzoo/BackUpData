package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/iShinzoo/BackUpData/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedBackupServiceServer
}

func (s *server) StreamProgress(
	req *pb.BackupRequest,
	stream pb.BackupService_StreamProgressServer,
) error {

	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("Progress %d%%", i*20)

		stream.Send(&pb.ProgressUpdate{
			Message: msg,
		})

		time.Sleep(time.Second)
	}

	return nil
}

func (s *server) RunBackup(ctx context.Context, req *pb.BackupRequest) (*pb.BackupResponse, error) {

	log.Println("Backup requested for:", req.Database)

	return &pb.BackupResponse{
		Status: "backup started",
	}, nil
}

func main() {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterBackupServiceServer(grpcServer, &server{})

	log.Println("Backup Daemon running on port 50051")

	grpcServer.Serve(lis)
}
