package postgres

import (
	"context"
	"fmt"

	"github.com/iShinzoo/BackUpData/internal/compression"
	"github.com/iShinzoo/BackUpData/internal/core"
	"github.com/iShinzoo/BackUpData/internal/storage/local"
)

type Executor struct{}

func (e *Executor) Run(
	ctx context.Context,
	job core.BackupJob,
) core.BackupResult {

	fmt.Println("Starting backup:", job.Name)

	pg := New(job.URL)

	dumpStream, err := pg.RunDump(ctx, job.URL)
	if err != nil {
		return core.BackupResult{
			Name:  job.Name,
			Error: err,
		}
	}

	filepath := "./backups/" + job.Name + ".sql.gz"

	fileWriter, err := local.CreateFile(filepath)
	if err != nil {
		return core.BackupResult{
			Name:  job.Name,
			Error: err,
		}
	}

	err = compression.CompressStream(dumpStream, fileWriter)
	if err != nil {
		return core.BackupResult{
			Name:  job.Name,
			Error: err,
		}
	}

	fmt.Println("Finished backup:", job.Name)

	return core.BackupResult{
		Name:   job.Name,
		Status: "success",
	}
}
