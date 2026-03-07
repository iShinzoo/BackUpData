package worker

import (
	"context"
	"sync"

	"github.com/iShinzoo/BackUpData/internal/core"
)

type WorkerPool struct {
	Workers int
}

func (p *WorkerPool) Run(
	ctx context.Context,
	jobs <-chan core.BackupJob,
	results chan<- core.BackupResult,
	handler func(context.Context, core.BackupJob) core.BackupResult,
) {

	var wg sync.WaitGroup

	for i := 0; i < p.Workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for job := range jobs {

				select {
				case <-ctx.Done():
					return
				default:
					result := handler(ctx, job)
					results <- result
				}
			}
		}()
	}

	wg.Wait()
	close(results)
}
