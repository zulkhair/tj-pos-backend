# Default values for OS and architecture
OS ?= linux
ARCH ?= amd64

# Binary output directory
BUILD_DIR = build

# Build command with OS and architecture options
.PHONY: build
build:
	@echo "Building for $(OS)/$(ARCH)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(OS) GOARCH=$(ARCH) go build -o $(BUILD_DIR)/tj-pos-$(OS)-$(ARCH) ./main.go
	@echo "Build complete: $(BUILD_DIR)/tj-pos-$(OS)-$(ARCH)"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
