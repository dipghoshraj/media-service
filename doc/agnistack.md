# Agnistack — Privacy-First Reverse Tunnel Network

Agnistack is a privacy-first, decentralized application deployment network. It is built to expose **private servers** to the internet through a distributed routing network — without TLS termination, without data inspection, and without giving up control of your certificates or domain.

You bring your own domain, your own SSL certs, and your own server. Agnistack routes the raw TCP stream directly to it. Nothing more.

While it works on localhost for development and testing, the primary use case is exposing a **privately hosted server** — a VPS, a bare-metal machine behind NAT, or a server on a restricted network — without opening firewall ports or changing your network setup.

> **Built for zero-trust, anonymous access.** Originally designed for privacy-first applications like end-to-end encrypted chat, Agnistack lets your server stay reachable even in environments with bans, port blocks, or network restrictions.

---

## How it Works

```
External Client
      │
      ▼
 agni-nova          ← Entry point: peeks SNI, routes to correct router
      │
      ▼
 agni-router        ← Receives TCP stream; maps SNI → agent session
      │  (gRPC bidirectional stream)
      ▼
 agni-agent         ← Runs on your private server; forwards to your app
      │
      ▼
 Your Application   ← localhost:<port>
```

1. An external client connects to **agni-nova** (the front-door TCP proxy). Nova peeks the SNI from the TLS handshake and consults the seeder to find the right router.
2. **agni-router** receives the TCP stream, looks up which agent session owns that SNI domain, and sends a `TunnelOpen` signal over an already-established gRPC stream.
3. **agni-agent** (running next to your app) dials `localhost:<your-port>`, then relays raw bytes in both directions over the gRPC stream.
4. The response travels back the same way — nova → router → agent → your app.

**There is no TLS termination at any hop.** Your certs stay yours.

---

## Components

| Component | Role |
|-----------|------|
| **agni-nova** | Front-door TCP proxy; peeks SNI and forwards traffic to the correct router |
| **agni-router** | Accepts edge TLS traffic; maps SNI → agent gRPC session; relays bytes |
| **agni-seeder** | Registry service; holds metadata about agents and routers (gateways) |
| **agni-agent** | Runs on your private server; connects to the network and tunnels traffic to your app |

All components are open source. You can run your own seeder and nova, or connect to the hosted network.

---

## Identity & Security

- Each agent generates a **self-signed TLS certificate**. The SHA-256 fingerprint of that cert becomes the agent's identity.
- The fingerprint is registered with the seeder on connect. The router verifies it on every connection.
- **TLS 1.3 only.** `InsecureSkipVerify: true` is used for the gRPC dial, but a custom `VerifyPeerCertificate` callback rejects any peer whose fingerprint doesn't match what the seeder returned.
- No CA chain, no certificate authority to trust or be compromised.

---

## Getting Started with agni-agent

### Prerequisites

