.PHONY: all delegate test mocks clean clearlocal app

all: test delegate app

delegate:
	cd dphoto && go build

test:
	docker-compose up -d
	go test ./... -race -cover

mocks:
	mockery --all -r

clean:
	go clean -testcache

clearlocal:
	aws --endpoint "http://localhost:4566" --region eu-west-1 s3 rm --recursive "s3://dphoto-local"
	aws --endpoint http://localhost:4566 --region eu-west-1 dynamodb delete-table --table dphoto-local
	rm -rf .build/volumes/*

app:
	cd app && make 
