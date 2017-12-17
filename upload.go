package attache

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

func (s Server) handleUpload(w http.ResponseWriter, r *http.Request) (result uploadResponse, err error) {
	// 1. dump into tempfile
	file, err := ioutil.TempFile(os.TempDir(), "upload")
	if err != nil {
		return result, errors.Wrapf(err, "tempfile")
	}
	_, err = io.Copy(file, r.Body)
	if err != nil {
		return result, errors.Wrapf(err, "copy")
	}
	defer os.Remove(file.Name())
	defer r.Body.Close()

	stat, err := file.Stat()
	if err != nil {
		return result, errors.Wrapf(err, "stat")
	}

	file.Seek(0, 0)
	uniqueKey, err := s.Storage.Upload(file)
	if err != nil {
		return result, errors.Wrapf(err, "upload")
	}

	file.Seek(0, 0)
	first512Bytes := make([]byte, 512)
	_, err = io.ReadFull(file, first512Bytes)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return result, errors.Wrapf(err, "readfull")
	}

	result = uploadResponse{
		Path:        uniqueKey,
		ContentType: http.DetectContentType(first512Bytes),
		Bytes:       stat.Size(),
	}

	// rotate
	// exif
	file.Seek(0, 0)
	img, _, err := image.DecodeConfig(file)
	if err == nil {
		result.Geometry = fmt.Sprintf("%dx%d", img.Width, img.Height)
	}

	return result, nil
}
