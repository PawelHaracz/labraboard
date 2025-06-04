SHELL := /bin/bash

VERSION ?= v0.0.1
REGISTRY ?= ghcr.io
IMAGE_BUILDER ?= docker
IMAGE_BUILD_CMD ?= build
IMAGE_NAME ?= pawelharacz/labraboard

export IMG = $(REGISTRY)/$(IMAGE_NAME):$(VERSION)

.PHONY: docker-build clean clean-all install install-dependencies install-frontend-dependencies fmt lint vet doc update-dependencies test-unit test-cover test build-api build-handlers build-frontend build docker-push docker-compose-build docker-compose-up docker-compose-stop build-swagger helm-render helm-push security-scan dependency-check release-prepare release-publish

install-dependencies: ## Install dependencies
	go install -v github.com/go-delve/delve/cmd/dlv@latest
	go mod download
	@echo "Go dependencies installed"

install-frontend-dependencies: ## Install frontend dependencies
	cd client && yarn install
	@echo "Frontend dependencies installed"
		
install: install-dependencies install-frontend-dependencies ## Install application and dependencies
	go install -v ./...
	@echo "Application installed"

# Code Quality
fmt:    ## format the go source files
	go fmt ./...

lint:   ## run go lint on the source files
	go tool golangci-lint run

vet:    ## run go vet on the source files
	go vet ./...

# Security
security-scan: ## Run security scans
	go tool gosec ./...
	go tool trivy image $(IMG)

# Dependency Management
dependency-check: ## Check for outdated dependencies
	go list -u -m all

update-dependencies:    ## update golang dependencies
	go get -u ./...
	go mod tidy

# Testing
test-unit: ## Run unit tests
	go test ./...

test-cover:     ## Run test coverage and generate html report
	rm -fr coverage
	mkdir coverage
	go list -f '{{if gt (len .TestGoFiles) 0}}"go test -covermode count -coverprofile {{.Name}}.coverprofile -coverpkg ./... {{.ImportPath}}"{{end}}' ./... | xargs -I {} bash -c {}
	echo "mode: count" > coverage/cover.out
	grep -h -v "^mode:" *.coverprofile >> "coverage/cover.out"
	rm *.coverprofile
	go tool cover -html=coverage/cover.out -o=coverage/cover.html

test: test-unit test-cover ## Run all tests

# Documentation
doc:    ## generate godocs and start a local documentation webserver on port 8085
	go tool godoc -http=:8085
	@echo "Godoc documentation generated"

doc-local: ## generate godocs for local project only
	@echo "Starting godoc server for labraboard project..."
	@echo "Navigate to http://localhost:8085/pkg/labraboard/ to see your API docs"
	go tool godoc -http=:8085

build-swagger: ## Generate Swagger documentation
	go tool swag init -g ./cmd/api/main.go -o ./docs
	@echo "Swagger documentation generated"

# Building
build-frontend: install-frontend-dependencies ## Build frontend application
	cd client && yarn build
	@echo "Frontend built"

build-api:  ## Build API server
	cd cmd/api && go build -o ../../bin/api
	@echo "API server built"

build-handlers:  ## Build handlers
	cd cmd/handlers && go build -o ../../bin/handlers
	@echo "Handlers built"

build: build-frontend build-api build-handlers ## Build all components
	@echo "All components built"

# Docker Operations
docker-build: ## Build Docker image
	CGO_ENABLED=0 GOOS=linux $(IMAGE_BUILDER) $(IMAGE_BUILD_CMD) -t $(IMG) .
	@echo "Docker image built"

docker-push: ## Push Docker image to registry
	$(IMAGE_BUILDER) push $(IMG)
	@echo "Docker image pushed"

docker-compose-build: ## Build Docker Compose services
	docker compose build
	@echo "Docker Compose services built"

docker-compose-up: ## Start Docker Compose services
	docker compose up
	@echo "Docker Compose services started"

docker-compose-stop: ## Stop Docker Compose services
	docker compose stop
	@echo "Docker Compose services stopped"

# Helm Operations
helm-render: ## Render Helm templates
	helm template charts/labraboard --set image.repository=$(REGISTRY)/$(IMAGE_NAME) --set image.tag=$(VERSION)
	@echo "Helm templates rendered"
helm-push: ## Package and push Helm chart
	helm package charts/labraboard --app-version $(VERSION) --version $(VERSION)
	helm push labraboard-*.tgz oci://ghcr.io/$(IMAGE_NAME)
	@echo "Helm chart pushed"
# Release Management
release-prepare: ## Prepare release artifacts
	@echo "Preparing release $(VERSION)"
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)
	@echo "Release prepared"	
release-publish: ## Publish release
	@echo "Publishing release $(VERSION)"
	@gh release create $(VERSION) --generate-notes
	@echo "Release published"
# Maintenance
clean: ## Clean build artifacts
	rm -rf bin/*
	go clean -i ./...
	@echo "Build artifacts cleaned"
clean-all: clean ## Remove all generated artifacts
	rm -rf coverage/
	rm -rf dist/
	rm -f *.tgz
	@echo "All generated artifacts cleaned"
# Help
help: ## Display this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk -F ':|##' '/^[^\t].+?:.*?##/ { printf "  %-20s %s\n", $$1, $$NF }' $(MAKEFILE_LIST)
