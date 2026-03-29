package postgres

import (
	"context"
	"io"

	"github.com/iShinzoo/BackUpData/internal/compression"
	"github.com/iShinzoo/BackUpData/internal/core"
	"github.com/iShinzoo/BackUpData/pkg/logger"
	"go.uber.org/zap"
)

type Executor struct {
	storage core.Storage
}

func NewExecutor(storage core.Storage) *Executor {
	return &Executor{
		storage: storage,
	}
}

func (e *Executor) Run(
	ctx context.Context,
	job core.BackupJob,
) core.BackupResult {

	logg, err := logger.New()
	if err != nil {
		panic(err)
	}

	logg.Info("Starting backup:", zap.String("job_name", job.Name))

	pg := New(job.URL)

	dumpStream, cmd, err := pg.RunDump(ctx, job.URL)
	if err != nil {
		return core.BackupResult{
			Name:  job.Name,
			Error: err,
		}
	}

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		err := compression.CompressStream(dumpStream, pw)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	err = e.storage.Save(ctx, job.Name+".sql.gz", pr)
	if err != nil {
		return core.BackupResult{
			Name:  job.Name,
			Error: err,
		}
	}

	dumpStream.Close()

	if err := cmd.Wait(); err != nil {
		return core.BackupResult{
			Name:  job.Name,
			Error: err,
		}
	}
	logg.Info("Finished backup:", zap.String("job_name", job.Name))

	return core.BackupResult{
		Name:   job.Name,
		Status: "success",
	}
}
