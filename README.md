```
make
```

will run tests, and if it passes, compile the code to generate a `build/upload-lambda` executable

NOTE: to cross compile for other architecture, setup the correct env

```
make clean
GOOS=linux GOARCH=amd64 make
```

## Usage

``` go
package main

import (
	"log"
	"net/http"
	"os"

	attache "github.com/winston/attache-lambda"
	"github.com/winston/attache-lambda/s3store"
)

func main() {
	http.Handle("/", attache.Server{
		Storage: s3store.Store{
			Bucket: os.Getenv("AWS_BUCKET"),
		},
	})

	log.Printf("Listening to %s...", os.Getenv("PORT"))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
```
