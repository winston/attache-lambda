package main

import (
	"bytes"
	"encoding/json"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type uploadResponse struct {
	Path        string
	ContentType string
	Bytes       int
	Geometry    *string
}

type uploadServer struct {
	bucket string
	region string
}

func main() {
	http.Handle("/", uploadServer{bucket: os.Getenv("AWS_BUCKET"), region: os.Getenv("AWS_REGION")})

	log.Printf("Listening to %s...", os.Getenv("PORT"))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (s uploadServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST", "PUT", "PATCH":
		stream, file := upload(w, r)

		s3Result := toS3(s.bucket, file)

		// rotate
		// exif

		result := uploadResponse{
			// Path:        fmt.Sprintf("some/path/%s", r.URL.Query().Get("file")),
			Path:        s3Result.String(),
			ContentType: http.DetectContentType(stream.Bytes()),
			Bytes:       stream.Len(),
			Geometry:    aws.String("4x2"),
		}

		json.NewEncoder(w).Encode(result)

	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, OPTIONS")

	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func upload(w http.ResponseWriter, r *http.Request) (*bytes.Buffer, *bytes.Reader) {
	stream := &bytes.Buffer{}
	_, err := io.Copy(stream, r.Body)
	if err != nil {
		log.Println(err.Error())
	}

	file := bytes.NewReader(stream.Bytes())

	return stream, file
}

func toS3(bucket string, file *bytes.Reader) *s3.PutObjectOutput {
	s3Service := s3.New(session.New())
	s3Options := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Body:   file,
		Key:    aws.String("abcd"),
	}

	s3Result, err := s3Service.PutObject(s3Options)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Println(awsErr.Error())
		} else {
			log.Println(err.Error())
		}
	}

	log.Println(s3Result.String())

	return s3Result
}
