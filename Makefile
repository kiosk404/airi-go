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

# 定义thriftgo的通用参数
THRIFT_OUT_GO := ./backend/api/model
THRIFT_PACKAGE_PREFIX_GO := github.com/kiosk404/airi-go/backend/api/model
THRIFT_OUT_JS := ./frontend/src/api
THRIFT_PACKAGE_PREFIX_JS :=

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
	@-rm -vrf $(BACKEND_OUTPUT_DIR)
	@-rm -vrf $(BACKEND_DIR)/crash.log

## format: Gofmt (reformat) package sources (exclude vendor dir if existed).
.PHONY: format
format: tools.verify.golines tools.verify.goimports
	@echo "===========> Formating codes"
	@$(FIND) -type f -name '*.go' | $(XARGS) gofmt -s -w
	@$(FIND) -type f -name '*.go' | $(XARGS) goimports -w -local $(BACKEND_DIR)
	@$(FIND) -type f -name '*.go' | $(XARGS) golines -w --max-len=240 --reformat-tags --shorten-comments --ignore-generated .
	@cd $(BACKEND_DIR) && $(GO) mod edit -fmt && cd -

## idl-go: Generate Go code from Thrift IDL files in ./idl directory.PHONY: idl-go
idl-go:
	@thriftgo -r -gen go:package_prefix=$(THRIFT_PACKAGE_PREFIX_GO)/,template=slim,with_context=true \
		-out $(THRIFT_OUT_GO) ./idl/api.thrift
