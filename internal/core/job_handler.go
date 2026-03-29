package core

import (
	"context"
	"fmt"
	"time"

	"github.com/iShinzoo/BackUpData/internal/metrics"
)

type BackupExecutor interface {
	Run(ctx context.Context, job BackupJob) BackupResult
}

func BackupHandler(
	ctx context.Context,
	job BackupJob,
	executor BackupExecutor,
	notifer Notifier,
	storage Storage,
) BackupResult {

	start := time.Now()

	result := executor.Run(ctx, job)

	duration := time.Since(start)

	var msg string

	if result.Error != nil {

		metrics.BackupFailure.Inc()

		msg = fmt.Sprintf(
			"❌ Backup Failed\nDatabase: %s\nError: %v",
			job.Name,
			result.Error,
		)

	} else {

		metrics.BackupSuccess.Inc()

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
