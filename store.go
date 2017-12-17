package attache

import "bytes"

type Store interface {
	Upload(file *bytes.Reader, fileType string) (filePath string, err error)
}
