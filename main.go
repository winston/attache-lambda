// Inspired by:
// - https://github.com/apex/up-examples/blob/master/oss/golang-uploads/main.go
// - https://github.com/choonkeat/attache/blob/master/lib/attache/upload.rb

// todos:
// 1. tests
// 2. apex logging?
// 3. validation? prevent image bomb?
// 4. protect API?
// 5. exif?
// 6. documentation

package main

import (
	// "encoding/base64"
	// "io/ioutil"
	"fmt"
	"log"
	"net/http"
	"os"
  "mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	// "github.com/apex/log"
	// "github.com/apex/log/handlers/json"
	// "github.com/apex/log/handlers/text"
)

// Use JSON logging when run by Up (including `up start`).
// func init() {
//   if os.Getenv("UP_STAGE") == "" {
//     log.SetHandler(text.Default)
//   } else {
//     log.SetHandler(json.Default)
//   }
// }

func main() {
	http.HandleFunc("/", apiRequest)

	log.Printf("Listening to %s...", os.Getenv("PORT"))
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatalf(fmt.Sprintf("Error Listening - %s", err))
	}
}

func apiRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s", r)

	// Validation

	switch r.Method {
	case "POST", "PUT", "PATCH":
		file := readFile(w, r) // is this refactoring ok?

		toS3(w, r, file)

	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, OPTIONS")

	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func readFile(w http.ResponseWriter, r *http.Request) multipart.File {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		// log.WithError(err).Error("Parsing Form")
		http.Error(w, fmt.Sprintf("Error Parsing Form - %s", err), http.StatusBadRequest)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Fprintf(w, fmt.Sprintf("Uploaded: %s, of type %s", handler.Filename, handler.Header.Get("Content-Type")))

	return file
}

func toS3(w http.ResponseWriter, r *http.Request, file multipart.File) {
	bucket := r.FormValue("bucket")
	log.Printf(bucket)

	sess := session.New()
	svc := s3manager.NewUploader(sess)

	filename := "abc"

	_, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   file})

	if err != nil {
		log.Fatalf("Unable to Upload %s to %s, %s", filename, bucket, err)
	}

	fmt.Fprintf(w, "Successfully Uploaded %s to %s", filename, bucket)
}
