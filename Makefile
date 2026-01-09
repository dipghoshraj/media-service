build-agent:
	cd agni-agent && go build -o ../bin/agni-agent.exe main.go

router-certs:
	cd agni-router && go run certmanger/certrouter.go


.PHONY: build router-certs