# indraNet Ingress-Tunnel - High-Level Design

## System Overview

The indraNet Ingress-Tunnel is a secure tunneling system designed to expose internal services to the internet through a gateway-agent architecture. It enables HTTP traffic forwarding from external clients to internal services running on edge devices through encrypted WebSocket connections.

## Architecture Components

### 1. Gateway-Tunnel (Server Component)
**Location**: `indraNet/tunnel/gateway-tunnel/`
**Purpose**: Acts as the public-facing proxy server that receives HTTP requests and forwards them to connected agents.

**Key Components**:
- **HTTP Server**: Listens on port 8082 for incoming HTTP requests
- **WebSocket Server**: Listens on port 8080 for agent connections
- **Session Manager**: Manages active agent connections and their lifecycle
- **In-Flight Manager**: Handles request-response correlation and streaming
- **Request Handler**: Processes incoming HTTP requests and routes them to appropriate agents

### 2. Agent-Tunnel (Client Component)
**Location**: `indraNet/tunnel/agent-tunnel/`
**Purpose**: Runs on edge devices, establishes secure connections to the gateway and forwards requests to local services.

**Key Components**:
- **CLI Interface**: Command-line interface using Cobra for user interactions
- **WebSocket Client**: Establishes and maintains connection to the gateway
- **Authentication Handler**: Manages HMAC-based authentication
- **Request Processor**: Handles incoming tunnel requests and forwards them locally
- **Connection Manager**: Manages local TCP connections to target services

## Communication Protocol

### Protocol Buffer Schema
```proto
message Envelope {
  oneof message {
    ConnectRequest connect = 1;
    TunnelRequest request = 2;
    TunnelResponse response = 3;
    ControlMessage control = 4;
    TunnelData data = 5;
    TunnelClose close = 6;
    ConnectError connect_error = 7;
  }
}
```

### Message Types:
1. **ConnectRequest**: Agent authentication and connection establishment
2. **TunnelRequest**: HTTP request forwarding from gateway to agent
3. **TunnelResponse**: HTTP response headers and initial body
4. **TunnelData**: Streaming response body chunks
5. **TunnelClose**: Connection termination signals
6. **ControlMessage**: Heartbeat and control signals

## System Flow

### Connection Establishment Flow
```
1. Agent starts with: agent-tunnel connect --gateway wss://gateway --token <token> --secret <secret> --id <agent-id>
2. Agent creates WebSocket connection to gateway
3. Agent sends ConnectRequest with HMAC signature
4. Gateway validates token and signature
5. Gateway registers agent session
6. Bidirectional communication channel established
```

### Request Processing Flow
```
1. External HTTP request arrives at Gateway (port 8082)
2. Gateway extracts agent ID from Host header/domain
3. Gateway looks up agent session in Session Registry
4. Gateway creates TunnelRequest with unique ID
5. Gateway registers request in In-Flight Manager
6. Gateway sends TunnelRequest to Agent via WebSocket
7. Agent receives request and forwards to local service
8. Agent sends TunnelResponse with headers and status
9. Agent streams response body via TunnelData messages
10. Gateway streams response back to original HTTP client
11. Gateway cleans up in-flight request when complete
```

## Key Features

### Security
- **HMAC Authentication**: Token-based authentication with signature verification
- **WebSocket Encryption**: Secure WebSocket (WSS) connections
- **Request Isolation**: Each request gets unique ID for correlation

### Scalability
- **Concurrent Handling**: Goroutine-based concurrent request processing
- **Session Management**: Efficient agent session tracking
- **Buffer Management**: Buffered channels for request/response handling

### Reliability
- **Heartbeat Mechanism**: Connection health monitoring
- **Timeout Handling**: Request timeout management (10 seconds)
- **Graceful Shutdown**: Proper cleanup of connections and resources
- **Error Handling**: Comprehensive error reporting and recovery

## Data Flow Architecture

```
┌─────────────────┐    HTTP Request    ┌──────────────────┐
│   External      │ ─────────────────→ │   Gateway        │
│   Client        │                    │   Tunnel         │
│                 │ ←───────────────── │   (Port 8082)    │
└─────────────────┘    HTTP Response   └──────────────────┘
                                               │
                                               │ WebSocket
                                               │ (Port 8080)
                                               ▼
                                       ┌──────────────────┐
                                       │   Agent          │
                                       │   Tunnel         │
                                       │   (Edge Device)  │
                                       └──────────────────┘
                                               │
                                               │ Local HTTP
                                               ▼
                                       ┌──────────────────┐
                                       │   Local          │
                                       │   Service        │
                                       │   (Target App)   │
                                       └──────────────────┘
```

## Technical Details

### Port Configuration
- **Gateway HTTP Server**: 8082 (for incoming requests)
- **Gateway WebSocket Server**: 8080 (for agent connections)
- **Agent**: Configurable local service ports

### Session Management
- **Registry Pattern**: Central session registry for agent tracking
- **Connection State**: Active connection monitoring with LastSeen timestamps
- **Automatic Cleanup**: Session cleanup on disconnection

### Request Correlation
- **UUID-based IDs**: Unique request identification
- **In-Flight Tracking**: Active request state management
- **Response Streaming**: Chunked response handling for large payloads

### Error Handling
- **Connection Errors**: WebSocket connection failure handling
- **Authentication Errors**: Invalid token/signature responses
- **Timeout Errors**: Request timeout management
- **Service Errors**: Target service unavailability handling

## Deployment Architecture

### Gateway Deployment
- Can be deployed as containerized service (Dockerfile included)
- Requires persistent connection handling
- Stateful service managing agent sessions

### Agent Deployment
- CLI-based deployment on edge devices
- Lightweight client with minimal dependencies
- Can be deployed as system service/daemon

## Security Considerations

1. **Authentication**: HMAC-SHA256 based token authentication
2. **Encryption**: TLS/WSS for all communications
3. **Authorization**: Agent ID based access control
4. **Input Validation**: Request sanitization and validation
5. **Rate Limiting**: Channel-based request queuing

## Monitoring and Observability

- **Connection Logging**: Agent connect/disconnect events
- **Request Logging**: HTTP request/response logging
- **Error Logging**: Comprehensive error tracking
- **Health Checks**: Connection health monitoring via heartbeats

This architecture provides a robust, secure, and scalable solution for exposing internal services through a tunneling mechanism, enabling secure remote access to edge device services.
