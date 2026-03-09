package core

import "context"

type BackupExecutor interface {
	Run(ctx context.Context, job BackupJob) BackupResult
}

func BackupHandler(
	ctx context.Context,
	job BackupJob,
	executor BackupExecutor,
) BackupResult {

	return executor.Run(ctx, job)
}