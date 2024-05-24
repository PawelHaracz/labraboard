VERSION ?= v0.0.1
REGISTRY ?= ghcr.io
IMAGE_BUILDER ?= docker
IMAGE_BUILD_CMD ?= build
IMAGE_NAME ?= labraboard/labraboard

export IMG = $(REGISTRY)/$(IMAGE_NAME):$(VERSION)
export CGO_ENABLED=0
export GOOS=linux

.PHONY: docker-build

mod:
	go mod download && go mod verify

test:
	  for PACKAGE in $(go list ./...); do go test ${PACKAGE}; done;

build-api:
	cd cmd/api && go build -o ../../bin/api

build-handlers:
	cd cmd/handlers && go build -o ../../bin/handlers

build: mod test build-api build-handlers

clean:
	rm -rf bin/*

docker-build:
	$(IMAGE_BUILDER) $(IMAGE_BUILD_CMD) -t $(IMG) .

docker-push:
	$(IMAGE_BUILDER) push $(IMG)

docker-compose-build:
	docker compose build

docker-compose-up:
	docker compose up

docker-compose-stop:
	docker compose stop