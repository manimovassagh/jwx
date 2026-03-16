VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    ?= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS  = -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

.PHONY: build test lint fmt check install clean setup

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
