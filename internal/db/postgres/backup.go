package postgres

import (
	"context"
	"io"
	"os/exec"
)

func (p *PostgresAdapter) RunDump(ctx context.Context, dbURL string) (io.ReadCloser, *exec.Cmd, error) {

	cmd := exec.CommandContext(
		ctx,
		"docker",
		"exec",
		"-i",
		"backup-postgres",
		"pg_dump",
		"-U",
		"backup",
		"testdb",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	// log pg_dump errors if any
	go io.Copy(io.Discard, stderr)

	return stdout, cmd, nil
}
