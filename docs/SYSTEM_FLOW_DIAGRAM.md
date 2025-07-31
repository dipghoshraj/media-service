# indraNet Ingress-Tunnel - System Flow Diagram

## Complete System Flow Visualization

### 1. Agent Connection Establishment Flow

```mermaid
sequenceDiagram
    participant Agent as Agent-Tunnel<br/>(Edge Device)
    participant Gateway as Gateway-Tunnel<br/>(Server)
    participant Registry as Session Registry
    participant Auth as Auth Validator

    Agent->>Gateway: WebSocket Connection Request (WSS)
    Gateway-->>Agent: WebSocket Connection Established

    Agent->>Agent: Generate HMAC Signature<br/>(token + timestamp + nonce)
    Agent->>Gateway: ConnectRequest<br/>(agent_id, token, timestamp, nonce, signature)
    
    Gateway->>Auth: Validate Token & Signature
    Auth-->>Gateway: Validation Result
    
    alt Authentication Success
        Gateway->>Registry: Register Agent Session<br/>(agent_id, connection, channels)
        Registry-->>Gateway: Session Registered
        Gateway-->>Agent: Connection Accepted
        
        par Read Loop
            Gateway->>Gateway: Start Read Loop<br/>(Listen for Agent Messages)
        and Write Loop
            Gateway->>Gateway: Start Write Loop<br/>(Send Requests to Agent)
        end
    else Authentication Failed
        Gateway-->>Agent: ConnectError
        Gateway->>Gateway: Close Connection
    end
```

### 2. HTTP Request Processing Flow

```mermaid
sequenceDiagram
    participant Client as External Client
    participant Gateway as Gateway-Tunnel<br/>(Port 8082)
    participant Registry as Session Registry
    participant InFlight as In-Flight Manager
    participant Agent as Agent-Tunnel
    participant Service as Local Service

    Client->>Gateway: HTTP Request<br/>(GET/POST myapp.domain.com/api)
    Gateway->>Gateway: Extract Agent ID<br/>from Host Header
    Gateway->>Registry: Lookup Agent Session<br/>(agent_id)
    
    alt Agent Session Found
        Registry-->>Gateway: Agent Session<br/>(connection, channels)
        Gateway->>Gateway: Generate Request ID<br/>(UUID)
        Gateway->>InFlight: Register Request<br/>(id, ResponseWriter)
        InFlight-->>Gateway: Request Registered
        
        Gateway->>Gateway: Create TunnelRequest<br/>(id, method, path, headers, body)
        
        par Send to Agent
            Gateway->>Agent: TunnelRequest via WebSocket
        and Start Streaming Goroutine
            Gateway->>InFlight: Start StreamToClient(id)
        end
        
        Agent->>Service: Forward HTTP Request<br/>(to localhost:port)
        Service-->>Agent: HTTP Response<br/>(headers + body)
        
        Agent->>Gateway: TunnelResponse<br/>(id, status, headers, initial_body)
        Gateway->>InFlight: Resolve(id, response)
        InFlight->>Client: Write Headers & Status
        
        loop Response Body Streaming
            Agent->>Gateway: TunnelData<br/>(id, chunk)
            Gateway->>InFlight: Stream(id, chunk)
            InFlight->>Client: Write Chunk
        end
        
        Agent->>Gateway: TunnelClose<br/>(id)
        Gateway->>InFlight: Close(id)
        InFlight->>InFlight: Cleanup Request
        
    else Agent Session Not Found
        Gateway-->>Client: 404 Agent Not Found
    end
```

### 3. Detailed Component Interaction Flow

```mermaid
graph TB
    subgraph "External World"
        Client[External HTTP Client]
    end
    
    subgraph "Gateway-Tunnel Server"
        HTTPServer[HTTP Server<br/>:8082]
        WSServer[WebSocket Server<br/>:8080]
        Handler[Request Handler]
        SessionMgr[Session Manager<br/>Registry]
        InFlightMgr[In-Flight Manager]
        AuthValidator[Auth Validator]
    end
    
    subgraph "Agent-Tunnel Client"
        CLI[CLI Interface<br/>Cobra Commands]
        WSClient[WebSocket Client]
        TunnelClient[Tunnel Client]
        ConnMgr[Connection Manager]
    end
    
    subgraph "Edge Device"
        LocalService[Local HTTP Service<br/>localhost:port]
    end
    
    Client -->|HTTP Request| HTTPServer
    HTTPServer --> Handler
    Handler --> SessionMgr
    Handler --> InFlightMgr
    
    CLI -->|connect command| WSClient
    WSClient -->|WebSocket| WSServer
    WSServer --> AuthValidator
    WSServer --> SessionMgr
    
    SessionMgr <-->|TunnelRequest| WSServer
    WSServer <-->|WebSocket Messages| WSClient
    WSClient --> TunnelClient
    TunnelClient --> ConnMgr
    ConnMgr -->|HTTP Request| LocalService
    LocalService -->|HTTP Response| ConnMgr
    
    InFlightMgr -->|Stream Response| HTTPServer
    HTTPServer -->|HTTP Response| Client
    
    style Gateway-Tunnel fill:#e1f5fe
    style Agent-Tunnel fill:#f3e5f5
    style External fill:#fff3e0
    style Edge fill:#e8f5e8
```

