.PHONY: all clean setup test build deploy test-go

all: clean test build

clean: clean-web clean-api
	go clean -testcache

setup: setup-infra-data setup-app

test: test-infra-data test-go test-web

build: build-go build-app

deploy: deploy-infra-data deploy-app install-cli

#######################################
## INFRA DATA
#######################################

.PHONY: setup-infra-data test-infra-data deploy-infra-data
INFRA_DATA := deployments/infra-data

setup-infra-data:
	command -v tfenv > /dev/null \
		cd $(INFRA_DATA) && \
		tfenv install && tfenv use
	cd $(INFRA_DATA) && \
		terraform init

test-infra-data:
	cd $(INFRA_DATA) && \
		terraform validate

deploy-infra-data:
	cd $(INFRA_DATA) && \
		terraform apply

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

test-go:
	AWS_PROFILE="" go test ./... -race -cover

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

clean-web:
	cd web && yarn clean

setup-web:
	cd web && yarn

test-web:
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
## API
#######################################

.PHONY: clean-api test-api build-api

clean-api:
	rm -rf ./bin ./api/vendor

test-api: test-go
	cd api/lambdas && AWS_PROFILE="" go test ./...

build-api:
	cd api/lambdas && \
		mkdir -p ../../bin && \
		export GO111MODULE=on && \
		env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w $(call unquote,$(BUILD_LD_FLAGS))" -o ../../bin ./...

#######################################
## APP = WEB + API
#######################################

.PHONY: setup-app test-app build-app deploy-app bg down

setup-app: setup-web
	cd deployments/sls && yarn

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
	mockery --all --dir pkg -r --output internal/mocks
	mockery --all --dir cmd -r --output internal/mocks
	git add internal/mocks

clearlocal: dcdown dcup

dcdown:
	AWS_ACCESS_KEY_ID="localstack" AWS_SECRET_ACCESS_KEY="localstack" aws --endpoint "http://localhost:4566" --region eu-west-1 s3 rm --recursive "s3://dphoto-local" | cat || echo "skipping"
	AWS_ACCESS_KEY_ID="localstack" AWS_SECRET_ACCESS_KEY="localstack" aws --endpoint "http://localhost:4566" --region eu-west-1 dynamodb delete-table --table dphoto-local | cat || echo "skipping"
	docker-compose down -v || echo "skipping"
	rm -rf .build/localstack

dcup:
	docker-compose up -d
	AWS_ACCESS_KEY_ID=localstack AWS_SECRET_ACCESS_KEY=localstack aws --endpoint http://localhost:4566 --region eu-west-1 s3 mb s3://dphoto-local | cat
	AWS_ACCESS_KEY_ID=localstack AWS_SECRET_ACCESS_KEY=localstack aws --endpoint http://localhost:4566 --region eu-west-1 sns create-topic --name dphoto-local-archive-jobs | cat
