package core

import "context"

type BackupService struct {
	Runner func(context.Context, BackupJob) BackupResult
}

func (s *BackupService) Handle(ctx context.Context, job BackupJob) BackupResult {
	return s.Runner(ctx, job)
}