- A domain with DNS access
- A private server or local machine running your application
- `agni-agent` binary ([download](#build) or build from source)

---

### Step 1 — Discover Available Seeders

Run the scan command to list available seeders in the network:

```bash
agni-agent scan
```

This prints a table of seeders (address, region, fingerprint). Pick one to use in your config.

---

### Step 2 — Configure `agni-config.yaml`

Create or edit `agni-config.yaml` in the directory where you'll run the agent:

```yaml
version: v1

Agent:
  name: "my-agent"
  domain: "myapp.example.com"   # Your domain — used for SNI routing
  forward: 443                  # Local port your app listens on
  host: "localhost"             # Local host to dial
  region: "global"
  certs: "./"                   # Directory containing client.pem + client-key.pem
  Seeder:
    address: "45.130.164.217:8080"      # Seeder address from scan output
    fingureprint: "<seeder-fingerprint>" # Seeder cert fingerprint from scan output
```

| Field | Description |
|-------|-------------|
| `domain` | The domain your users connect to (must match your DNS record) |
| `forward` | The port your application listens on (on the private server) |
| `host` | Host to dial on the private server — `localhost` if the app runs on the same machine as the agent, or a LAN IP if it runs elsewhere |
| `certs` | Path to your `client.pem` and `client-key.pem` files |
| `region` | Region for gateway lookup (e.g., `global`) |
| `Seeder.address` | Address of the seeder service |
| `Seeder.fingureprint` | SHA-256 fingerprint of the seeder's TLS cert |

---

### Step 3 — Generate Agent Credentials

Generate a self-signed TLS certificate for your agent's identity:

```bash
agni-agent gen-creds --dns myapp.example.com --name client
```

This creates `client.pem` and `client-key.pem` in the current directory (or wherever `certs` points in your config).

| Flag | Required | Description |
|------|----------|-------------|
| `--dns` | Yes | DNS SAN to embed in the certificate (use your domain or agent ID) |
| `--name` | Yes | Base filename for the generated files (`<name>.pem`, `<name>-key.pem`) |

---

### Step 4 — Point Your Domain to the Network

After connecting (Step 5), the agent prints the nova address. Add a DNS `A` record for your domain pointing to that address.

```
myapp.example.com.  A  <nova-ip>
```

---

### Step 5 — Connect

Make sure your application is running on the private server, then run the agent on that same machine (or anywhere with network access to `host:forward`):

```bash
agni-agent connect
```

The agent will:
1. Read config from `agni-config.yaml`
2. Compute the certificate fingerprint and register with the seeder
3. Discover the assigned router (gateway)
4. Open a persistent gRPC tunnel to the router
5. Begin relaying traffic to your application

Leave it running. Any external request to your domain will be tunneled through, end-to-end, to your private server in real time.

> **Localhost also works.** If you are running both the agent and your app on a local machine (e.g., for development), set `host: localhost` and `forward` to your dev port. The tunnel behaves identically.

---

## CLI Reference

### `agni-agent scan`

List available seeders in the network.

```bash
agni-agent scan
```

### `agni-agent connect`

Register, discover a gateway, and open a persistent tunnel.

```bash
agni-agent connect
```

Reads `agni-config.yaml` from the current working directory.

### `agni-agent gen-creds`

Generate self-signed TLS certificates for agent identity.

```bash
agni-agent gen-creds --dns <domain-or-id> --name <filename-base>
```

### `agni-agent version`

Print the current agent version.

```bash
agni-agent version
```

---

## Build from Source

### Prerequisites

- Go 1.21+
- `make`

### Quick build (from repo root)

```bash
make agent-windows   # → release/windows/agni-agent.exe
make agent-linux     # → release/linux/agni-agent
make agent-darwin    # → release/darwin/agni-agent
make agent-all       # All platforms
```

### From the agent directory

```bash
cd agni-agent
make build-windows
make build-linux
make build-darwin
```

### Run without building

```bash
cd agni-agent
go run . connect
go run . gen-creds --dns myapp.example.com --name client
```

---

## Full Example

```bash
# 1. Scan for seeders
agni-agent scan

# 2. Edit agni-config.yaml with your domain, port, and chosen seeder

# 3. Generate certs
agni-agent gen-creds --dns myapp.example.com --name client

# 4. Start your application on the private server (port 3000 in this example)
node server.js

# 5. On the same machine, update forward: 3000 in agni-config.yaml, then connect
agni-agent connect
```

Once connected, set your domain's A record to the nova address printed by the agent. Your private server is now reachable from the internet — no open ports, no firewall changes, no TLS termination.

> **For local development:** the same steps apply. Set `host: localhost` and `forward` to your dev port. The agent tunnels traffic to whatever is running on that port.

---

## Running Your Own Infrastructure

All components are open source. You can self-host the full stack:

- **agni-seeder** — run your own registry
- **agni-router** — run your own edge router
- **agni-nova** — run your own front-door proxy

Connect your components to the broader network for better distribution, or run a fully isolated private deployment.

See [agni-router-quickstart.md](agni-router-quickstart.md) for router setup.
