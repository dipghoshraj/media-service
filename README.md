# agni-Agent — CLI

The **agni-agent** is the client-side component of the agni-tunnels reverse tunneling stack. It runs alongside your private/local application and opens a persistent gRPC tunnel to an agni-router, allowing external traffic to reach services that aren't directly exposed to the internet.

The agent uses TLS with SHA-256 certificate fingerprint pinning (no CA chain) and discovers routers through a seeder service (`mem-sdk`).

## How it works

### 1) Discovery and registration

1. Agent reads `Agent` settings from `agni-config.yaml` in the current working directory.
2. Connects to the **seeder** service and looks up available gateways for the configured `region`.
3. Computes the local certificate's SHA-256 fingerprint and registers itself with the seeder (`ConnectAgent`).
4. Receives routing metadata — gateway IP, ports, and the gateway's certificate fingerprint.

### 2) Tunnel establishment

1. Agent dials the router over gRPC with mTLS-style identity. `InsecureSkipVerify: true` is set, but a custom `VerifyPeerCertificate` callback rejects any server whose certificate fingerprint doesn't match the gateway identity returned by the seeder.
2. Opens a bidirectional `AgniTunnel.Connect` stream.
3. Sends a `ConnectRequest` envelope containing agent ID, domain, timestamp, and signature.
4. Waits for a `ConnectAck` from the router.

### 3) Data path (router → local app)

1. When an external client connects (via router/nova SNI routing), the router sends a `TunnelOpen` envelope with a unique `connection_id`.
2. Agent dials `host:forward` (from config) over TCP and maps the local connection to the `connection_id`.
3. A `LocaltoRpc` goroutine reads bytes from the local connection and forwards them as `TunnelData` envelopes to the router.
4. Incoming `TunnelData` envelopes from the router are written to the matching local connection.
5. On close or error, a `TunnelClose` envelope is sent and the local connection is cleaned up.

## Repository layout

```
agni-agent/
├── main.go                   # Entry point — calls cmd.Execute()
├── go.mod
├── Makefile                  # Cross-platform build targets
├── cmd/
│   ├── rootcmd.go            # Cobra root command (agent-tunnel)
│   ├── connect.go            # `connect` — register + open tunnel
│   ├── gencreds.go           # `gen-creds` — generate TLS certificates
│   └── version.go            # `version` — print version (currently 0.1.2)
├── pkg/
│   ├── bridge/
│   │   ├── load-yaml.go      # Config struct + YAML loading (init)
│   │   ├── agentconnect.go   # Seeder calls: fingerprint, gateway lookup, registration
│   │   ├── build_creds.go    # Delegates to certengine for self-signed cert generation
│   │   └── logger.go         # Structured JSON logger (log/slog) + FormatError helper
│   ├── rpc/
│   │   ├── client.go         # gRPC dial with TLS fingerprint verification
│   │   ├── connect.go        # TunnelSession: stream setup, ConnectRequest, signal handling
│   │   ├── polling.go        # PollStream: recv loop dispatching ConnectAck/Open/Data/Close
│   │   └── streamdata.go     # HandleStream (router→local write) + LocaltoRpc (local→router)
│   └── connector/
│       ├── connector.go      # TCP dial to local app (host:forward from config)
│       └── close.go          # SendClose helper for TunnelClose envelopes
doc/
├── agni-agent-quickstart.md
└── agni-router-quickstart.md
agni-config.yaml              # Shared config (Agent, Router, Nova sections)
Makefile                      # Root-level shortcuts for agent build targets
```

## Configuration

The agent reads `agni-config.yaml` from the current working directory at startup. The config is loaded once in a package `init()` function; any error causes an immediate exit.

```yaml
version: v1

Agent:
  name: "agent-agni"
  domain: "agni.local.internal"   # Domain used for SNI routing
  forward: 5050                   # Local port your app listens on
  host: "localhost"               # Local host to dial
  region: "global"                # Region for gateway lookup
  certs: "./"                     # Directory containing client.pem + client-key.pem
  Seeder:
    address: "localhost:8080"     # Seeder service address
    fingureprint: "<seeder-cert-fingerprint>"  # Note: typo is intentional, preserved for compat
```

