.PHONY: all test

default: test

all: test bench

bench:
	go test ./... -bench=. -benchmem

cyclo:
	gocyclo -over 13 ./*/*.go

fmt:
	go fmt $(shell go list ./...)

gen:
	mockgen -destination mocks/input/mock_reader.go github.com/jedib0t/go-prompter/input Reader
	mockgen -destination mocks/prompt/mock_prompter.go github.com/jedib0t/go-prompter/prompt Prompter

run:
	go run ./examples/prompt/sql

test: gen fmt vet cyclo
	go test -cover -coverprofile=.coverprofile -race $(shell go list ./...)

tools:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@v0.5.1
	go install go.uber.org/mock/mockgen@latest

vet:
	go vet $(shell go list ./...)

