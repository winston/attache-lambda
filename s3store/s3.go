package s3store

import (
	"io"
	"math"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Store uses s3 as backing store
type Store struct {
	Bucket string
}

// Upload fulfills attache.Store interface
func (s Store) Upload(r io.ReadSeeker) (string, error) {
	uniqueKey := strconv.FormatInt((math.MaxInt64 - time.Now().UnixNano()), 10)
	// unsure about how long we can cache `svc` or must we really
	// session.New everytime?
	svc := s3.New(session.New())
	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Body:   r,
		Key:    &uniqueKey,
	})
	return uniqueKey, err
}
