package core

import (
	"context"
	"fmt"
	"time"
)

type BackupExecutor interface {
	Run(ctx context.Context, job BackupJob) BackupResult
}

func BackupHandler(
	ctx context.Context,
	job BackupJob,
	executor BackupExecutor,
	notifer Notifier,
) BackupResult {

	start := time.Now()

	result := executor.Run(ctx, job)

	duration := time.Since(start)

	var msg string

	if result.Error != nil {

		msg = fmt.Sprintf(
			"❌ Backup Failed\nDatabase: %s\nError: %v",
			job.Name,
			result.Error,
		)

	} else {

		msg = fmt.Sprintf(
			"✅ Backup Completed\nDatabase: %s\nDuration: %s\nSize: %d bytes",
			job.Name,
			duration,
			result.Size,
		)
	}

	if notifer != nil {
		notifer.Notify(ctx, msg)
	}

	return result
}
