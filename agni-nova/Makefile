BIN_DIR  := ../release
BINARY   := agni-nova

.PHONY: build-linux build-darwin build-windows build-all help

## Build for Linux (amd64)
build-linux:
	-mkdir "$(BIN_DIR)" 2>nul
	cmd /C "set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0&& go build -o $(BIN_DIR)/$(BINARY)-linux-amd64 ."

## Build for macOS (amd64)
build-darwin:
	-mkdir "$(BIN_DIR)" 2>nul
	cmd /C "set GOOS=darwin&& set GOARCH=amd64&& set CGO_ENABLED=0&& go build -o $(BIN_DIR)/$(BINARY)-darwin-amd64 ."

## Build for Windows (amd64)
build-windows:
	-mkdir "$(BIN_DIR)" 2>nul
	cmd /C "set GOOS=windows&& set GOARCH=amd64&& set CGO_ENABLED=0&& go build -o $(BIN_DIR)/$(BINARY)-windows-amd64.exe ."

## Build all platform targets
build-all: build-linux build-darwin build-windows

## Show this help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build-linux    Build agni-nova for Linux   -> $(BIN_DIR)/$(BINARY)-linux-amd64"
	@echo "  build-darwin   Build agni-nova for macOS   -> $(BIN_DIR)/$(BINARY)-darwin-amd64"
	@echo "  build-windows  Build agni-nova for Windows -> $(BIN_DIR)/$(BINARY)-windows-amd64.exe"
	@echo "  build-all      Build all platform targets"
