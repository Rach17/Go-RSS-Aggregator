# Go RSS Aggregator Makefile

# Variables
API_BINARY=tmp/api
SCRAPER_BINARY=tmp/scraper
API_CMD=./api
SCRAPER_CMD=./scraper
GO_VERSION=1.21

# Default target
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make check         - Check dependencies and environment"
	@echo "  make build-api     - Build the API server"
	@echo "  make build-scraper - Build the scraper"
	@echo "  make build-all     - Build both API and scraper"
	@echo "  make run-api       - Run the API server"
	@echo "  make run-scraper   - Run the scraper"
	@echo "  make run-both      - Run both API and scraper in background"
	@echo "  make dev           - Run API with hot reload using air"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make test          - Run tests"
	@echo "  make deps          - Download dependencies"
	@echo "  make verify        - Verify dependencies"
	@echo "  make stop          - Stop all running processes"

# Check dependencies and environment
.PHONY: check
check:
	@echo "Checking Go installation..."
	@if ! command -v go >/dev/null 2>&1; then \
		echo "❌ Go is not installed. Please install Go $(GO_VERSION) or later."; \
		exit 1; \
	fi
	@echo "✅ Go is installed: $$(go version)"
	@echo "Checking Go version..."
	@GO_VER=$$(go version | cut -d' ' -f3 | cut -d'o' -f2); \
	if ! printf '%s\n%s\n' "$(GO_VERSION)" "$$GO_VER" | sort -V -C; then \
		echo "⚠️  Go version $$GO_VER found, $(GO_VERSION) or later recommended"; \
	else \
		echo "✅ Go version $$GO_VER is compatible"; \
	fi
	@echo "Checking go.mod file..."
	@if [ ! -f go.mod ]; then \
		echo "❌ go.mod not found. Run 'go mod init' first."; \
		exit 1; \
	fi
	@echo "✅ go.mod found"
	@echo "Checking for required directories..."
	@if [ ! -d "api" ]; then \
		echo "❌ api directory not found"; \
		exit 1; \
	fi
	@if [ ! -d "scraper" ]; then \
		echo "❌ scraper directory not found"; \
		exit 1; \
	fi
	@echo "✅ Required directories found"

# Verify and download dependencies
.PHONY: deps
deps: check
	@echo "Downloading dependencies..."
	go mod download
	@echo "Tidying dependencies..."
	go mod tidy
	@echo "✅ Dependencies updated"

# Verify dependencies
.PHONY: verify
verify: check
	@echo "Verifying dependencies..."
	go mod verify
	@echo "Checking for vulnerabilities..."
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "⚠️  govulncheck not installed. Install with: go install golang.org/x/vuln/cmd/govulncheck@latest"; \
	fi
	@echo "✅ Dependencies verified"

# Create tmp directory if it doesn't exist
tmp:
	@mkdir -p tmp

# Build targets with dependency check
.PHONY: build-api
build-api: check tmp
	@echo "Building API server..."
	@if [ ! -f "api/main.go" ]; then \
		echo "❌ api/main.go not found"; \
		exit 1; \
	fi
	go build -v -o $(API_BINARY) $(API_CMD)
	@if [ ! -f "$(API_BINARY)" ]; then \
		echo "❌ Build failed: $(API_BINARY) not created"; \
		exit 1; \
	fi
	chmod +x $(API_BINARY)
	@echo "✅ API server built successfully"

.PHONY: build-scraper
build-scraper: check tmp
	@echo "Building scraper..."
	@if [ ! -f "scraper/main.go" ]; then \
		echo "❌ scraper/main.go not found"; \
		exit 1; \
	fi
	go build -v -o $(SCRAPER_BINARY) $(SCRAPER_CMD)
	@if [ ! -f "$(SCRAPER_BINARY)" ]; then \
		echo "❌ Build failed: $(SCRAPER_BINARY) not created"; \
		exit 1; \
	fi
	chmod +x $(SCRAPER_BINARY)
	@echo "✅ Scraper built successfully"

.PHONY: build-all
build-all: build-api build-scraper
	@echo "✅ All components built successfully"

# Run targets
.PHONY: run-api
run-api: build-api
	@echo "Starting API server..."
	./$(API_BINARY)

.PHONY: run-scraper
run-scraper: build-scraper
	@echo "Starting scraper..."
	./$(SCRAPER_BINARY)

.PHONY: run-both
run-both: build-all
	@echo "Starting both API and scraper..."
	./$(API_BINARY) & ./$(SCRAPER_BINARY) &
	@echo "✅ Both services started in background"

# Development with hot reload
.PHONY: dev
dev: check
	@echo "Checking if air is installed..."
	@if ! command -v air >/dev/null 2>&1; then \
		echo "❌ Air not found. Installing air..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	@echo "Starting API with hot reload..."
	air

# Utility targets
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf tmp/
	go clean
	@echo "✅ Clean completed"

.PHONY: test
test: check
	@echo "Running tests..."
	go test -v ./...
	@echo "✅ Tests completed"

.PHONY: stop
stop:
	@echo "Stopping all Go processes..."
	pkill -f "$(API_BINARY)" || true
	pkill -f "$(SCRAPER_BINARY)" || true
	pkill -f "go run ./api" || true
	pkill -f "go run ./scraper" || true
	@echo "✅ All processes stopped"

# Install development tools
.PHONY: install-tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "✅ Development tools installed"

# Full setup for new developers
.PHONY: setup
setup: check deps verify install-tools
	@echo "✅ Setup completed successfully!"
	@echo "You can now run:"
	@echo "  make build-all  - To build everything"
	@echo "  make dev        - For development with hot reload"
	@echo "  make run-both   - To run both services"

# Run both services in foreground with visible logs
.PHONY: run-both-logs
run-both-logs: build-all
	@echo "Starting both API and scraper with logs..."
	@echo "Press Ctrl+C to stop both services"
	./$(API_BINARY) & ./$(SCRAPER_BINARY) & wait