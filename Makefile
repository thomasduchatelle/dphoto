.PHONY: all clean setup test build deploy test-go

all: clean test build

clean: clean-web clean-api
	go clean -testcache

setup: setup-cdk setup-app

test: test-cdk test-go test-web

build: build-go build-app

deploy: deploy-cdk deploy-app install-cli

#######################################
## CDK
#######################################

.PHONY: setup-cdk test-cdk deploy-cdk
CDK_DIR := deployments/cdk

setup-cdk:
	command -v npm > /dev/null || (echo "npm is required" && exit 1)
	cd $(CDK_DIR) && npm install
	command -v cdk > /dev/null || npm install -g aws-cdk

test-cdk:
	cd $(CDK_DIR) && npm test

deploy-cdk:
	cd $(CDK_DIR) && cdk deploy --context environment=next

#######################################
## PKG & CLI
#######################################

.PHONY: setup-go test-go build-go install-go
unquote = $(patsubst "%,%,$(patsubst %",%,$(1)))

APPLICATION_VERSION ?= ""
APPLICATION_VERSION_SNAPSHOT ?= "true"
BUILD_LD_FLAGS = "-X 'github.com/thomasduchatelle/dphoto/pkg/meta.SemVer=$(call unquote,$(APPLICATION_VERSION))' -X 'github.com/thomasduchatelle/dphoto/pkg/meta.Snapshot=$(call unquote,$(APPLICATION_VERSION_SNAPSHOT))'"

setup-go:
	docker-compose pull
	docker-compose up -d

test-pkg:
	AWS_PROFILE="" go test ./... -race -cover

test-go: test-pkg test-api

build-go:
	go build -ldflags="-s -w $(call unquote,$(BUILD_LD_FLAGS))"  -o ./ ./cmd/...

build-cli:
	env GOARCH=amd64 GOOS=linux  CGO_ENABLED=0 go build -ldflags="-s -w $(call unquote,$(BUILD_LD_FLAGS))" -o ./bin-cli/dphoto-amd64-linux  ./cmd/dphoto
	env GOARCH=amd64 GOOS=darwin CGO_ENABLED=0 go build -ldflags="-s -w $(call unquote,$(BUILD_LD_FLAGS))" -o ./bin-cli/dphoto-amd64-darwin ./cmd/dphoto
	env GOARCH=arm64 GOOS=darwin CGO_ENABLED=0 go build -ldflags="-s -w $(call unquote,$(BUILD_LD_FLAGS))" -o ./bin-cli/dphoto-arm64-darwin ./cmd/dphoto

install-cli:
	go install ./cmd/...

#######################################
## WEB
#######################################

.PHONY: clean-web setup-web test-web build-web update-snapshots

clean-web: clean-waku
	cd web && yarn clean

setup-web: setup-waku
	cd web && yarn

test-web: test-waku
	cd web && yarn test:ci

update-snapshots:
	@echo "Update snapshots [should only be used on CI]"
	rm -rf web/src/stories/__image_snapshots__ && cd web && CI=true yarn test:ci -u

build-web:
	cd web && CI=true yarn build

start:
	docker-compose up -d wiremock && \
		cd web && DANGEROUSLY_DISABLE_HOST_CHECK=true yarn start

storybook:
	cd web && yarn storybook

test-web-ci:
	docker build -t dphoto-puppeteer ./tools/puppeteer/
	docker run --rm -v "$(shell pwd):/app" -it dphoto-puppeteer yarn test:ci

#######################################
## WAKU
#######################################

.PHONY: clean-waku setup-waku test-waku build-waku

clean-waku:
	cd web-waku && rm -rf dist/

setup-waku:
	cd web-waku && npm install

test-waku:
	@echo "Waku tests - placeholder (no tests configured yet)"
	cd web-waku && npm run build

build-waku:
	cd web-waku && npm run build

start-waku:
	cd web-waku && npm run dev

#######################################
## API
#######################################

.PHONY: clean-api test-api build-api

clean-api:
	rm -rf ./bin ./api/vendor

test-api:
	cd api/lambdas && AWS_PROFILE="" go test ./...

build-api:
	cd api/lambdas && \
		rm -rf ../../bin/* && \
		mkdir -p ../../bin && \
		export GO111MODULE=on && \
		env GOARCH=arm64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w $(call unquote,$(BUILD_LD_FLAGS))" -o ../../bin ./...

	cd bin/ && \
		for bin in * ; do mv "$$bin" bootstrap && zip "$$bin.zip" bootstrap ; done

#######################################
## APP = WEB + API
#######################################

.PHONY: setup-app test-app build-app deploy-app bg down

setup-app: setup-web

clean-app: clean-api clean-web

test-app: test-api test-web

build-app: build-api build-web

AWS_PROFILE ?= dphoto
deploy-app: clean-web clean-api build-app
	export AWS_PROFILE="$(AWS_PROFILE)" && cd deployments/sls && sls deploy

bg:
	docker-compose --profile bg up -d

down:
	docker-compose --profile bg down

#######################################
## UTILS
#######################################

.PHONY: mocks clearlocal dcdown dcup

mocks:
	rm -rf internal/mocks
	mockery --all --dir pkg -r --with-expecter --output internal/mocks
	mockery --all --dir cmd -r --with-expecter --output internal/mocks
	git add internal/mocks

clearlocal:
	AWS_ACCESS_KEY_ID="localstack" AWS_SECRET_ACCESS_KEY="localstack" aws --endpoint "http://localhost:4566" --region us-east-1 s3 rm --recursive "s3://dphoto-local" | cat || echo "skipping"
	AWS_ACCESS_KEY_ID="localstack" AWS_SECRET_ACCESS_KEY="localstack" aws --endpoint "http://localhost:4566" --region us-east-1 dynamodb delete-table --table dphoto-local | cat || echo "skipping"

dcdown:
	AWS_ACCESS_KEY_ID="localstack" AWS_SECRET_ACCESS_KEY="localstack" aws --endpoint "http://localhost:4566" --region us-east-1 s3 rm --recursive "s3://dphoto-local" | cat || echo "skipping"
	AWS_ACCESS_KEY_ID="localstack" AWS_SECRET_ACCESS_KEY="localstack" aws --endpoint "http://localhost:4566" --region us-east-1 dynamodb delete-table --table dphoto-local | cat || echo "skipping"
	docker-compose down -v || echo "skipping"
	rm -rf .build/localstack

dcup:
	docker-compose up -d
	AWS_ACCESS_KEY_ID=localstack AWS_SECRET_ACCESS_KEY=localstack aws --endpoint http://localhost:4566 --region us-east-1 s3 mb s3://dphoto-local | cat
	AWS_ACCESS_KEY_ID=localstack AWS_SECRET_ACCESS_KEY=localstack aws --endpoint http://localhost:4566 --region us-east-1 sns create-topic --name dphoto-local-archive-jobs | cat
