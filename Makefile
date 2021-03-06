SOURCES := $(shell find . -name '*.go')
BINARY := scanner-trivy
IMAGE_TAG := dev
IMAGE := aquasec/harbor-scanner-trivy:$(IMAGE_TAG)

build: $(BINARY)

test: build
	GO111MODULE=on go test -v -short -race -coverprofile=coverage.txt -covermode=atomic ./...

test-integration: build
	GO111MODULE=on go test -count=1 -v -tags=integration ./test/integration/...

test-component: container
	GO111MODULE=on go test -count=1 -v -tags=component ./test/component/...

$(BINARY): $(SOURCES)
	GOOS=linux GO111MODULE=on CGO_ENABLED=0 go build -o $(BINARY) cmd/scanner-trivy/main.go

container: build
	docker build --no-cache -t $(IMAGE) .

lint:
	./bin/golangci-lint --build-tags component,integration run -v

setup:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.21.0

PHONY: setup
