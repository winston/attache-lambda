package attache

import (
	"encoding/json"
	"net/http"
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

// Server handles upload and download
type Server struct {
	Storage Store
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST", "PUT", "PATCH":
		result, err := s.handleUpload(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(result)

	case "GET":
		s.handleDownload(w, r)

	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, OPTIONS")

	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}
