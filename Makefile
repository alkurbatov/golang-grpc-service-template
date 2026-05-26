MAKEFLAGS := --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c

SERVICE_NAME = templatesrv
BUILD_DIR = dist
SEED := on

PROTOC = protoc --proto_path=./proto/ --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative
COMPOSE = docker compose -f deployments/docker-compose.yaml -f deployments/docker-compose.dev.yaml
COMPOSE_DEBUG = docker compose -f deployments/docker-compose.yaml -f deployments/docker-compose.debug.yaml

# NB (alkurbatov): Although this template has small coverage threshold
# production-ready services must have >= 80% coverage.
TEST_COVERAGE_THRESHOLD = 70

.DEFAULT_GOAL := help
.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-38s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: download
download: ## Download go.mod dependencies
	@echo Downloading go.mod dependencies
	go mod download -x

.PHONY: update
update: ## Update all Golang modules at once
	go get -u ./...
	go mod tidy

.PHONY: install-tools
install-tools: download ## Install dev tools
# NB (alkurbatov): Add dev packages (e.g. protoc) you want to install here.
#
# (!) Do not forget to specify version, commit or tag.
	@echo Installing dev tools
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

.PHONY: proto
proto: ## Generate gRPC protobuf bindings
	$(PROTOC) \
		--go_out=./pkg/echopb/v1 \
		--go-grpc_out=./pkg/echopb/v1 \
		./proto/echo.proto

.PHONY: build
build:
	./scripts/build cmd/$(SERVICE_NAME) $(BUILD_DIR)/$(SERVICE_NAME)

.PHONY: run
run: stop ## Run the project in docker compose
	$(COMPOSE) up -d --build
	$(COMPOSE) logs -f

.PHONY: debug ## Run the project in docker compose with remote debugger
debug:
	$(COMPOSE_DEBUG) up -d --build
	@echo "Now run 'dlv connect :2345' to attach debugger to the service"
	$(COMPOSE_DEBUG) logs -f

.PHONY: stop
stop: ## Stop the running project and destroy containers
	$(COMPOSE) down

.PHONY: clean
clean: stop
	rm -rf $(BUILD_FOLDER)

.PHONY: lint-golang
lint-golang: ## Lint Golang source code
	golangci-lint run
	go tool deadcode -test ./... | tee deadcode.out && [ ! -s deadcode.out ]

.PHONY: lint-shell
lint-shell: ## Lint shell scripts
	shellcheck --severity=warning ./scripts/*

.PHONY: lint-docker
lint-docker: ## Lint Dockerfile manifests
	hadolint -c .hadolint.yaml Dockerfile*

.PHONY: lint
lint: lint-golang lint-shell lint-docker ## Lint project source

.PHONY: fmt
fmt: ## Format the source code
	go run mvdan.cc/gofumpt@latest -l -w -extra .
	go run golang.org/x/tools/cmd/goimports@latest -l -w .
	go run github.com/daixiang0/gci@latest write \
		--skip-generated \
		--custom-order \
		-s standard \
		-s default \
		-s prefix\(github.com/alkurbatov/golang-grpc-service-template\) \
		-s blank \
		-s dot \
		.

.PHONY: unit-tests
unit-tests: ## Run unit tests
	go test -v -race -shuffle=$(SEED) ./{internal,pkg}/... -coverprofile=coverage.out -covermode atomic
	@grep -v -E "(_mock|.pb).go" coverage.out > coverage.out.tmp
	@mv coverage.out.tmp coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out
	@./scripts/check-coverage coverage.out $(TEST_COVERAGE_THRESHOLD)

.PHONY: update-snapshots
update-snapshots: ## Update snapshots used in unit tests
	@UPDATE_SNAPS=true go test -v ./{internal,pkg}/...

.PHONY: smoke-tests
smoke-tests: ## Run smoke tests
	./scripts/run-smoke-tests

.PHONY: docs
docs: ## View project documentation
	@echo "Project and packages documentation available at:"
	@echo -e "\thttp://127.0.0.1:3000/pkg/"
	@go tool godoc -http=:3000 -index > /dev/null
