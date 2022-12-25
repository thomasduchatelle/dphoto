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

setup-go:
	docker-compose pull
	docker-compose up -d

test-go:
	AWS_PROFILE="" go test ./... -race -cover

build-go:
	go build ./cmd/...

install-cli:
	go install ./cmd/...

#######################################
## WEB
#######################################

.PHONY: clean-web setup-web test-web build-web

clean-web:
	cd web && yarn clean

setup-web:
	cd web && npm install
	cd web && yarn install

test-web:
	+echo "No tests implemented for WEB"

build-web:
	cd web && CI=true yarn build

start:
	cd web && DANGEROUSLY_DISABLE_HOST_CHECK=true yarn start


#######################################
## API
#######################################

.PHONY: clean-api test-api build-api

clean-api:
	cd api && rm -rf ./bin ./vendor

test-api:
	cd api && AWS_PROFILE="" go test ./...

build-api:
	cd api/lambdas && \
		mkdir -p ../../bin && \
		export GO111MODULE=on && \
		env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ../../bin ./...

#######################################
## APP = WEB + API
#######################################

.PHONY: setup-app test-app build-app deploy-app

setup-app: setup-web
	cd deployments/sls && npm install

test-app: test-api test-web

build-app: build-api build-web

AWS_PROFILE ?= dphoto
deploy-app: clean-web clean-api test-app build-app
	export AWS_PROFILE="$(AWS_PROFILE)" && cd deployments/sls && sls deploy --debug

#######################################
## UTILS
#######################################

.PHONY: mocks clearlocal dcdown dcup

mocks:
	mockery --all -r

clearlocal: dcdown dcup

dcdown:
	AWS_ACCESS_KEY_ID="localstack" AWS_SECRET_ACCESS_KEY="localstack" aws --endpoint "http://localhost:4566" --region eu-west-1 s3 rm --recursive "s3://dphoto-local" | cat || echo "skipping"
	AWS_ACCESS_KEY_ID="localstack" AWS_SECRET_ACCESS_KEY="localstack" aws --endpoint "http://localhost:4566" --region eu-west-1 dynamodb delete-table --table dphoto-local | cat || echo "skipping"
	docker-compose down -v || echo "skipping"
	rm -rf .build/localstack

dcup:
	docker-compose up -d
