# Curriculum map

Source: Sriniously, **Backend from First Principles** (YouTube playlist). Titles are as listed in the playlist; check numbering against the playlist itself, since a few videos may have been added or renamed since this map was written.

| # | Video | Hookline version | What it becomes in Hookline |
|---|---|---|---|
| 1 | Roadmap for backend from first principles | v0.0 | Project orientation |
| 2 | Walk the path of a true backend engineer | v0.0 | Mindset and method |
| 3 | What is a backend, how do they work and why do we need them? | v0.0 | Request-journey diagram |
| 4 | Benefits of learning backend engineering from first principles | v0.0 | Why Go + standard library first |
| 5 | Understanding HTTP for backend engineers | v0.1 | Raw TCP HTTP server, then `net/http` |
| 6 | What is routing in backend? | v0.2 | `/v1` route table with `ServeMux` |
| 7 | Serialization and deserialization | v0.3 | Strict JSON decoding, raw payload storage |
| 8 | Authentication and authorization | v0.4 | API keys, tenant isolation |
| 9 | Validations and transformations | v0.5 | Boundary validation, field errors |
| 10 | Controllers, services, repositories, middlewares and request context | v0.6 | Layered refactor, middleware chain |
| 11 | Complete REST API design | v0.7 | Resource model, pagination, idempotency keys |
| 12 | Mastering databases with Postgres | v0.8 | Schema, migrations, indexes |
| 13 | Caching, the secret behind it all | v0.9 | API key and endpoint caches |
| 14 | Task queues and background jobs | v0.10 | Postgres queue + delivery workers |
| 15 | Full-text search using Elasticsearch | v0.11 | Message search |
| 16 | Error handling and building fault-tolerant systems | v0.12 | Retries, backoff, circuit breakers |
| 17 | Production-grade configuration management | v0.13 | Typed, validated config |
| 18 | Logging, monitoring and observability | v0.14 | slog, metrics, traces, SLO |
| 19 | Graceful shutdown | v0.15 | Draining API and workers |
| 20 | Backend security | v0.16 → **v1.0** | SSRF guard, rate limits, secrets |
| 21 | Backend scaling and performance engineering | v1.1 | Load tests, profiling, replicas |
| 22 | Concurrency and parallelism: IO bound vs CPU bound | v1.2 | Fair worker pool, ordering |
| 23–24 | Object storage | v1.3 | Large payloads, archival |
| 25 | Real-time backends | v1.4 | SSE live delivery feed |
| 26 | Testing for backend engineers | v1.5 | Test strategy, chaos test (tests start in v0.1) |
| 27 | The twelve-factor app | v1.6 | Audit, Docker, CI/CD, deploy |
| 28 | OpenAPI | v1.7 | Spec, docs, generated client |
| 29 | Webhooks | v1.8 | Standard Webhooks signing, receiver library |
| 30 | WebSockets | **v2.0** | WebSocket transport + `hookline listen` CLI |

## Adding other sources

The playlist is the only source for now. When a concept needs more depth, add the extra resource here next to the version it supports, with one line on why it was added.

| Version | Extra resource | Why |
|---|---|---|
| | | |