### 4. Message Flow Architecture

```mermaid
graph LR

%% Subgraph: Protocol Buffer Messages
subgraph "Protocol Buffer Messages"
    ConnectReq["ConnectRequest\n• agent_id\n• token\n• timestamp\n• nonce\n• signature"]
    TunnelReq["TunnelRequest\n• id\n• method\n• path\n• headers\n• body"]
    TunnelResp["TunnelResponse\n• id\n• status\n• headers\n• body"]
    TunnelData["TunnelData\n• id\n• chunk"]
    TunnelClose["TunnelClose\n• id"]
    ControlMsg["ControlMessage\n• type (PING/PONG)\n• payload"]
end

%% Subgraph: Message Flow
subgraph "Message Flow"
    Agent[Agent]
    Gateway[Gateway]

    Agent -->|1. ConnectRequest| Gateway
    Gateway -->|2. TunnelRequest| Agent
    Agent -->|3. TunnelResponse| Gateway
    Agent -->|4. TunnelData| Gateway
    Agent -->|5. TunnelClose| Gateway

    Agent -->|PING| Gateway
    Gateway -->|PONG| Agent
end

%% Reference arrows for message roles
ConnectReq -.->|Authentication| Agent
TunnelReq -.->|Request Forwarding| Gateway
TunnelResp -.->|Response Headers| Agent
TunnelData -.->|Response Streaming| Agent
TunnelClose -.->|Connection Cleanup| Agent
ControlMsg -.->|Health Monitoring| Agent

```

### 5. Error Handling Flow

```mermaid
graph TD
    Start[Start Operation] --> CheckAuth{Authentication<br/>Valid?}
    
    CheckAuth -->|No| AuthError[Send ConnectError<br/>Close Connection]
    CheckAuth -->|Yes| CheckSession{Agent Session<br/>Exists?}
    
    CheckSession -->|No| SessionError[Return 404<br/>Agent Not Found]
    CheckSession -->|Yes| SendRequest[Send TunnelRequest<br/>to Agent]
    
    SendRequest --> CheckTimeout{Response<br/>within 10s?}
    
    CheckTimeout -->|No| TimeoutError[Return 504<br/>Gateway Timeout<br/>Cleanup Request]
    CheckTimeout -->|Yes| ProcessResponse[Process TunnelResponse<br/>Stream to Client]
    
    ProcessResponse --> CheckStream{Streaming<br/>Complete?}
    CheckStream -->|No| StreamError[Handle Stream Error<br/>Cleanup Resources]
    CheckStream -->|Yes| Success[Request Complete<br/>Cleanup Resources]
    
    AuthError --> End[End]
    SessionError --> End
    TimeoutError --> End
    StreamError --> End
    Success --> End
    
    style AuthError fill:#ffcdd2
    style SessionError fill:#ffcdd2
    style TimeoutError fill:#ffcdd2
    style StreamError fill:#ffcdd2
    style Success fill:#c8e6c9
```

### 6. Connection Lifecycle Management

```mermaid
stateDiagram-v2
    [*] --> Disconnected
    
    Disconnected --> Connecting: agent-tunnel connect
    Connecting --> Authenticating: WebSocket Established
    
    Authenticating --> Connected: Auth Success
    Authenticating --> Disconnected: Auth Failed
    
    Connected --> Processing: Receive TunnelRequest
    Processing --> Connected: Send TunnelResponse
    
    Connected --> Heartbeat: Periodic Ping
    Heartbeat --> Connected: Pong Received
    Heartbeat --> Disconnected: Timeout
    
    Connected --> Disconnected: Connection Lost
    Processing --> Disconnected: Connection Error
    
    Disconnected --> [*]: Process Exit
```

## Key Flow Characteristics

### Performance Characteristics
- **Concurrent Processing**: Each request handled in separate goroutine
- **Streaming Support**: Chunked response streaming for large payloads
- **Connection Pooling**: WebSocket connection reuse for multiple requests
- **Request Correlation**: UUID-based request/response matching

### Reliability Features
- **Timeout Management**: 10-second request timeout
- **Connection Health**: Periodic heartbeat monitoring
- **Graceful Shutdown**: Proper resource cleanup on disconnection
- **Error Recovery**: Comprehensive error handling at each layer

### Security Features
- **HMAC Authentication**: Cryptographic request signing
- **TLS Encryption**: End-to-end encrypted communication
- **Request Isolation**: Each request processed independently
- **Session Validation**: Continuous session state monitoring

This flow diagram provides a comprehensive view of how the indraNet ingress-tunnel system operates, from initial connection establishment through request processing and error handling.
