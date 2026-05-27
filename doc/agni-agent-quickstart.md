# Quick Start: Agni Agent

This guide helps you quickly set up and run the Agni Agent, configure `agni-config.yaml`, and connect using the CLI.

## 1. Prepare `agni-config.yaml`

Ensure your config file includes the following under the `Agent` section:

```yaml
Agent:
  name: "agent-agni"
  domain: "agni.local.internal"
  forward: 5050           # Local port to forward traffic to
  region: "global"
  certs: "agni-agent/certs"   # Path to agent certificate files
  Seeder:
    address: "localhost:8080" # Seeder service address
    fingureprint: "<router-cert-fingerprint>"
```

- `domain`: The domain mapped for SNI routing.
- `forward`: Local port your app listens on (e.g., 5050).
- `certs`: Directory for agent TLS certificates.
- `Seeder`: Discovery service for router/gateway lookup.

## 2. Generate Agent Certificates

```bash
cd agni-agent
go run . gen-creds --dns <agent-id-or-domain> --name client
```

## 3. Start Your Local App

Run your app on the port specified in `Agent.forward` (e.g., 5050).

## 4. Start Agni Agent and Connect

Build and run the agent:

```bash
make build-agent
./bin/agni-agent.exe connect
```

Or use `go run` during development:

```bash
cd agni-agent
go run . connect
```

This registers the agent with the seeder, discovers the router, and establishes a tunnel.
