package attache

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	geometry4x3 = "4x3"
	geometry1x1 = "1x1"
)

func TestServerServeHTTP(t *testing.T) {
	testCases := []struct {
		givenURI       string
		givenFile      string
		expectedStatus int
		expectedJSON   uploadResponse
	}{
		{
			givenURI:       "/?file=x.jpg",
			givenFile:      "testdata/landscape.jpg",
			expectedStatus: http.StatusOK,
			expectedJSON: uploadResponse{
				Bytes:       425,
				ContentType: "image/jpeg",
				Geometry:    geometry4x3,
			},
		},
		{
			givenURI:       "/?file=y.jpg",
			givenFile:      "testdata/transparent.gif",
			expectedStatus: http.StatusOK,
			expectedJSON: uploadResponse{
				Bytes:       42,
				ContentType: "image/gif",
				Geometry:    geometry1x1,
			},
		},
		{
			givenURI:       "/?file=z.jpg",
			givenFile:      "testdata/sample.txt",
			expectedStatus: http.StatusOK,
			expectedJSON: uploadResponse{
				Bytes:       42,
				ContentType: "image/gif",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			input, err := os.Open(tc.givenFile)
			if err != nil {
				t.Fatalf("os.Open(%s): %s", tc.givenFile, err.Error())
			}
			defer input.Close()

			r := httptest.NewRequest("POST", tc.givenURI, input)
			w := httptest.NewRecorder()
			s := Server{Storage: newDummyStore()}
			s.ServeHTTP(w, r)

			result := w.Result()
			assert.Equal(t, tc.expectedStatus, result.StatusCode, "http status")
			if result.StatusCode != http.StatusOK {
				body, err := ioutil.ReadAll(result.Body)
				t.Fatalf("%s %#v", body, err)
			}

			var actual uploadResponse
			if err = json.NewDecoder(result.Body).Decode(&actual); err != nil {
				t.Errorf("result body: %s", err.Error())
			}
			defer result.Body.Close()

			// since `path` is uniquely generated, we can only test for "presence" first
			// then stuff it into `tc.expectedJSON` in order to perform a lazy whole-object comparison
			assert.NotEmpty(t, actual.Path, "path should not be empty")
			tc.expectedJSON.Path = actual.Path
			assert.Equal(t, tc.expectedJSON, actual)
		})
	}
}
