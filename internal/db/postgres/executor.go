package postgres

import (
	"context"
	"os"
	"path/filepath"

	"github.com/iShinzoo/BackUpData/internal/compression"
	"github.com/iShinzoo/BackUpData/internal/core"
	"github.com/iShinzoo/BackUpData/internal/storage/local"
	"github.com/iShinzoo/BackUpData/pkg/logger"
	"go.uber.org/zap"
)

type Executor struct{}

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

	wd, err := os.Getwd()
	if err != nil {
		return core.BackupResult{
			Name:  job.Name,
			Error: err,
		}
	}

	backupDir := filepath.Join(wd, "backups")

	if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
		return core.BackupResult{
			Name:  job.Name,
			Error: err,
		}
	}

	filePath := filepath.Join(backupDir, job.Name+".sql.gz")

	fileWriter, err := local.CreateFile(filePath)
	if err != nil {
		return core.BackupResult{
			Name:  job.Name,
			Error: err,
		}
	}

	defer fileWriter.Close()

	err = compression.CompressStream(dumpStream, fileWriter)
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
