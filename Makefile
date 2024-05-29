# Makefile for building a Go program for multiple architectures

# Program name
BINARY_NAME=sunat

# Source file
SRC_FILE=main.go

# Build directory
BUILD_DIR=bin

# Architectures to build for
ARCHS = \
    darwin/amd64 \
    darwin/arm64 \
    linux/386 \
    linux/amd64 \
    linux/arm64 \
    windows/386 \
    windows/amd64

# Default target: build for all architectures
all: $(ARCHS)

# Rule to build for a specific architecture
$(ARCHS):
	@GOOS=$(word 1,$(subst /, ,$@)) GOARCH=$(word 2,$(subst /, ,$@)) \
	go build -o $(BUILD_DIR)/$(BINARY_NAME)_$(word 1,$(subst /, ,$@))_$(word 2,$(subst /, ,$@)) $(SRC_FILE)
	@echo "Built $(BINARY_NAME) for $(word 1,$(subst /, ,$@))/$(word 2,$(subst /, ,$@))"

# Clean up build artifacts
clean:
	rm -rf $(BUILD_DIR)

.PHONY: all clean $(ARCHS)
