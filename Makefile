SHELL := /bin/bash

VERSION ?= v0.0.1
REGISTRY ?= ghcr.io
IMAGE_BUILDER ?= docker
IMAGE_BUILD_CMD ?= build
IMAGE_NAME ?= pawelharacz/labraboard

export IMG = $(REGISTRY)/$(IMAGE_NAME):$(VERSION)

.PHONY: docker-build clean clean-all install mod fmt lint vet doc update-dependencies test-unit test-cover test build-api build-handlers build docker-push docker-compose-build docker-compose-up docker-compose-stop build-swagger helm-render helm-push

install:    ## build and install go application executable
	go install -v ./...

clean-all:  ## remove all generated artifacts and clean all build artifacts
	go clean -i ./...
	rm -rf bin/*

mod:
	go mod download && go mod verify

fmt:    ## format the go source files
	go fmt ./...

lint:   ## run go lint on the source files
	golangci-lint run

vet:    ## run go vet on the source files
	go vet ./...

doc:    ## generate godocs and start a local documentation webserver on port 8085
	godoc -http=:8085 -index

update-dependencies:    ## update golang dependencies
	dep ensure

test-unit:
	go test ./...

# Generate test coverage
test-cover:     ## Run test coverage and generate html report
	rm -fr coverage
	mkdir coverage
	go list -f '{{if gt (len .TestGoFiles) 0}}"go test -covermode count -coverprofile {{.Name}}.coverprofile -coverpkg ./... {{.ImportPath}}"{{end}}' ./... | xargs -I {} bash -c {}
	echo "mode: count" > coverage/cover.out
	grep -h -v "^mode:" *.coverprofile >> "coverage/cover.out"
	rm *.coverprofile
	go tool cover -html=coverage/cover.out -o=coverage/cover.html

test: test-unit test-cover

build-api:
	cd cmd/api && go build -o ../../bin/api

build-handlers:
	cd cmd/handlers && go build -o ../../bin/handlers

build: mod test build-api build-handlers

docker-build:
	CGO_ENABLED=0 GOOS=linux $(IMAGE_BUILDER) $(IMAGE_BUILD_CMD) -t $(IMG) .

docker-push:
	$(IMAGE_BUILDER) push $(IMG)

docker-compose-build:
	docker compose build

docker-compose-up:
	docker compose up

docker-compose-stop:
	docker compose stop

build-swagger:
	swag init -g ./cmd/api/main.go -o ./docs

helm-render:
	helm template charts/labraboard --set image.repository=$(REGISTRY)/$(IMAGE_NAME) --set image.tag=$(VERSION)

helm-push:
	helm package charts/labraboard --app-version $(VERSION) --version $(VERSION)
	helm push labraboard-*.tgz oci://ghcr.io/$(IMAGE_NAME)
