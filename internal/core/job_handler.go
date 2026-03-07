package core

import (
	"context"
	"fmt"
	"time"
)

func BackupHandler(ctx context.Context, job BackupJob) BackupResult {

	fmt.Println("Starting backup:", job.Name)

	time.Sleep(3 * time.Second)

	fmt.Println("Finished backup:", job.Name)

	return BackupResult{
		Name:   job.Name,
		Status: "success",
	}
}
