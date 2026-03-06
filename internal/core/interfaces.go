package core

import (
	"context"
	"io"
)

type BackUpOptions struct {
	DatabaseName string
	OutputName   string
}

type Database interface {
	Connect(ctx context.Context) error
	Backup(ctx context.Context, opts BackUpOptions) (io.ReadCloser, error)
	Restore(ctx context.Context, r io.Reader) error
}

type Storage interface {
	Save(ctx context.Context, name string, r io.Reader) error
}

type Notifier interface {
	Notify(ctx context.Context, msg string) error
}
