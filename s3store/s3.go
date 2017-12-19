package s3store

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/oklog/ulid"
	"github.com/pkg/errors"
	attache "github.com/winston/attache-lambda"
)

// Store uses s3 as backing store
type Store struct {
	Bucket string
}

// Upload fulfills attache.Store interface
func (s Store) Upload(ctx context.Context, file io.ReadSeeker, fileType string) (string, error) {
	fileName := filename(fileType)
	filePath := fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s", os.Getenv("AWS_REGION"), s.Bucket, fileName)

	// unsure about how long we can cache `svc` or must we really
	// session.New everytime?
	svc := s3.New(session.New())
	_, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Body:   file,
		Key:    &fileName,
	})

	return filePath, err
}

// Download fulfills attache.Store interface
func (s Store) Download(ctx context.Context, filePath string) (io.ReadCloser, error) {
	svc := s3.New(session.New())
	result, err := svc.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    &filePath,
	})
	if e, ok := err.(awserr.Error); ok {
		if e.Code() == s3.ErrCodeNoSuchKey ||
			e.Code() == s3.ErrCodeNoSuchBucket ||
			e.Code() == s3.ErrCodeNoSuchUpload {
			return nil, nil
		}
	}
	if err != nil {
		return nil, errors.Wrapf(err, "get object: %#v", filePath)
	}

	return result.Body, nil
}

func filename(fileType string) string {
	current := time.Now()
	entropy := rand.New(rand.NewSource(current.UnixNano()))
	key := ulid.MustNew(ulid.Timestamp(current), entropy)

	ext := strings.TrimPrefix(fileType, "image/")

	name := fmt.Sprintf("%s.%s", key, ext)

	return name
}

// compile-time check that we implement attache.Store interface
var _ attache.Store = Store{}
