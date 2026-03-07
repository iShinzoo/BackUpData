package core

import "context"

type FullBackupStrategy struct{}

func (f *FullBackupStrategy) Execute(
	ctx context.Context,
	job BackupJob,
) BackupResult {

	return BackupResult{
		Name:   job.Name,
		Status: "full backup-done",
	}
}
