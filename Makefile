
ROOT_DIR ?= $(CURDIR)
BUILD_DIR := $(ROOT_DIR)/_output
BIN_DIR := $(BUILD_DIR)/bin

BUILD_SCRIPT := $(ROOT_DIR)/build/build.sh
DEB_BUILD_SCRIPT := $(ROOT_DIR)/build/build_deb.sh

GO_INSTALL := go install 
GO_FLAGS := GOFLAGS=-mod=mod


PKGS := go list ./...

install:
	@for PKG in $$( $(PKGS) ); do \
		echo $$PKG; \
		$(GO_INSTALL) $$PKG; \
	done

static-files:
	@statik -src=./web -f

build-bin:
	$(GO_FLAGS) go build -o $(BIN_DIR)/npv main.go
	