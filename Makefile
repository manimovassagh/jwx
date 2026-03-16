VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    ?= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS  = -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

.PHONY: build test lint fmt check install clean setup help

build:
	go build -ldflags "$(LDFLAGS)" -o bin/jwx ./cmd/jwx

test:
	go test ./...

fmt:
	gofmt -w .

lint:
	golangci-lint run

check: fmt lint test

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/jwx

clean:
	rm -rf bin/

setup:
	cp hooks/pre-commit .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
	@echo "✓ Pre-commit hook installed"

help:
	@echo "Available targets:"
	@echo "  build    - Build the jwx binary"
	@echo "  test     - Run tests"
	@echo "  lint     - Run golangci-lint"
	@echo "  fmt      - Format code with gofmt"
	@echo "  check    - Run fmt + lint + test"
	@echo "  install  - Install jwx to GOPATH"
	@echo "  clean    - Remove build artifacts"
	@echo "  setup    - Install pre-commit hooks"
	@echo "  help     - Show this help"
