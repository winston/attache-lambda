package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type uploadResponse struct {
	Path        string
	ContentType string
	Bytes       int
	Geometry    *string
}

type uploadServer struct {
	bucket string
	region string
}

func (s uploadServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result := uploadResponse{
		Path: fmt.Sprintf("some/path/%s", r.URL.Query().Get("file")),
	}

	json.NewEncoder(w).Encode(result)
}
