# Variables
APP_NAME = Beelzebub
BUILD_DIR = bin
SRC_DIR = cmd
MAIN_FILE = main.go

# Commands
GO = go
RM = rm -rf
MKDIR = mkdir -p

# Flags
GOFLAGS = -mod=vendor
LDFLAGS = -s -w
BUILD_FLAGS = $(GOFLAGS) -ldflags "$(LDFLAGS)"
TEST_FLAGS = -v -cover

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	$(MKDIR) $(BUILD_DIR)
	$(GO) build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)/$(MAIN_FILE)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

# Run the application
.PHONY: run
run: build
	$(BUILD_DIR)/$(APP_NAME)

# Clean build artifacts
.PHONY: clean
clean:
	$(RM) $(BUILD_DIR)
	@echo "Clean complete."

# Run tests
.PHONY: test
test:
	$(GO) test $(GOFLAGS) $(TEST_FLAGS) ./...

# Format code
.PHONY: fmt
fmt:
	$(GO) fmt ./...

# Initialize the database
.PHONY: db-init
db-init:
	sqlite3 storage/Beelzebub.db < pkg/db/schema.sql
	@echo "Database initialized."

# Generate documentation (requires godoc)
.PHONY: docs
docs:
	godoc -http=:6060
	@echo "Documentation available at http://localhost:6060"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all        - Build the application (default)"
	@echo "  build      - Build the application"
	@echo "  run        - Build and run the application"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"
	@echo "  fmt        - Format code"
	@echo "  db-init    - Initialize the database"
	@echo "  docs       - Start local documentation server"

