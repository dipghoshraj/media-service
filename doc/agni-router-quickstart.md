# Quick Start: Agni Router

This guide covers setting up and running the Agni Router, generating TLS certificates, and configuring `agni-config.yaml`.

## 1. Prepare `agni-config.yaml`

Ensure your config file includes the following under the `Router` section:

```yaml
Router:
  name: "agni-Router"
  router_ip: "127.0.0.1"    # IP to include in the server certificate's SAN
  rpc_port: "9000"           # gRPC port that agents connect to
  proxy_port: "8000"         # TCP proxy port for edge traffic
  certs: "agni-router/certmanger"  # Directory where certificates are stored/read from
  region: "global"
  dns: "localhost"           # DNS name to include in the server certificate's SAN
  Seeder:
    address: "localhost:8080"
    fingureprint: "<seeder-cert-fingerprint>"  # Note: field is 'fingureprint' (preserve typo)
```

Key fields used during cert generation:

| Field | Purpose |
|-------|---------|
| `name` | Common name (CN) written into the certificate |
| `router_ip` | IP SAN added to the certificate |
| `dns` | DNS SAN added to the certificate |
| `certs` | Output directory where `server.pem` / `server-key.pem` are written |

## 2. Generate Router Certificates

Certificates are generated using the `gen-certs` subcommand built into the router binary. Values are read from `agni-config.yaml` in the current working directory.

**Using `make` (recommended):**

```bash
make router-certs
```

**Using the built binary directly:**

```bash
make build-router
./bin/agni-router gen-certs
```

**Using `go run` during development:**

```bash
cd agni-router
go run . gen-certs
```

This generates `server.pem` and `server-key.pem` in the directory specified by `Router.certs`.

## 3. Start the Router

Build and run:

```bash
make build-router
./bin/agni-router
```

Or run directly during development:

```bash
cd agni-router
go run . 
```

The router starts:
- A **gRPC server** on `Router.rpc_port` (agents connect here)
- A **TCP SNI proxy** on `Router.proxy_port` (edge traffic enters here)

## 4. Verify the Router Fingerprint

After generating certificates, capture the router's certificate fingerprint to configure agents and nova. The fingerprint is logged at startup under the `[Agni Router]` prefix.

Copy the fingerprint into `agni-config.yaml` for each agent:

```yaml
Agent:
  Seeder:
    fingureprint: "<router-cert-fingerprint>"
```

## 5. Subcommand Reference

| Subcommand | Description |
|------------|-------------|
| *(none)* | Start the router server (default behaviour) |
| `gen-certs` | Generate self-signed TLS certificates from `agni-config.yaml` and exit |
