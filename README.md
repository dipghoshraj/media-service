# agni-tunnels

agni-tunnels is a Go-based reverse tunneling stack that lets a private/local service be exposed through an edge router using TLS + SNI-based routing.

At a high level it has three binaries:

- **agni-agent**: runs near your private app and opens a persistent gRPC tunnel to a router.
- **agni-router**: accepts edge TLS traffic, maps SNI -> agent session, and relays bytes over gRPC.
- **agni-nova**: a front-door TCP proxy that reads client SNI and forwards to the right router/gateway address discovered from a seeder service.

The project currently uses a shared config file (`agni-config.yaml`) and external service discovery/identity APIs from `mem-sdk`.

## How the system works

### 1) Discovery and registration

1. The **router** starts and reads `Router` settings from `agni-config.yaml`.
2. It computes its certificate fingerprint and registers itself in the seeder (`SeedGatway`).
3. The **agent** starts, reads `Agent` settings, asks the seeder for a router/gateway in a region, and registers itself.
4. The agent receives routing metadata including gateway address/ports and connects over gRPC with mTLS-like identity checks based on certificate fingerprint.

### 2) Tunnel establishment

1. Agent opens `AgniTunnel.Connect` bidirectional stream to router.
2. First message is `ConnectRequest`; router replies with `ConnectAck`.
3. Router stores the agent stream in an in-memory session registry keyed by mapped domain.

### 3) Data path (client -> local app)

1. External client reaches nova (or router directly) with TLS SNI (example: `agni.local.internal`).
2. Router-side TCP proxy (`RouterServer`) accepts raw TCP, peeks SNI, and finds matching agent session.
3. Router sends `TunnelOpen{connection_id}` to agent.
4. Agent dials local app on `localhost:<forward>` from config.
5. Router and agent exchange `TunnelData` frames over gRPC, forwarding bytes both directions until close/error.

## Repository layout

- `agni-agent/`
  - CLI app (`connect`, `gen-creds`, `version`) with Cobra.
  - `pkg/bridge`: config loading, seeder/registry calls, cert fingerprint helpers.
  - `pkg/rpc`: gRPC client stream lifecycle and bidirectional forwarding logic.
  - `pkg/connector`: local TCP connection to your private app.

- `agni-router/`
  - gRPC server with TLS cert loading and custom peer certificate validation.
  - Registers itself in seeder and keeps in-memory agent session/domain maps.
  - `server/`: TCP listener that receives inbound SNI traffic and hooks it to agent stream.

- `agni-nova/`
  - Lightweight TCP proxy.
  - Peeks SNI and asks seeder for the correct gateway/router address.
  - Pipes bytes between client and selected backend router.

- `example/`
  - Sample Flask + Gunicorn app behind nginx TLS.
  - Useful for simulating the local service that agent forwards to.

- `agni-config.yaml`
  - Shared config for Agent, Router, and Nova.

## Configuration quick notes

`agni-config.yaml` currently defines:

- `Agent.forward`: local port agent connects to (e.g. `5050`).
- `Agent.domain`: domain used for mapping.
- `Router.rpc_port`: router gRPC tunnel port.
- `Router.proxy_port`: router inbound proxy/TCP listener port.
- `certs` folders for agent/router certificate files.
- `Seeder.address` + `Seeder.fingureprint` for service discovery/identity checks.

## Build commands

Use the Makefile targets:

```bash
make build-agent
make router-certs
make build-router
make build-nova
```

Outputs are placed in `bin/` with `.exe` suffix in current Make targets.

## Typical local run flow

1. Prepare/update `agni-config.yaml` for your environment.
2. Generate router certs:
   ```bash
   make router-certs
   ```
3. Generate agent certs:
   ```bash
   cd agni-agent
   go run . gen-creds --dns <agent-id-or-domain> --name client
   ```
4. Start router:
   ```bash
   make build-router && ./bin/agni-router.exe
   ```
5. Start your local app on `localhost:<Agent.forward>`.
6. Start agent:
   ```bash
   make build-agent && ./bin/agni-agent.exe connect
   ```
7. (Optional) Start nova front proxy:
   ```bash
   make build-nova && ./bin/agni-nova-proxy.exe
   ```

## Current caveats / things to know

- Some naming/typos exist in code (`fingureprint`, `ProxtPort`, etc.) but are consistent with config tags.
- Router `Connect` currently ends with `select {}` (placeholder), so stream lifecycle cleanup is minimal.
- Authentication logic is partially stubbed (`checkConnect` always true), while cert fingerprint checks are active.
- Error handling is still evolving in several places.

