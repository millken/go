# PlutoVG Go Bindings Makefile
# Builds static libraries for PlutoVG and manages source synchronization

# Project configuration
PROJECT_NAME = plutovg
PLUTOVG_VERSION = main
PLUTOVG_REPO = https://github.com/sammycage/plutovg.git

# Directory structure
SRC_DIR = src
INCLUDE_DIR = include
LIBS_DIR = _libs
BUILD_DIR = build

# Platform detection
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

ifeq ($(UNAME_S),Linux)
    PLATFORM = linux
    CC = gcc
    AR = ar
    CFLAGS = -O3 -fPIC -DPLUTOVG_BUILD_STATIC
    LDFLAGS = -lm
endif

ifeq ($(UNAME_S),Darwin)
    PLATFORM = darwin
    CC = clang
    AR = ar
    CFLAGS = -O3 -fPIC -DPLUTOVG_BUILD_STATIC
    LDFLAGS = -lm
endif

ifeq ($(UNAME_M),x86_64)
    ARCH = amd64
endif

ifeq ($(UNAME_M),arm64)
    ARCH = arm64
endif

ifeq ($(UNAME_M),aarch64)
    ARCH = arm64
endif

# Target directories
PLATFORM_DIR = $(PLATFORM)/$(ARCH)
TARGET_LIB_DIR = $(LIBS_DIR)/$(PLATFORM_DIR)
TARGET_LIB = $(TARGET_LIB_DIR)/lib$(PROJECT_NAME).a

# Default target
.PHONY: all
all: sync build

# Sync PlutoVG source code
.PHONY: sync
sync:
	@echo "Syncing PlutoVG source code..."
	@if [ ! -d "$(SRC_DIR)" ]; then \
		echo "Cloning PlutoVG repository..."; \
		git clone $(PLUTOVG_REPO) $(SRC_DIR); \
	else \
		echo "Updating PlutoVG repository..."; \
		cd $(SRC_DIR) && git fetch origin && git reset --hard origin/$(PLUTOVG_VERSION); \
	fi
	@echo "Copying headers to include directory..."
	@mkdir -p $(INCLUDE_DIR)
	@cp $(SRC_DIR)/*.h $(INCLUDE_DIR)/ 2>/dev/null || true
	@echo "Source sync completed."

# Build static library
.PHONY: build
build:
	@echo "Building static library..."
	cd $(SRC_DIR) && \
	cmake -B build . && \
	cmake --build build
	@echo "Static library created successfully."

