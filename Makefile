.PHONY: build clean deploy

build:
	env CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/receiver cmd/receiver/main.go
	env CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/writer cmd/writer/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
