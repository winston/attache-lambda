package attache

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"
)

func (s Server) handleUpload(w http.ResponseWriter, r *http.Request) (result uploadResponse, err error) {
	stream := &bytes.Buffer{}
	_, err = io.Copy(stream, r.Body)
	if err != nil {
		return result, errors.Wrapf(err, "copy")
	}
	defer r.Body.Close()

	file := bytes.NewReader(stream.Bytes())
	fileType := http.DetectContentType(stream.Bytes())
	fileLen := stream.Len()

	filePath, err := s.Storage.Upload(r.Context(), file, fileType)
	if err != nil {
		return result, errors.Wrapf(err, "upload")
	}

	fileMeta := extractFileMeta(file, fileType)

	result = uploadResponse{
		Path:        filePath,
		ContentType: fileType,
		Bytes:       fileLen,
		Meta:        fileMeta,
	}

	return result, nil
}

func extractFileMeta(file *bytes.Reader, fileType string) uploadMeta {
	fileMeta := uploadMeta{}

	if strings.HasPrefix(fileType, "image/") {
		imageMeta(file, &fileMeta)
	}

	return fileMeta
}

func imageMeta(file *bytes.Reader, fileMeta *uploadMeta) {
	file.Seek(0, 0)
	x, err := exif.Decode(file)
	if err == nil {
		xDateTime, _ := x.DateTime()
		fileMeta.DateTime = xDateTime.String()

		xLat, xLong, _ := x.LatLong()
		fileMeta.LatLong = fmt.Sprintf("%fx%f", xLat, xLong)
	}

	file.Seek(0, 0)
	imageSrc, _, err := image.DecodeConfig(file)
	if err == nil {
		fileMeta.Geometry = fmt.Sprintf("%dx%d", imageSrc.Width, imageSrc.Height)
	}
}
