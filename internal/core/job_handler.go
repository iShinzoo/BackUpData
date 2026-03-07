package core

import "context"

func BackupHandler(ctx context.Context, job BackupJob) BackupResult {

	return BackupResult{
		Name:   job.Name,
		Status: "success",
	}
}
