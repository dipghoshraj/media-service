# agni-Agent

**agni-agent** runs alongside your private application and opens a persistent gRPC tunnel to an agni-router, letting external traffic reach services that aren't directly exposed to the internet.

## Installation

### Download pre-built binary (recommended)

Pre-built binaries for Linux, macOS, and Windows are attached to every [GitHub Release](https://github.com/dipghoshraj/agni-stack/releases).

1. Go to the [Releases page](https://github.com/dipghoshraj/agni-stack/releases) and download the binary for your platform:

| Platform | File |
|----------|------|
| Linux (amd64) | `agni-agent` |
| macOS (amd64) | `agni-agent` |
| Windows (amd64) | `agni-agent.exe` |

2. Make it executable (Linux / macOS):
   ```bash
   chmod +x agni-agent
   sudo mv agni-agent /usr/local/bin/agni-agent
   ```

### Build from source

Requires Go 1.21+ and `make`.

```bash
git clone https://github.com/dipghoshraj/agni-stack.git
cd agni-stack

make agent-linux      # → release/linux/agni-agent
make agent-darwin     # → release/darwin/agni-agent
make agent-windows    # → release/windows/agni-agent.exe
make agent-all        # all platforms at once
```

Run `make help` to list all available targets.

## Configuration

Place an `agni-config.yaml` in the directory where you run the agent:

```yaml
version: v1

Agent:
  name: "agent-agni"
  domain: "agni.local.internal"   # SNI domain for routing
  forward: 5050                   # Local port your app listens on
  host: "localhost"               # Local host to dial
  region: "global"
  certs: "./"                     # Directory with client.pem + client-key.pem
  Seeder:
    address: "localhost:8080"
    fingureprint: "<seeder-cert-fingerprint>"
```

| Field | Description |
|-------|-------------|
| `domain` | SNI domain the router uses to route traffic to this agent |
| `forward` | TCP port of your local application |
| `host` | Hostname/IP for the local dial |
| `certs` | Path containing `client.pem` and `client-key.pem` |
| `Seeder.address` | Address of the seeder/discovery service |
| `Seeder.fingureprint` | SHA-256 fingerprint of the seeder's TLS certificate |

## CLI commands

### `connect`

Register with the seeder and open a persistent tunnel:

```bash
agni-agent connect
```

### `gen-creds`

Generate self-signed TLS certificates using values from `agni-config.yaml` (`Agent.domain` and `Agent.name`):

```bash
agni-agent gen-creds
```

To point at a different config file:

```bash
agni-agent gen-creds -f /path/to/agni-config.yaml
```

| Flag | Description |
|------|-------------|
| `-f` | Path to config file (default: `agni-config.yaml` in the current directory) |

### `version`

```bash
agni-agent version
```

## Quick start

1. **Edit `agni-config.yaml`** with your seeder address, domain, forward port, and cert path.

2. **Generate certificates**:
   ```bash
   agni-agent gen-creds
   ```

3. **Start your local app** on the port set in `Agent.forward`.

4. **Run the agent**:
   ```bash
   agni-agent connect
   ```

