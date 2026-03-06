package postgres

import (
	"context"
	"io"

	"github.com/iShinzoo/BackUpData/internal/core"
)

type PostgresAdapter struct {
	connString string
}

func New(conn string) *PostgresAdapter {
	return &PostgresAdapter{
		connString: conn,
	}
}

func (p *PostgresAdapter) Connect(ctx context.Context) error {
	return nil
}

func (p *PostgresAdapter) Backup(ctx context.Context, opts core.BackUpOptions) (io.ReadCloser, error) {
	return nil, nil
}

func (p *PostgresAdapter) Restore(ctx context.Context, r io.Reader) error {
	return nil
}
