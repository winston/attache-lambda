Checkout https://github.com/GoogleCloudPlatform/cloud-functions-go

Change the `main` function in `main.go` to become

``` go
func main() {
	flag.Parse()

	http.Handle(nodego.HTTPTrigger, attache.Server{
		Storage: gcloudstore.NewStore("your-bucket-name"),
	})

	nodego.TakeOver()
}
```
