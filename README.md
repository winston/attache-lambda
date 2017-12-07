```
make
```

will run tests, and if it passes, compile the code to generate a `build/upload-lambda` executable

NOTE: to cross compile for other architecture, setup the correct env

```
make clean
GOOS=linux GOARCH=amd64 make
```
