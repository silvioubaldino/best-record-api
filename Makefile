export GO111MODULE ?= on
PACKAGES = $(shell go list ./...)
PACKAGES_PATH = $(shell go list -f '{{ .Dir }}' ./...)
# FLAGS = --max-same-issues 0  --max-issues-per-linter 0 --new-from-rev HEAD~1

.PHONY: all
all: check_tools ensure-deps fmt imports linter

.PHONY: check_tools
check_tools:
	@type "golangci-lint" > /dev/null 2>&1 || echo 'Please install golangci-lint: https://golangci-lint.run/usage/install/#local-installation'
	@type "goimports" > /dev/null 2>&1 || echo 'Please install goimports: go get golang.org/x/tools/cmd/goimports'
	@type "govulncheck" > /dev/null 2>&1 || echo 'Please install govulncheck: go install golang.org/x/vuln/cmd/govulncheck@latest'

.PHONY: ensure-deps
ensure-deps:
	@echo "=> Syncing dependencies with go mod tidy"
	@go mod tidy

.PHONY: fmt
fmt:
	@echo "=> Executing gofumpt"
	@gofumpt -l -w .

.PHONY: imports
imports:
	@echo "=> Executing goimports"
	@goimports -w $(PACKAGES_PATH)

# Runs golangci-lint with arguments if provided.
.PHONY: linter
linter:
	@echo "=> Executing golangci-lint$(if $(FLAGS), with flags: $(FLAGS))"
	@golangci-lint run -c .golangci.yml ./... $(FLAGS)

.PHONY: build
build:
	@echo "=> Running build command for windows..."
	cd server && GOOS=windows GOARCH=amd64 go build -o best-record.exe ./main.go