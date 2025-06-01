.PHONY: check test clean

test: check

check:
	@go test ./... -benchmem -race

clean:
	@rm -f *.test */*.test

GOLANGCI_LINT_VERSION = "v2.1.6"
lint:
#	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run
	docker run --rm -t -v $$(pwd):/app -w /app \
--user $$(id -u):$$(id -g) \
-v $$(go env GOCACHE):/.cache/go-build -e GOCACHE=/.cache/go-build \
-v $$(go env GOMODCACHE):/.cache/mod -e GOMODCACHE=/.cache/mod \
-v ~/.cache/golangci-lint:/.cache/golangci-lint -e GOLANGCI_LINT_CACHE=/.cache/golangci-lint \
golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) golangci-lint run
