# Agnistack — Privacy-First fabric for exposing private servers to the internet

Modern infrastructure is becoming increasingly centralized. Developers depend on third-party platforms to expose services, route traffic, and manage networking — often at the cost of privacy, control, and resilience

Agnistack is a privacy-first, decentralized application deployment network. It is built to expose **private servers** to the internet through a distributed routing network — without TLS termination, without data inspection, and without giving up control of your certificates or domain.

You bring your own domain, your own SSL certs, and your own server. Agnistack routes the raw TCP stream directly to it. Nothing more.

AgniStack explores a different direction:

- **Distributed infrastructure** — no single point of failure or control
- **Community-operated networking** — anyone can run a router, seeder, or proxy node
- **Privacy-first connectivity** — raw TCP relay, no TLS termination, no payload inspection
- **Zero-trust identity** — certificate fingerprint pinning, not CA-chain trust
- **Edge-native** — built to work in restricted, firewalled, or censored environments

The goal is not just tunneling. The goal is foundational infrastructure for the next generation of private, distributed systems.

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

## Security

AgniStack uses a zero-trust-inspired model with no dependency on certificate authorities.

- **Outbound-only connections.** The agent always initiates. No inbound ports are opened on your server.
- **Certificate fingerprint identity.** Each agent generates a self-signed TLS certificate. The SHA-256 fingerprint is the agent's identity — registered with the seeder and verified by the router on every connection.
- **TLS 1.3 only.** All gRPC connections enforce `MinVersion: tls.VersionTLS13`. A custom `VerifyPeerCertificate` callback rejects any peer whose fingerprint doesn't match the seeder's response, even with `InsecureSkipVerify: true` set for the dial.
- **No CA chain.** No certificate authority to trust, rotate, or compromise.
- **No payload inspection.** Traffic is forwarded as raw bytes — AgniStack has no visibility into connection content.

---



## Getting Started with agni-agent

### Prerequisites

- A domain with DNS access
- A private server or local machine running your application
- `agni-agent` binary ([download](#build-from-source) or build from source)

---

### Step 1 — Discover Available Seeders

```bash
agni-agent scan
```

Prints a table of seeders (address, region, fingerprint). Pick one for your config.

---

### Step 2 — Configure `agni-config.yaml`

```yaml
version: v1

Agent:
  name: "my-agent"
  domain: "myapp.example.com"   # Your domain — used for SNI routing
  forward: 443                  # Local port your app listens on
  host: "localhost"
  region: "global"
  certs: "./"                   # Directory containing client.pem + client-key.pem
  Seeder:
    address: "45.130.164.217:8080"
    fingureprint: "<seeder-fingerprint>"
```

See [agni-agent-quickstart.md](agni-agent-quickstart.md) for a full field reference.

---

### Step 3 — Generate Agent Credentials

```bash
agni-agent gen-creds --dns myapp.example.com --name client
```

Creates `client.pem` and `client-key.pem` in the current directory.

---

### Step 4 — Connect

Start your application, then run the agent on the same machine:

```bash
agni-agent connect
```

The agent registers with the seeder, discovers the assigned router, and opens a persistent gRPC tunnel. Once connected, add a DNS `A` record for your domain pointing to the nova address the agent prints.

> **Local development:** set `host: localhost` and `forward` to your dev port. The tunnel behaves identically.

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

## Current Capabilities

- Secure application exposure through outbound-only tunnels
- SNI-based distributed ingress routing
- Persistent gRPC bidirectional streams between agent and router
- Certificate fingerprint identity — no CA chain
- TLS 1.3 enforcement across all connections
- Decentralized seeder discovery
- Multi-platform agent (`linux`, `darwin`, `windows`)
- Self-hostable full stack

---

## In Progress

- Global routing optimization across multiple regions
- Observability layer (connection metrics, tunnel health)
- Access control and policy enforcement
- Multi-region failover
- Dynamic service discovery

---

## Future Direction

AgniStack is the foundational infrastructure layer for future privacy-first systems. Long-term directions include:

- Zero-trust communication infrastructure
- Community-backed anonymous networking
- Decentralized edge-native connectivity protocols
- Privacy-first routing with no central registry dependency

---

## Running Your Own Infrastructure

All components are open source. You can self-host the full stack:

- **agni-seeder** — run your own registry
- **agni-router** — run your own edge router
- **agni-nova** — run your own front-door proxy

Connect your components to the broader network for better distribution, or run a fully isolated private deployment.

See [agni-router-quickstart.md](agni-router-quickstart.md) for router setup.
