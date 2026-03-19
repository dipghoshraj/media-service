RELEASE_DIR := ../release
BINARY      := agni-nova

.PHONY: build-linux build-darwin build-windows build-all help

## Build for Linux (amd64)
build-linux:
	-mkdir "$(RELEASE_DIR)\linux" 2>nul
	cmd /C "set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0&& go build -o $(RELEASE_DIR)/linux/$(BINARY) ."

## Build for macOS (amd64)
build-darwin:
	-mkdir "$(RELEASE_DIR)\darwin" 2>nul
	cmd /C "set GOOS=darwin&& set GOARCH=amd64&& set CGO_ENABLED=0&& go build -o $(RELEASE_DIR)/darwin/$(BINARY) ."

## Build for Windows (amd64)
build-windows:
	-mkdir "$(RELEASE_DIR)\windows" 2>nul
	cmd /C "set GOOS=windows&& set GOARCH=amd64&& set CGO_ENABLED=0&& go build -o $(RELEASE_DIR)/windows/$(BINARY).exe ."

## Build all platform targets
build-all: build-linux build-darwin build-windows

## Show this help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build-linux    Build agni-nova for Linux   -> $(RELEASE_DIR)/linux/$(BINARY)"
	@echo "  build-darwin   Build agni-nova for macOS   -> $(RELEASE_DIR)/darwin/$(BINARY)"
	@echo "  build-windows  Build agni-nova for Windows -> $(RELEASE_DIR)/windows/$(BINARY).exe"
	@echo "  build-all      Build all platform targets"
