.PHONY: build clean deploy gomodgen

build: test
	mkdir -p bin
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin ./...

test: 
	go test ./...

clean:
	rm -rf ./bin ./vendor

deploy: clean build
	sls deploy --verbose

