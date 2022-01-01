.PHONY: testApi buildApi testUi buildUi clean build deploy test

testApi:
	AWS_PROFILE="" go test ../domain/...
	AWS_PROFILE="" go test ./viewer_api/...

buildApi:
	mkdir -p bin
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin ./...

testUi:
	echo "TODO - note: 'yarn test' is interactive"

buildUi:
	cd viewer_ui && yarn build

clean:
	rm -rf ./bin ./vendor
	cd viewer_ui && yarn clean

build: buildApi buildUi

deploy: clean test build
	sls deploy --verbose

test: testApi testUi
