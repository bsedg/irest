PACKAGES=$(shell go list ./... | grep -v /vendor/)
COVERPROFILE=cover.out

all: clean build fmt lint test vet

build:
	@echo "+ $@"
	@go build .

build_example:
	@echo "+ $@"
	@go build cmd/example_report/example_report.go

clean:
	@echo "+ $@"
	@rm -f example_report $(COVERPROFILE)

cover:
	@echo "+ $@"
	@go test -coverprofile=$(COVERPROFILE) .
	@go tool cover -html=$(COVERPROFILE)

fmt:
	@echo "+ $@"
	@gofmt -s -l . | grep -v vendor | tee /dev/stderr

lint:
	@echo "+ $@"
	@golint ./... | grep -v vendor | tee /dev/stderr

test: fmt lint vet
	@echo "+ $@"
	@go test -v $(shell go list ./... | grep -v vendor)

vet:
	@echo "+ $@"
	@go vet $(shell go list ./... | grep -v vendor)
