package attache

import "io"

type Store interface {
	Upload(w io.ReadSeeker) (uniqueKey string, err error)
}
