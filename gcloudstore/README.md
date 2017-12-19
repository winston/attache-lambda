Checkout https://github.com/GoogleCloudPlatform/cloud-functions-go

Change the `main` function in `main.go` to become

``` go
func main() {
	flag.Parse()

	http.Handle(nodego.HTTPTrigger, attache.Server{
		Storage: gcloudstore.NewStore("your-bucket-name"),
		GetPrefixPath: "/execute?",
	})

	nodego.TakeOver()
}
```

NOTE:
- Google Cloud Function http endpoint matches request path strictly, e.g. `/attache` works but `/attache/thing.jpg` is 404.
- To workaround that, we need to specify the url as `/attache?thing.jpg` instead
- But internally our http server sees the prefix as `/execute` ¯\_(ツ)_/¯, so we have to configure `GetPrefixPath: "/execute?"`
