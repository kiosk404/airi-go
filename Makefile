# Build all by default, even if it's not first
.DEFAULT_GOAL := all

.PHONY: all
#all: tidy format lint cover build
all: tidy format lint build clean

# ==============================================================================
# Build options

ROOT_PACKAGE=airi-go
FRONTEND_VERSION_PACKAGE=airi-go
BACKEND_VERSION_PACKAGE=airi-go/backend/pkg/version

# ==============================================================================
# Includes

include scripts/make-rules/common.mk # make sure include common.mk at the first include line
include scripts/make-rules/golang.mk
include scripts/make-rules/tools.mk

# ==============================================================================
# Usage

define USAGE_OPTIONS

Options:
  DEBUG            Whether to generate debug symbols. Default is 0.
  BINS             The binaries to build. Default is all of cmd.
                   This option is available when using: make build/build.multiarch
                   Example: make build BINS="goserver test"
  VERSION          The version information compiled into binaries.
                   The default is obtained from gsemver or git.
  V                Set to 1 enable verbose build. Default is 0.
endef
export USAGE_OPTIONS

# ==============================================================================
# Targets

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)
	@-rm -vrf $(OUTPUT_BIN_DIR)
