.PHONY: build clean deploy zip deploy-sls

build:
	env CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -tags lambda.norpc -o build/auth/bootstrap cmd/auth/main.go
	env CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -tags lambda.norpc -o build/receiver/bootstrap cmd/receiver/main.go
	env CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -tags lambda.norpc -o build/writer/bootstrap cmd/writer/main.go

zip:
	zip -j build/auth.zip build/auth/bootstrap
	zip -j build/receiver.zip build/receiver/bootstrap
	zip -j build/writer.zip build/writer/bootstrap

clean:
	rm -rf ./build

deploy-sls:
	sls deploy --verbose

deploy: clean build zip deploy-sls