| Field | Purpose |
|-------|---------|
| `domain` | The SNI domain the router uses to map traffic to this agent |
| `forward` | TCP port of the local application to tunnel to |
| `host` | Hostname/IP for the local dial (combined with `forward`) |
| `certs` | Path to `client.pem` and `client-key.pem` certificate files |
| `region` | Passed to the seeder when looking up gateways |
| `Seeder.address` | Address of the seeder/discovery service |
| `Seeder.fingureprint` | SHA-256 fingerprint of the seeder's TLS cert (typo preserved) |

## CLI commands

The agent CLI is built with [Cobra](https://github.com/spf13/cobra). The root command is `agent-tunnel`.

### `connect`

Register the agent with the seeder, discover a gateway, and open a persistent tunnel:

```bash
./release/windows/agni-agent.exe connect
```

### `gen-creds`

Generate self-signed TLS certificates for the agent:

```bash
./release/windows/agni-agent.exe gen-creds --dns <agent-id-or-domain> --name client
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--dns` | Yes | — | DNS SAN to include in the certificate |
| `--name` | Yes | `agent-tunnel` | Base filename for the generated cert/key (`<name>.pem`, `<name>-key.pem`) |

### `version`

Print the current version:

```bash
./release/windows/agni-agent.exe version
```

## Build

### From the repo root (shortcuts)

```bash
make agent-linux      # Linux amd64   → release/linux/agni-agent
make agent-darwin     # macOS amd64   → release/darwin/agni-agent
make agent-windows    # Windows amd64 → release/windows/agni-agent.exe
make agent-all        # All platforms
```

### From the agent directory

```bash
cd agni-agent
make build-linux
make build-darwin
make build-windows
make build-all
```

### Generate certs via Make

```bash
cd agni-agent
make gen-creds DNS=my-agent-id NAME=client
```

### During development

```bash
cd agni-agent
go run . connect
go run . gen-creds --dns my-agent-id --name client
```

## Quick start

1. **Configure** — edit `agni-config.yaml` with your seeder address, agent domain, forward port, and host.

2. **Generate certificates**:
   ```bash
   cd agni-agent
   go run . gen-creds --dns <agent-domain> --name client
   ```

3. **Start your local app** on the port specified in `Agent.forward` (e.g. `:5050`).

4. **Run the agent**:
   ```bash
   make agent-windows
   ./release/windows/agni-agent.exe connect
   ```
   The agent will register with the seeder, resolve a gateway, verify the gateway's certificate fingerprint, and establish a bidirectional gRPC tunnel.

## Key dependencies

| Package | Purpose |
|---------|---------|
| `github.com/odio4u/agni-schema/tunnel` | Proto-generated tunnel message types (Envelope, ConnectRequest, TunnelData, etc.) |
| `github.com/odio4u/mem-sdk/memsdk` | Seeder client: gateway lookup, agent registration |
| `github.com/odio4u/mem-sdk/certengine` | Self-signed certificate generation |
| `github.com/spf13/cobra` | CLI framework |
| `google.golang.org/grpc` | gRPC client |
| `gopkg.in/yaml.v3` | YAML config parsing |

## Logging

The agent uses Go's `log/slog` package with a JSON handler. All log entries carry a `"component": "agni-agent"` field. Example output:

```json
{"time":"...","level":"INFO","msg":"config loaded","component":"agni-agent","file":"agni-config.yaml","version":"v1"}
{"time":"...","level":"INFO","msg":"agent registered","component":"agni-agent","agent_id":"...","domain":"agni.local.internal"}
```

The `FormatError` helper (`pkg/bridge/logger.go`) is available for building prefixed error strings:

```go
bridge.FormatError("message", err)  // → "[Agni Agent] message -- <err>"
```

## Notes

- The YAML field `fingureprint` (not `fingerprint`) is an intentional preserved typo — keep it consistent in config files and struct tags.
- TLS verification uses `InsecureSkipVerify: true` paired with a custom `VerifyPeerCertificate` callback that checks the SHA-256 fingerprint. This is by design (fingerprint pinning, not CA-based trust).
- Stream sends are protected by a `sync.Mutex` (`sendMu`) to prevent concurrent writes on the gRPC stream.
- Graceful shutdown is handled via `os.Signal` (interrupt) — the agent cancels the context, closes the gRPC connection, and waits for the poll goroutine to finish.

