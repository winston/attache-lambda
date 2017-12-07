build/upload-lambda: test
	go build -o build/upload-lambda .

test:
	go test ./...

clean:
	rm -f build/upload-lambda
