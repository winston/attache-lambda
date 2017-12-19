package attache

import (
	"io"
	"net/http"
	"path/filepath"
)

func (s Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Base(r.URL.Path)
	stream, err := s.Storage.Download(r.Context(), filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if stream == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	io.Copy(w, stream)
}
