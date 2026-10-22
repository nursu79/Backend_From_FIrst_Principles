# The Journey of a Request

This document details every hop a request takes from an external client to **Hookline** and back, identifying where latency and failure modes live at each stage.

---

## 1. Sequence Diagram

```mermaid
sequenceDiagram
    autonumber
    actor Client as External Client<br/>(e.g., Shopify / Stripe / curl)
    participant DNS as DNS Resolver<br/>(UDP/DoH)
    participant RP as Reverse Proxy / LB<br/>(e.g., NGINX / Cloudflare)
    participant GoAPI as Hookline API<br/>(cmd/api : Go Process)
    participant DB as PostgreSQL<br/>(System of Record)

    %% Step 1: DNS
    rect rgb(240, 248, 255)
    Note over Client,DNS: Phase 1: Name Resolution (Cold Start)
    Client->>+DNS: 1. Resolve domain: api.hookline.dev
    DNS-->>-Client: 2. IP Address: 140.82.121.6 (Cached via TTL)
    end

    %% Step 2 & 3: TCP & TLS
    rect rgb(255, 250, 240)
    Note over Client,RP: Phase 2: Connection Setup (Over Public Internet)
    Client->>+RP: 3. TCP SYN
    RP-->>Client: 4. TCP SYN-ACK
    Client->>RP: 5. TCP ACK (Socket Connected)
    
    Client->>RP: 6. TLS 1.3 ClientHello (Cipher suites, ALPN: h2, http/1.1)
    RP-->>Client: 7. TLS 1.3 ServerHello, Certificate & Finished
    Client->>-RP: 8. Session Keys Established (Encrypted Tunnel)
    end

    %% Step 4: HTTP Request & Reverse Proxy
    rect rgb(245, 255, 245)
    Note over Client,GoAPI: Phase 3: Application Ingestion
    Client->>+RP: 9. POST /v1/events (Encrypted HTTP/2 or HTTP/1.1)
    Note over RP: Reverse Proxy terminates TLS,<br/>buffers slow request body,<br/>forwards raw HTTP over internal LAN
    RP->>+GoAPI: 10. HTTP Request to internal Go socket (127.0.0.1:8080)
    
    Note over GoAPI: Go HTTP server routes to internal/http<br/>Decodes & validates JSON body<br/>Passes to internal/app service
    GoAPI->>+DB: 11. INSERT event & queue task (SQL Transaction)
    Note over DB: Disk I/O: Append to WAL & write table row
    DB-->>-GoAPI: 12. Commit Confirmed (id = 42)
    
    GoAPI-->>-RP: 13. HTTP 202 Accepted {"id": "evt_42", "status": "queued"}
    RP-->>-Client: 14. Encrypted HTTP Response to Client
    end
```

---

## 2. Latency & Failure Modes by Hop

| Hop | Technical Step | Typical Latency | What Can Fail? | How We Mitigate |
|---|---|---|---|---|
| **1. DNS** | Domain lookup via UDP port 53 | **5ms – 100ms** (0ms if local DNS cache hit) | DNS timeout, hijacked records (NXDOMAIN), misconfigured TTL | Client DNS caching, Anycast DNS providers (Cloudflare/Route53). |
| **2. TCP** | 3-way handshake (`SYN` → `SYN-ACK` → `ACK`) | **1 RTT** (10ms–150ms depending on geographic distance) | Packet loss, SYN flood attacks, network route changes | TCP Keep-Alive and connection pooling reuse established sockets. |
| **3. TLS** | TLS 1.3 cryptographic handshake & cert verification | **1 RTT** (~20ms–150ms) | Expired certificates, untrusted CA, cipher mismatch, clock drift | TLS Session Resumption (0-RTT), automated cert renewal (Let's Encrypt). |
| **4. Reverse Proxy** | TLS decryption, slow-client buffering, load balancing | **0.5ms – 2ms** | Proxy out of worker connections, 502 Bad Gateway if Go server is down | Reverse proxy absorbs slowloris/DDoS; health checks route traffic to healthy nodes. |
| **5. Go Process** | Router parsing, JSON decoding, input validation | **0.1ms – 1ms** (In-memory CPU) | Malformed JSON (400), payload too large (413), panic/crash | Strict payload limits (`http.MaxBytesReader`), error handling, goroutine recovery. |
| **6. Database** | PostgreSQL insert into events & job queue | **2ms – 50ms+** (Disk I/O) | Connection pool exhaustion, lock contention, slow disk, database crash | Connection pooling (`pgxpool`), database indexing, lean synchronous write transactions. |
