package attache

import (
	"bytes"

	"golang.org/x/net/context"
)

type Store interface {
	Upload(ctx context.Context, file *bytes.Reader, fileType string) (filePath string, err error)
}
