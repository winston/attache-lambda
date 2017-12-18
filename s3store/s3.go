package s3store

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	attache "github.com/winston/attache-lambda"
)

// Store uses s3 as backing store
type Store struct {
	Bucket string
}

// Upload fulfills attache.Store interface
func (s Store) Upload(file *bytes.Reader, fileType string) (string, error) {
	fileName := filename(fileType)
	filePath := fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s", os.Getenv("AWS_REGION"), s.Bucket, fileName)

	// unsure about how long we can cache `svc` or must we really
	// session.New everytime?
	svc := s3.New(session.New())
	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Body:   file,
		Key:    &fileName,
	})

	return filePath, err
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
