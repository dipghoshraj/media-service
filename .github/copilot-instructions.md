# agni-tunnels Project Guidelines

## Overview

Go-based reverse tunneling stack: a private/local service is exposed through an edge router using TLS + SNI-based routing. Three independent Go modules — keep them separate (no cross-module imports).

| Binary | Module | Role |
|--------|--------|------|
| `agni-agent` | `github.com/odio4u/agni-tunnels/agni-agent` | Runs near the private app; opens gRPC tunnel to router |
| `agni-router` | `github.com/odio4u/agni-tunnels/agni-router` | Accepts edge TLS, maps SNI → agent session, relays bytes |
| `agni-nova` | `github.com/odio4u/agni-tunnels/agni-nova` | Front-door TCP proxy; peeks SNI, forwards to correct router |

## Build and Test

```bash
make build-agent      # → bin/agni-agent.exe
make build-router     # → bin/agni-router.exe
make build-nova       # → bin/agni-nova-proxy.exe
make router-certs     # Generate router TLS certs (certmanger/certrouter.go)
make help             # Show all targets
```

Generate agent certs manually:
```bash
cd agni-agent
go run . gen-creds --dns <agent-id-or-domain> --name client
```

**There are no automated tests** (`*_test.go` files). Validate changes by building and running locally with the example stack in `example/`.

## Architecture

### Data path
1. External client → **nova** (SNI peek → seeder lookup → forward to router)
2. Router TCP proxy peeks SNI → finds agent session → sends `TunnelOpen{connection_id}`
3. Agent dials `localhost:<forward>` → exchanges `TunnelData` frames over gRPC bidirectional stream

### Tunnel protocol (all from `github.com/odio4u/agni-schema/tunnel`)
No local `.proto` files exist. All generated types live in the external `agni-schema` package. Key message variants on `tunnel.Envelope`:

| Variant | Direction | Meaning |
|---------|-----------|---------|
| `Envelope_Connect` | Agent → Router | Identity handshake (ID, Token, Timestamp, Signature) |
| `Envelope_ConnectAck` | Router → Agent | Accept/reject |
| `Envelope_Open` / `TunnelOpen` | Router → Agent | Dial local app with `ConnectionId` |
| `Envelope_Data` / `TunnelData` | Bidirectional | Raw payload bytes |
| `Envelope_Close` / `TunnelClose` | Bidirectional | Signal teardown |

### Identity / mTLS
Both sides use **self-signed certs + SHA-256 fingerprint pinning** (not a CA chain). `InsecureSkipVerify: true` combined with a custom `VerifyPeerCertificate` callback that compares the hex fingerprint against a value from config/seeder. **TLS 1.3 only** (`MinVersion: tls.VersionTLS13`).

Auth entry points: `agni-router/pkg/config/agentauth.go`, `agni-agent/pkg/rpc/client.go`.

## Key Conventions

### Configuration
- All three binaries read `agni-config.yaml` from the **current working directory** on startup.
- Config is loaded once at startup with `os.ReadFile` → `yaml.Unmarshal`; any error is fatal (`log.Fatalf`).
- **Important typo in YAML key**: `fingureprint` (not `fingerprint`) — preserve this in both YAML and struct tags to keep compatibility.

### Error handling
- Use the `FormatError` helper in `agni-agent/pkg/bridge/logger.go` for agent-side errors:
  ```go
  bridge.FormatError("message", err)  // → "[Agni Agent] message -- <err>"
  ```
- Router errors use inline `fmt.Errorf("[Agni Router] ...")` prefixes.
- Config/startup errors: `log.Fatalf(...)` — process terminates immediately.
- Stream errors: log then return; caller decides whether to reconnect.
- Always handle `io.EOF` explicitly on `stream.Recv()` as a clean shutdown signal.
- No custom error types; use the built-in `error` interface throughout.

### Logging
- Standard `log` package only (`log.Println`, `log.Printf`). No structured logger.
- Log prefix pattern: `[Agni Agent]`, `[Agni Router]`, `[Agni Nova]` — keep consistent.
- router gRPC server has panic recovery middleware (logs stack trace + returns `"internal server error"`).

### Concurrency
- Agent bidirectional relay is two goroutines: `PollStream` (router→local) and `LocaltoRpc` (local→router), both sharing a `context.Context` for cancellation.
- Router tracks active agent sessions in an in-memory map (see `agni-router/pkg/session/`).
- Use `sync.Mutex` or channels for session map access — check existing patterns in `session.go` before adding new state.

### Module layout
```
agni-agent/pkg/bridge/   — config loading, seeder calls, cert fingerprint
agni-agent/pkg/rpc/      — gRPC stream lifecycle + bidirectional forwarding
agni-agent/pkg/connector/ — local TCP dial to private app
agni-agent/cmd/          — Cobra CLI commands (connect, gen-creds, version)

agni-router/pkg/config/  — config loading, agent cert auth
agni-router/pkg/rpc/     — gRPC server implementation
agni-router/pkg/session/ — in-memory agent session registry
agni-router/server/      — TCP SNI listener
agni-router/certmanger/  — cert generation (note: folder name typo, preserve it)
```

## Key External Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/odio4u/agni-schema/tunnel` | Proto-generated tunnel message types |
| `github.com/odio4u/mem-sdk/memsdk` | Seeder client: register, lookup agents/gateways |
| `github.com/odio4u/mem-sdk/certengine` | Self-signed cert generation (agent only) |
| `github.com/odio4u/mem-sdk/sni` | SNI peeking (nova + router) |
| `github.com/spf13/cobra` | CLI framework (agent only) |
| `google.golang.org/grpc` | gRPC framework |
| `github.com/grpc-ecosystem/go-grpc-middleware` | Recovery/panic interceptors (router only) |
| `github.com/google/uuid` | Connection IDs (router only) |
| `gopkg.in/yaml.v3` | Config parsing (all modules) |
