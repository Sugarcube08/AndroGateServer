project_name = sms-gateway
image_name = androgate-server:latest

extension=
ifeq ($(OS),Windows_NT)
	extension = .exe
endif

.PHONY: \
	all fmt lint test coverage benchmark deps gen clean help \
	init init-dev air db-upgrade db-upgrade-raw run build install \
	docker-build docker docker-clean

all: gen fmt lint coverage ## Run all tests and checks

fmt: ## Format the code
	golangci-lint fmt

lint: ## Lint the code
	golangci-lint run --timeout=5m

test: ## Run tests
	go test -race -shuffle=on -count=1 -covermode=atomic -coverpkg=./... -coverprofile=coverage.out ./...

coverage: test ## Generate coverage
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

benchmark: ## Run benchmarks
	go test -run=^$$ -bench=. -benchmem ./... | tee benchmark.txt

deps: ## Install dependencies
	go mod download

gen: ## Generate code
	PATH="$(shell go env GOPATH)/bin:$$PATH" go generate ./...

clean: ## Remove build artifacts
	rm -f coverage.* benchmark.txt
	rm -rf dist tmp

help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

###

init: deps

init-dev: init
	go install github.com/air-verse/air@latest \
		&& go install github.com/swaggo/swag/cmd/swag@latest \
		&& go install github.com/pressly/goose/v3/cmd/goose@latest

ngrok:
	ngrok http 3000

air: ## Run development server
	@command -v air >/dev/null 2>&1 || { \
      echo "Please install air: go install github.com/air-verse/air@latest"; \
      exit 1; \
    }
	@echo "Starting development server with air..."
	TZ=UTC DEBUG=1 air

db-upgrade:
	go run ./cmd/$(project_name)/main.go db:migrate

db-upgrade-raw:
	go run ./cmd/$(project_name)/main.go db:auto-migrate
	
run:
	go run cmd/$(project_name)/main.go

test-e2e: test
build:
	go build -o tmp/$(project_name) ./cmd/$(project_name)
	
install:
	go install ./cmd/$(project_name)

docker-build:
	docker build -t $(image_name) .

docker:
	docker compose -f deployments/docker-compose/docker-compose.yml up --build

docker-clean:
	docker compose -f deployments/docker-compose/docker-compose.yml down --volumes
