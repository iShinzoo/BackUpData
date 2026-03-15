package main

import (
	"context"
	"log"
	"net"

	pb "github.com/iShinzoo/BackUpData/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedBackupServiceServer
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
