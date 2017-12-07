package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	http.Handle("/", uploadServer{})

	log.Printf("Listening to %s...", os.Getenv("PORT"))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (s uploadServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST", "PUT", "PATCH":
		buffer := upload(w, r)

		svc := s3.New(session.New())
		s3i := &s3.PutObjectInput{
			Body:   buffer,
			Bucket: aws.String(),
			Key:    aws.String(),
		}

		result, err := svc.PutObject(s3i)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				fmt.Println(aerr.Error())
			}
			return
		}

		fmt.Println(result)

		result := uploadResponse{
			Path:        fmt.Sprintf("some/path/%s", r.URL.Query().Get("file")),
			ContentType: http.DetectContentType(buffer.Bytes()),
			Bytes:       buffer.Len(),
			// Geometry: "4x3"
		}

		json.NewEncoder(w).Encode(result)

	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, OPTIONS")

	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func upload(w http.ResponseWriter, r *http.Request) *bytes.Buffer {
	buffer := &bytes.Buffer{}
	_, err := io.Copy(buffer, r.Body)
	if err != nil {
		log.Printf(err.Error())
	}

	return buffer
}
