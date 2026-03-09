package postgres

import (
	"context"
	"io"
	"os/exec"
)

func (p *PostgresAdapter) RunDump(ctx context.Context, dbURL string) (io.ReadCloser, error) {

	cmd := exec.CommandContext(
		ctx,
		"docker",
		"exec",
		"backup-postgres",
		"pg_dump",
		dbURL,
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	go cmd.Wait()

	return stdout, nil
}
