build-agent:
	cd agni-agent && go build -o ../bin/agni-agent.exe main.go


.PHONY: build