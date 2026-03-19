BIN_DIR  := ../release
BINARY   := agni-agent

.PHONY: build-linux build-darwin build-windows build-all gen-creds help

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

## Generate agent TLS credentials (pass DNS and --name flags)
## Usage: make gen-creds DNS=my-agent-id NAME=client
gen-creds:
	go run . gen-creds --dns $(DNS) --name $(NAME)

## Show this help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build-linux    Build agni-agent for Linux   -> $(BIN_DIR)/$(BINARY)-linux-amd64"
	@echo "  build-darwin   Build agni-agent for macOS   -> $(BIN_DIR)/$(BINARY)-darwin-amd64"
	@echo "  build-windows  Build agni-agent for Windows -> $(BIN_DIR)/$(BINARY)-windows-amd64.exe"
	@echo "  build-all      Build all platform targets"
	@echo "  gen-creds      Generate agent TLS credentials"
	@echo "                 Required vars: DNS=<agent-id-or-domain> NAME=<cert-name>"
	@echo "  Example: make gen-creds DNS=my-agent NAME=client"
