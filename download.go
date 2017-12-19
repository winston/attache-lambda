package attache

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

func (s Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	fullpath := strings.TrimPrefix(r.URL.RequestURI(), s.GetPrefixPath)
	objectKey := filepath.Base(fullpath)
	stream, err := s.Storage.Download(r.Context(), objectKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if stream == nil {
		http.Error(w, fullpath, http.StatusNotFound)
		return
	}
	io.Copy(w, stream)
}
