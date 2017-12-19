package gcloudstore

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"

	"cloud.google.com/go/storage"

	"github.com/pkg/errors"
	attache "github.com/winston/attache-lambda"
)

// Store uses Google Cloud Storage as backing store
type Store struct {
	bucketName string
}

// NewStore returns *Store setup and ready to use
func NewStore(bucketName string) Store {
	return Store{
		bucketName: bucketName,
	}
}

// Upload fulfills attache.Store interface
func (s Store) Upload(ctx context.Context, src io.ReadSeeker, fileType string) (string, error) {
	fileName := filename(fileType)

	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "storage newclient")
	}
	dst := client.Bucket(s.bucketName).Object(fileName).NewWriter(ctx)
	if _, err := io.Copy(dst, src); err != nil {
		return "", errors.Wrapf(err, "io copy")
	}
	// must call `dst.Close()` otherwise file won't be available
	if err := dst.Close(); err != nil {
		return "", errors.Wrapf(err, "writer close")
	}

	return fileName, nil
}

func filename(fileType string) string {
	// Sorts in Reverse Chrono Order
	key := strconv.FormatInt((math.MaxInt64 - time.Now().UnixNano()), 10)
	ext := strings.TrimPrefix(fileType, "image/")

	name := fmt.Sprintf("%s.%s", key, ext)

	return name
}

// compile-time check that we implement attache.Store interface
var _ attache.Store = Store{}
