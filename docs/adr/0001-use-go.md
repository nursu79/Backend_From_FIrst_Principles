# ADR-0001: Use Go for Hookline

- **Status:** accepted
- **Version:** v0.0
- **Date:** 2026-10-22

## Context

Hookline is an infrastructure platform designed for high-throughput, reliable webhook ingestion and delivery. The system has two distinct workloads:
1. **High-throughput HTTP Ingestion:** Synchronously receiving events from upstream providers (e.g., Stripe, Shopify), parsing payloads, verifying auth, and persisting them with sub-20ms latency.
2. **Concurrent Webhook Delivery:** Operating thousands of concurrent outbound HTTP requests to external, third-party endpoints. These external endpoints are out of our control: they may be slow, timing out, or returning server errors.

We need a technology stack that provides:
- Lightweight concurrency to handle tens of thousands of concurrent I/O-bound outgoing HTTP requests without exhausting system memory.
- Predictable resource consumption and fast cold-start / restart times.
- Strong standard library primitives close to the operating system network stack (`net/http`, `net`, `context`).
- Low operational complexity (simple deployments, minimal runtime dependencies).

## Options considered

### 1. Node.js with TypeScript
- **Pros:**
  - Highly productive ecosystem, rich JSON manipulation primitives, and large talent pool.
  - Non-blocking asynchronous I/O via the event loop handles idle I/O connections efficiently.
- **Cons & Trade-offs:**
  - **Single-Threaded Event Loop Bottleneck:** CPU-bound tasks (such as HMAC-SHA256 signature generation or deserializing large JSON payloads) block the main event loop, causing latency spikes and dropped requests across all concurrent operations.
  - **Resource Footprint:** Each Node.js process runs Google's V8 engine, requiring 100MB–200MB of baseline memory. Scaling across CPU cores requires clustering multiple distinct processes, multiplying memory consumption.
  - **Dependency Bloat:** Heavy reliance on external npm packages with deep dependency trees increases security supply-chain risks.

### 2. Rust
- **Pros:**
  - Maximum raw execution performance, zero-cost abstractions, and guaranteed memory safety without a garbage collector.
  - Low and deterministic memory usage.
- **Cons & Trade-offs:**
  - Webhook delivery is fundamentally I/O-bound (dominated by network transit and recipient latency), meaning Rust's microsecond CPU performance gains yield negligible end-to-end advantages over Go.
  - Steep learning curve, slow compilation times, and high cognitive overhead managing asynchronous lifetimes and the borrow checker for standard web service workloads.

### 3. Go (Chosen)
- **Pros:**
  - **Goroutines & M:N Scheduler:** Goroutines begin with only ~2KB to 4KB of stack memory (compared to 1MB–8MB for operating system threads), allowing Hookline to run 50,000+ concurrent delivery workers with minimal memory overhead (~100MB). Go's runtime multiplexes these across CPU cores, preventing CPU-heavy crypto operations on one core from starving I/O on others.
  - **Single Static Binary:** Compiles directly to machine code with zero external runtime dependencies. Production deployments are atomic and can run in lightweight scratch containers.
  - **Standard Library Quality:** Go's standard library provides production-grade HTTP servers, HTTP clients, connection pooling, and context cancellation primitives out of the box, reducing external dependency risks.
- **Cons & Trade-offs:**
  - Explicit and verbose error handling (`if err != nil`) requires disciplined code organization.
  - Garbage collection introduces minor pauses, though modern Go GC pauses are sub-millisecond and acceptable for Hookline's latency profile.

## Decision

We chose **Go** as the primary implementation language for Hookline.

## Consequences

- **What gets easier:**
  - Building reliable concurrent delivery workers using native goroutines, channels, and context cancellation.
  - Deployment and operations via single static binaries with minimal Docker images.
  - Understanding low-level networking without framework abstractions obscuring the HTTP lifecycle.
- **What gets harder:**
  - Writing boilerplate error handling and serialization code compared to dynamic/expressive languages like TypeScript.
- **What we'll need to revisit and when:**
  - If CPU-intensive payload transformations or edge cryptography become a bottleneck at massive scale (e.g. v2.0+), we can evaluate isolated Rust extensions or dedicated worker services via an ADR.
