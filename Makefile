build-agent:
	cd agni-agent && go build -o ../bin/agni-agent main.go

router-certs:
	cd agni-router && go run certmanger/certrouter.go

build-router:
	cd agni-router && go build -o ../bin/agni-router main.go

build-nova:
	cd agni-nova && go build -o ../bin/agni-nova-proxy main.go


build-all:
	@echo "Building Router..."
	cmd /C "set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0&& cd agni-router && go build -o release/agni-router-linux-amd64 main.go"
	@echo "Building agent..."
	cmd /C "set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0&& cd agni-agent && go build -o release/agni-agent-linux-amd64 main.go"
	@echo "Building nova..."
	cmd /C "set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0&& cd agni-nova && go build -o release/agni-nova-linux-amd64 main.go"
	@echo "All builds completed."

## Show this help
help:
	@echo "Usage:"
	@echo "make build-agent     Build agni agent"
	@echo "make router-certs    Generate router certificates"
	@echo "make build-router    Build agni router"
	@echo "make build-nova      Build agni nova"


.PHONY: build router-certs help