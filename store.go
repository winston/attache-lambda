package attache

import (
	"io"

	"golang.org/x/net/context"
)

type Store interface {
	Upload(ctx context.Context, file io.ReadSeeker, fileType string) (filePath string, err error)
	Download(ctx context.Context, filePath string) (io.ReadCloser, error) // should return `nil, nil` for when `filePath` is not found
}
