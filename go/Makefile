.PHONY: checks
checks:
	@# 1. Checking code formatting
	go fmt ./...
	@# 2. Running linter
	golangci-lint run ./...
	@# 3. Running tests
	go test -v ./...
	@# 4. Running build
	go build ./...
	@# 5. Running vet
	go vet ./...
	@# 6. Running tidy
	@echo "Checking if go.mod is tidy..."
	@cp go.mod go.mod.tmp
	@cp go.sum go.sum.tmp
	@go mod tidy
	@if ! cmp -s go.mod go.mod.tmp || ! cmp -s go.sum go.sum.tmp; then \
		echo "Error: go.mod or go.sum were modified by 'go mod tidy'."; \
		mv go.mod.tmp go.mod; \
		mv go.sum.tmp go.sum; \
		exit 1; \
	fi
	@rm go.mod.tmp go.sum.tmp
