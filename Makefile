build-agent:
	cd agni-agent && go build -o ../bin/agni-agent.exe main.go

router-certs:
	cd agni-router && go run certmanger/certrouter.go

build-router:
	cd agni-router && go build -o ../bin/agni-router.exe main.go

build-nova:
	cd agni-nova && go build -o ../bin/agni-nova.exe main.go

## Show this help
help:
	@echo "Usage:"
	@echo "make build-agent     Build agni agent"
	@echo "make router-certs    Generate router certificates"
	@echo "make build-router    Build agni router"
	@echo "make build-nova      Build agni nova"


.PHONY: build router-certs help