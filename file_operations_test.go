package attache

import (
	"bytes"
	"io"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageMeta(t *testing.T) {
	testCases := []struct {
		givenFile    string
		expectedMeta uploadMeta
	}{
		{
			givenFile: "testdata/meta.jpg", // https://raw.githubusercontent.com/ianare/exif-samples/master/jpg/gps/DSCN0010.jpg
			expectedMeta: uploadMeta{
				DateTime: "2008-10-22 16:28:39 +0800 SGT",
				LatLong:  "43.467448x11.885127",
				Geometry: "640x480",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			log.Println("Running for", tc.givenFile)

			file, err := os.Open(tc.givenFile)
			if err != nil {
				t.Fatalf("os.Open(%s): %s", tc.givenFile, err.Error())
			}

			stream := &bytes.Buffer{}
			io.Copy(stream, file)

			result := uploadMeta{}
			ImageMeta(bytes.NewReader(stream.Bytes()), &result)

			assert.Equal(t, tc.expectedMeta, result)
		})
	}
}
