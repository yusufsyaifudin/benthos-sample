
# compress golang binary size: https://blog.filippo.io/shrink-your-go-binaries-with-this-one-weird-trick/
build: install
	ls -al
	rm -rf artifacts
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags="-s -w" -mod=mod -a -installsuffix cgo -o artifacts/raftsample-368 main.go
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -mod=mod -a -installsuffix cgo -o artifacts/raftsample main.go
	GO111MODULE=on CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -mod=mod -a -installsuffix cgo -o artifacts/raftsample-intel.osx main.go
	GO111MODULE=on CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -mod=mod -a -installsuffix cgo -o artifacts/raftsample-arm64.osx main.go
	GO111MODULE=on CGO_ENABLED=0 GOOS=windows go build -ldflags="-s -w" -mod=mod -a -installsuffix cgo -o artifacts/raftsample.exe main.go

install:
	go mod tidy
