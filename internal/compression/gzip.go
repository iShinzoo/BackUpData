package compression

import (
	"compress/gzip"
	"io"
)

func CompressStream(input io.Reader, output io.Writer) error {

	gz := gzip.NewWriter(output)
	defer gz.Close()

	_, err := io.Copy(gz, input)

	return err
}
