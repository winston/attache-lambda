package attache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"golang.org/x/net/context"
)

func TestHandleDownload(t *testing.T) {
	testCases := []struct {
		givenURI         string
		givenFile        string
		givenContentType string
		expectedStatus   int
		expectedFile     string
	}{
		{
			givenURI:         "/%s",
			givenFile:        "testdata/transparent.gif",
			givenContentType: "image/gif",
			expectedStatus:   http.StatusOK,
			expectedFile:     "testdata/transparent.gif",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			file, err := os.Open(tc.givenFile)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer file.Close()
			store := newDummyStore()
			store.Upload(context.Background(), file, tc.givenContentType)

			r := httptest.NewRequest("GET", fmt.Sprintf(tc.givenURI, store.LastUniqueKey), nil)
			w := httptest.NewRecorder()
			s := Server{Storage: store}
			s.ServeHTTP(w, r)

			result := w.Result()
			assert.Equal(t, tc.expectedStatus, result.StatusCode)

			expectedBytes, err := ioutil.ReadFile(tc.expectedFile)
			if err != nil {
				t.Fatal(err.Error())
			}
			actualBytes, err := ioutil.ReadAll(result.Body)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer result.Body.Close()
			assert.Equal(t, expectedBytes, actualBytes)
		})
	}
}
