.PHONY: all delegate test mocks clearlocal

all: test delegate

delegate:
	go build

test:
	go test ./... -race -cover

mocks:
	mockery --all -r

clearlocal:
	aws --endpoint "http://localhost:4566" --region eu-west-1 s3 rm --recursive "s3://dphoto-local"
	aws --endpoint http://localhost:4566 --region eu-west-1 dynamodb delete-table --table dphoto-local
	rm -rf .build/volumes/*
