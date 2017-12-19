// this isn't really a test against the `attache.Store` interface
// but it sets up a `dummyStore` implementation that keeps []byte
// in memory
package attache

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/net/context"

	uuid "github.com/satori/go.uuid"
)

type dummyStore struct {
	hash map[string][]byte // default is `nil`
}

func newDummyStore() *dummyStore {
	return &dummyStore{
		hash: map[string][]byte{},
	}
}

// Upload fulfills attache.Store interface
func (s *dummyStore) Upload(ctx context.Context, file *bytes.Reader, fileType string) (string, error) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	uniqueKey := uuid.NewV4().String()
	s.hash[uniqueKey] = data
	return uniqueKey, nil
}

// compile-time check that we implement attache.Store interface
var _ Store = newDummyStore()
