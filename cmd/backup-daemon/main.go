package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
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

		err := stream.Send(&pb.ProgressUpdate{
			Message: msg,
		})

		if err != nil {
			return err
		}

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

func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	log.Println("RPC called:", info.FullMethod)

	return handler(ctx, req)
}

func main() {

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	pb.RegisterBackupServiceServer(grpcServer, &server{})

	go func() {
		log.Println("Backup Daemon running on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutdown signal received")

	grpcServer.GracefulStop()
}
