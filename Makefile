.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/connectHandler connectHandler/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/disconnectHandler disconnectHandler/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/messageHandler messageHandler/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
