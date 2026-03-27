package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fmt"

	"github.com/iShinzoo/BackUpData/internal/core"
	"github.com/iShinzoo/BackUpData/internal/core/worker"
	"github.com/iShinzoo/BackUpData/internal/db/postgres"
	"github.com/iShinzoo/BackUpData/internal/metrics"
	"github.com/iShinzoo/BackUpData/internal/notification/slack"
	"github.com/iShinzoo/BackUpData/pkg/config"
	"github.com/iShinzoo/BackUpData/pkg/logger"
	pb "github.com/iShinzoo/BackUpData/proto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedBackupServiceServer
	cfg    *config.Config
	logger *zap.Logger
}

func (s *server) StreamProgress(
	req *pb.BackupRequest,
	stream pb.BackupService_StreamProgressServer,
) error {

	for i := 0; i < 5; i++ {

		s.logger.Info(
			"backup progress",
			zap.Int("progress_percent", i*20),
		)

		err := stream.Send(&pb.ProgressUpdate{
			Message: fmt.Sprintf("Progress: %d%%", i*20),
		})

		if err != nil {
			return err
		}

		time.Sleep(time.Second)
	}

	return nil
}

func (s *server) RunBackup(ctx context.Context, req *pb.BackupRequest) (*pb.BackupResponse, error) {

	s.logger.Info(
		"backup requested",
		zap.String("database", req.Database),
	)

	jobs := make(chan core.BackupJob, 10)
	results := make(chan core.BackupResult)

	pool := worker.WorkerPool{
		Workers: 3,
	}

	pgExecutor := postgres.Executor{}

	var notifier core.Notifier

	if s.cfg.SlackWebhook != "" {
		notifier = slack.New(s.cfg.SlackWebhook)
	}

	handler := func(ctx context.Context, job core.BackupJob) core.BackupResult {
		return core.BackupHandler(ctx, job, &pgExecutor, notifier)
	}

	go pool.Run(ctx, jobs, results, handler)

	jobs <- core.BackupJob{
		Name: fmt.Sprintf(
			"%s-%s",
			req.Database,
			time.Now().Format("20060102-150405"),
		),
		URL: s.cfg.PostgresURL,
	}

	close(jobs)

	res := <-results

	if res.Error != nil {

		s.logger.Error(
			"backup failed",
			zap.String("database", req.Database),
			zap.Error(res.Error),
		)

		return &pb.BackupResponse{
			Status: "backup failed",
		}, res.Error
	}

	s.logger.Info(
		"backup completed",
		zap.String("database", req.Database),
	)

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

	logg, _ := logger.New()

	logg.Info(
		"grpc request received",
		zap.String("method", info.FullMethod),
	)

	return handler(ctx, req)
}

func main() {

	logg, err := logger.New()
	if err != nil {
		panic(err)
	}

	defer logg.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logg.Fatal("failed to load config:", zap.Error(err))
	}

	if err := cfg.Validate(); err != nil {
		logg.Fatal("invalid configuration:", zap.Error(err))
	}

	metrics.Register()

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
		logg.Fatal("", zap.Error(err))
	}

	pb.RegisterBackupServiceServer(grpcServer, &server{
		cfg:    cfg,
		logger: logg,
	})

	go func() {
		logg.Info("Backup Daemon running on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			logg.Fatal("", zap.Error(err))
		}
	}()

	go func() {
		logg.Info("Prometheus metrics exposed on :9090/metrics")

		http.Handle("/metrics", promhttp.Handler())

		if err := http.ListenAndServe(":9090", nil); err != nil {
			logg.Fatal("metrics server error", zap.Error(err))
		}
	}()

	<-ctx.Done()

	logg.Info("Shutdown signal received")

	grpcServer.GracefulStop()
}
