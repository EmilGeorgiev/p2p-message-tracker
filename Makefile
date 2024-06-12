# Define the tools as variables
GOFMT := gofmt -s -w
GOIMPORTS := goimports -w
GOLINT := golint ./...
GOVET := go vet ./...
STATICCHECK := staticcheck ./...
GOLANGCI_LINT := golangci-lint run
ERRCHECK := errcheck ./...
GOSEC := gosec ./...

.PHONY: all fmt lint vet staticcheck golangci-lint errcheck gosec check

# Format the code
fmt:
	@echo "Running gofmt..."
	$(GOFMT) .
	@echo "Running goimports..."
	$(GOIMPORTS) .

# Lint the code
lint:
	@echo "Running golint..."
	$(GOLINT)

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET)

# Run staticcheck
staticcheck:
	@echo "Running staticcheck..."
	$(STATICCHECK)

# Run golangci-lint
golangci-lint:
	@echo "Running golangci-lint..."
	$(GOLANGCI_LINT)

# Run errcheck
errcheck:
	@echo "Running errcheck..."
	$(ERRCHECK)

# Run gosec
gosec:
	@echo "Running gosec..."
	$(GOSEC)

# Run all checks
all-checks: fmt lint vet staticcheck golangci-lint errcheck gosec
	@echo "All checks passed!"

# Run tests
test:
	@echo "Running tests...."
	go test ./...