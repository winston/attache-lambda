// Inspired by:
// - https://github.com/apex/up-examples/blob/master/oss/golang-uploads/main.go
// - https://github.com/choonkeat/attache/blob/master/lib/attache/upload.rb

package main

import (
  // "encoding/base64"
  // "io/ioutil"
  "net/http"
  "os"
  "fmt"
  "log"

  // "github.com/apex/log"
  // "github.com/apex/log/handlers/json"
  // "github.com/apex/log/handlers/text"
)

// Use JSON logging when run by Up (including `up start`).
// func init() {
//   if os.Getenv("UP_STAGE") == "" {
//     log.SetHandler(text.Default)
//   } else {
//     log.SetHandler(json.Default)
//   }
// }

func main() {
  http.HandleFunc("/upload", upload)

  if err := http.ListenAndServe(":" + os.Getenv("PORT"), nil); err != nil {
    log.Fatalf("Error Listening: %s", err)
  }
}

func upload(w http.ResponseWriter, r *http.Request) {
  log.Printf("%s", r)

  switch r.Method {
  case "POST", "PUT", "PATCH":

    r.ParseMultipartForm(32 << 20)
    file, handler, err := r.FormFile("file")
    if err != nil {
      // log.WithError(err).Error("Parsing Form")
      http.Error(w, "Error Parsing Form: %s", http.StatusBadRequest)
      return
    }
    defer file.Close()

    fmt.Printf("%s", handler.Filename)
    fmt.Printf("%s", handler.Header.Get("Content-Type"))

    // bytes, err := ioutil.ReadAll(file)
    // if err != nil {
    //   // log.WithError(err).Error("Reading File")
    //   http.Error(w, "Error Reading File.", http.StatusInternalServerError)
    //   return
    // }

  case "OPTIONS":
    w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, OPTIONS")

  default:
    http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
  }
}
