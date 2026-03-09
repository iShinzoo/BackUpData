package local

import (
	"context"
	"io"
	"os"
)

func Save(ctx context.Context, path string, r io.Reader) error {

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, r)

	return err
}

func CreateFile(path string) (*os.File, error) {
	return os.Create(path)
}
