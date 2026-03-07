package core

import "context"

type BackupStrategy interface {
	Execute(ctx context.Context, job BackupJob) BackupResult
}
