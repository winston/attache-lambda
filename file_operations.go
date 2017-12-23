package attache

import (
	"bytes"
	"fmt"
	"image"

	"github.com/rwcarlsen/goexif/exif"
)

func ImageMeta(file *bytes.Reader, fileMeta *uploadMeta) {
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
