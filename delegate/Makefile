all: test delegate

delegate:
	go build

test:
	go test ./... -race -cover

mocks:
	mockery --all -r --testonly --inpackage