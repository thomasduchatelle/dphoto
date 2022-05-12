.PHONY: all clean install test build deploy

all: test build

clean: clean-app
	go clean -testcache

install: install-infra-data install-domain install-app

test: test-infra-data test-domain test-cli test-app

build: build-cli build-app

deploy: deploy-infra-data deploy-app deploy-cli-local

#######################################
## INFRA DATA
#######################################

.PHONY: install-infra-data test-infra-data deploy-infra-data

install-infra-data:
	command -v tfenv > /dev/null \
		cd infra-data && \
		tfenv install && tfenv use
	cd infra-data && \
		terraform init

test-infra-data:
	cd infra-data && \
		terraform validate

deploy-infra-data:
	cd infra-data && \
		terraform apply

#######################################
## DOMAIN
#######################################

.PHONY: install-domain test-domain

install-domain:
	docker-compose pull
	docker-compose up -d

test-domain:
	AWS_PROFILE="" go test ./domain/... -race -cover

#######################################
## APP
#######################################

.PHONY: clean-app install-app test-app test-app-api test-app-ui build-app build-app-api build-app-ui deploy-app

clean-app:
	cd app && rm -rf ./bin ./vendor
	cd app/viewer_ui && yarn clean

install-app:
	cd app && npm install
	cd app/viewer_ui && yarn install

test-app: test-app-api test-app-ui

test-app-api: test-domain
	AWS_PROFILE="" go test ./app/viewer_api/...

test-app-ui:
	+echo "Implement tests for APP VIEWER UI"

build-app: build-app-api build-app-ui

build-app-api:
	cd app && \
		mkdir -p bin && \
		export GO111MODULE=on && \
		env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin ./...

build-app-ui:
	cd app/viewer_ui && CI=true yarn build

AWS_PROFILE ?= dphoto
deploy-app: clean-app test-app build-app
	export AWS_PROFILE="$(AWS_PROFILE)" && cd app && sls deploy --debug

start:
	cd app/viewer_ui && DANGEROUSLY_DISABLE_HOST_CHECK=true yarn start

#######################################
## CLI
#######################################

.PHONY: test-cli build-cli deploy-cli-local

test-cli: test-domain
	AWS_PROFILE="" go test ./dphoto/... -race -cover

build-cli:
	cd dphoto && go build

deploy-cli-local:
	cd dphoto && go install

#######################################
## UTILS
#######################################

.PHONY: mocks clearlocal

mocks:
	mockery --all -r

clearlocal:
	aws --endpoint "http://localhost:4566" --region eu-west-1 s3 rm --recursive "s3://dphoto-local"
	aws --endpoint http://localhost:4566 --region eu-west-1 dynamodb delete-table --table dphoto-local
	rm -rf .build/volumes/*
