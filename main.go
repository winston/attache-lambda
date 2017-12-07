package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rwcarlsen/goexif/exif"
)

type uploadResponse struct {
	Path        string
	ContentType string
	Bytes       int
	Meta        uploadMeta
}

type uploadMeta struct {
	DateTime string
	LatLong  string
	Geometry string
}

type uploadServer struct {
	region string
	bucket string
}

func main() {
	http.Handle("/", uploadServer{region: os.Getenv("AWS_REGION"), bucket: os.Getenv("AWS_BUCKET")})

	log.Printf("Listening to %s...", os.Getenv("PORT"))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (s uploadServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST", "PUT", "PATCH":
		file, fileType, fileLen := read(w, r)
		filePath := sendToS3(s.region, s.bucket, file, fileType)
		fileMeta := meta(file, fileType)

		log.Println(fileMeta)

		result := uploadResponse{
			Path:        filePath,
			ContentType: fileType,
			Bytes:       fileLen,
			Meta:        fileMeta,
		}

		json.NewEncoder(w).Encode(result)

	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, OPTIONS")

	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func read(w http.ResponseWriter, r *http.Request) (*bytes.Reader, string, int) {
	stream := &bytes.Buffer{}
	_, err := io.Copy(stream, r.Body)
	if err != nil {
		log.Println(err.Error())
	}

	file := bytes.NewReader(stream.Bytes())
	fileType := http.DetectContentType(stream.Bytes())
	fileLen := stream.Len()

	return file, fileType, fileLen
}

func sendToS3(region string, bucket string, file *bytes.Reader, fileType string) string {
	fileName := filename(fileType)
	filePath := fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s", region, bucket, fileName)

	s3Service := s3.New(session.New())
	s3Options := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Body:   file,
		Key:    aws.String(fileName),
	}

	_, err := s3Service.PutObject(s3Options)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Println(awsErr.Error())
		} else {
			log.Println(err.Error())
		}
	}

	return filePath
}

func filename(fileType string) string {
	// Sorts in Reverse Chrono Order
	key := strconv.FormatInt((math.MaxInt64 - time.Now().UnixNano()), 10)
	ext := strings.TrimPrefix(fileType, "image/")

	name := fmt.Sprintf("%s.%s", key, ext)

	return name
}

func meta(file *bytes.Reader, fileType string) uploadMeta {
	fileMeta := uploadMeta{DateTime: "", LatLong: "", Geometry: ""}

	if strings.Contains(fileType, "image") {
		imageMeta(file, &fileMeta)
	}

	return fileMeta
}

func imageMeta(file *bytes.Reader, fileMeta *uploadMeta) {
	x, err := exif.Decode(file)
	if err != nil {
		log.Println(err.Error())
	} else {
		xDateTime, xerr := x.DateTime()
		if xerr != nil {
			log.Println(xerr.Error())
		}
		fileMeta.DateTime = xDateTime.String()

		xLat, xLong, xerr := x.LatLong()
		if xerr != nil {
			log.Println(xerr.Error())
		}
		fileMeta.LatLong = fmt.Sprintf("%fx%f", xLat, xLong)
	}
	file.Seek(0, 0)

	imageSrc, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Println(err.Error())
	}
	file.Seek(0, 0)

	fileMeta.Geometry = fmt.Sprintf("%dx%d", imageSrc.Width, imageSrc.Height)
}
