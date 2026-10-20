# Hookline Roadmap

From an empty folder to a production-grade webhook platform in 26 versions. Every version follows the order of Sriniously's **Backend from First Principles** playlist, so the video is always the foundation and the code is where I apply it.

## The story that runs through every version

Hookline is a **post office**.

Customers (other companies' backends) drop off letters (events). The post office sorts them, records them, and sends couriers to deliver them to the addresses each customer registered (their webhook endpoints). Some addresses are closed, some are slow, some are fake. The post office has to deliver every letter at least once, prove it tried, and never let one customer's flood of mail delay everyone else's.

Each version adds one room, rule or role to this post office. The analogy comes first so the concept has somewhere to land, then the technical version replaces it.

## How each version works

Every version has the same shape, and the lecturer agent (see [`AGENTS.md`](./AGENTS.md)) walks through it in this order:

| Step | What happens |
|---|---|
| 1. Watch | Watch the listed video(s). Take notes in your own words. |
| 2. Story | The lecturer tells the post-office version of the concept. |
| 3. Pain | Why the current version of Hookline can't cope. No pain, no new version. |
| 4. Principle | The first principle underneath, tied back to the video. |
| 5. Design | You propose the design first. The lecturer questions it. Then you agree on one. |
| 6. Build | You write the code in small steps. The lecturer reviews; code from the lecturer only comes after you can explain why it's needed. |
| 7. Break it | Run the experiment that proves the version works, or shows exactly how it fails. |
| 8. Explain back | Explain the concept to the lecturer without notes. |
| 9. Record | Learning-log entry, ADR(s), CHANGELOG line, tests, then `git tag vX.Y.0`. |

### Versioning rules

- **Tags** follow SemVer: `v0.x` while the platform is being shaped, `v1.0.0` when it's safe to run in production, `v2.0.0` when it gains a second delivery transport.
- **Branches:** each version is built on a branch named `version/vX.Y-short-name` and merged into `main` through a pull request you review like a teammate's.
- **API versioning:** the public HTTP API lives under `/v1` from v0.2 onwards. Breaking it later requires an ADR.
- **Testing is a thread, not a chapter.** Video 26 is late in the playlist, but from v0.1 every version adds tests for what it built. v1.5 is where the testing strategy is formalized.

### Phase gates

At the end of each phase, stop and the supervisor agent reviews the work before continuing. The gate checks that the design still holds, the evidence is complete, and nothing was skipped.

| Gate | After | Question it answers |
|---|---|---|
| A | v0.0 | Do I understand the journey of a request and what Hookline is? |
| B | v0.6 | Is the code shaped so it can grow (layers, middleware, context)? |
| C | v0.10 | Does Hookline actually deliver webhooks end to end? |
| D | v1.0 | Is it safe to run in production? (errors, config, observability, shutdown, security) |
| E | v1.4 | Does it hold up under load, concurrency and real-time use? |
| F | v2.0 | Portfolio review: can I defend every design decision in an interview? |

---

## Phase 1: Foundations (Gate A)

### v0.0 — Orientation
**Videos:** 1. Roadmap for backend from first principles · 2. Walk the path of a true backend engineer · 3. What is a backend, how do they work and why do we need them? · 4. Benefits of learning backend engineering from first principles

**Story:** Before opening a post office, you walk the town. Where do letters come from? Who sends them, who receives them, what roads do they travel?

**The pain:** None yet. You can't feel the pain of a system you can't picture.

**Principles to grasp**
- A backend is a long-running process that listens on a port, receives requests, does work, and returns responses.
- The request journey: browser or client → DNS → TCP/TLS → reverse proxy / load balancer → your process → database and other services → back.
- Frameworks change; the concepts underneath (HTTP, state, concurrency, failure) don't.

**Build**
- Repo scaffold: Go module, `cmd/api/main.go` that prints a startup line, `Makefile` or `justfile`, `.gitignore`, `docker-compose.yml` placeholder.
- Draw (by hand or Excalidraw) the journey of one request to Hookline and commit the image to `docs/diagrams/`.
- Write a one-page "What is Hookline and who uses it?" in `docs/learning-log/v0.0.md`.

**Break it:** `curl -v https://api.github.com` and annotate every line of output: which part is DNS, TCP, TLS, HTTP?

**Done when**
- [ ] I can explain the request journey out loud in under 2 minutes.
- [ ] Scaffold builds with `go build ./...`.
- [ ] `v0.0.0` tag pushed.

**Evidence:** ADR-0001 "Use Go for Hookline", ADR-0002 "Record architecture decisions".

---

### v0.1 — Speak HTTP
**Videos:** 5. Understanding HTTP for backend engineers

**Story:** Every letter must follow a strict format: envelope (start line), stamps and labels (headers), a blank line, then the letter itself (body). A clerk who can't read the format can't do anything else.

**The pain:** Hookline can't talk to anyone yet.

**Principles to grasp**
- HTTP is text over TCP (for 1.1): request line, headers, blank line, body.
- Methods and their meaning (safe, idempotent), status code families, important headers (`Content-Type`, `Content-Length`, `Authorization`, `Cache-Control`), keep-alive, HTTP/1.1 vs 2 vs 3, CORS and preflight.
- Statelessness: each request carries everything needed to handle it.

**Build**
- Part 1: a raw TCP server with `net.Listen` that reads bytes, parses the request line and headers by hand, and writes a valid HTTP response. Accept `POST /events` and print the body.
- Part 2: replace it with `net/http`. Compare the two and write what the standard library now does for you.
- `GET /healthz` returning `200`.

**Break it:** Send a malformed request with `nc` (missing blank line, wrong `Content-Length`). Watch what your raw server does vs what `net/http` does.

**Done when**
- [ ] Raw server and `net/http` server both accept an event.
- [ ] I can explain what `Content-Length` and keep-alive do and why they matter.
- [ ] First test: `httptest` checks `/healthz`.

**Evidence:** learning log with an annotated raw HTTP request and response.

---

### v0.2 — Find the right door
**Videos:** 6. What is routing in backend? How requests find their way home

**Story:** The sorting room. Letters arrive in one pile; the sorter reads the address and drops each into the right pigeonhole.

**The pain:** v0.1 has one handler. Hookline needs applications, endpoints, and events, each with their own actions.

**Principles to grasp**
- A route is method + path. Static vs dynamic segments, path params vs query params.
- Route matching order and precedence, 404 vs 405, trailing slashes.
- Nesting and grouping (`/v1/apps/{app_id}/endpoints`), API versioning in the path.

**Build**
- Use Go's standard `http.ServeMux` (method + wildcard patterns) on purpose, so you see routing without a framework.
- Routes: `POST/GET /v1/apps`, `GET /v1/apps/{app_id}`, `POST/GET /v1/apps/{app_id}/endpoints`, `POST /v1/apps/{app_id}/events`.
- In-memory maps for storage.

**Break it:** Register two conflicting patterns and read the panic. Request a known path with the wrong method and check you get `405` with an `Allow` header.

**Done when**
- [ ] Route table documented in the learning log.
- [ ] Table-driven tests for every route, including 404 and 405.

---

### v0.3 — Speak the same language
**Videos:** 7. Serialization and deserialization for backend engineers

**Story:** Letters arrive in many languages. A translator turns them into the post office's internal language, and back again for replies. A sloppy translator invents words or drops sentences.

**The pain:** Hookline is passing raw bytes around. It needs structured data in and out.

**Principles to grasp**
- Serialization turns in-memory structures into bytes for the wire; deserialization goes the other way. JSON is text; formats like Protobuf are binary.
- Field naming, `omitempty`, null vs missing vs zero value, timestamps (RFC 3339, UTC), big numbers, raw payloads that must be stored untouched.
- Strict decoding: reject unknown fields and trailing data. Limit body size.

**Build**
- Request/response structs for apps, endpoints and events. Event `payload` stored as `json.RawMessage`, never reinterpreted.
- A `decodeJSON` helper with `DisallowUnknownFields`, `http.MaxBytesReader`, and clear error messages.
- Consistent JSON error shape `{ "error": { "code": "...", "message": "..." } }`.

**Break it:** Send a 10 MB body, a body with an extra field, a body with two JSON objects, and `"amount": 9007199254740993`. Record what happens to each.

**Done when**
- [ ] Every bad input above returns a clear `4xx`, never a `500`.
- [ ] Round-trip tests: a payload goes in and comes out byte-identical.

---

### v0.4 — Who's knocking?
**Videos:** 8. Authentication and authorization for backend engineers

**Story:** Two questions at the door: "Who are you?" (show your ID card: authentication) and "Are you allowed into this room?" (check the access list: authorization). A customer's ID gets them into their own mailroom, never someone else's.

**The pain:** Anyone can create apps and send events as anyone else.

**Principles to grasp**
- Authentication vs authorization. Stateful sessions vs stateless tokens (JWT), API keys, OAuth at a high level.
- Never store secrets in plain text: hash API keys (fast hash with a high-entropy key is fine; passwords need a slow hash like argon2id/bcrypt).
- Multi-tenancy: every resource belongs to an organization, and every query is scoped to it.
- 401 vs 403. Timing attacks and constant-time comparison.

**Build**
- Organizations and API keys: format `hk_live_<prefix>_<secret>`, store prefix + hash, show the full key once.
- An `Authorization: Bearer` check on all `/v1` routes, resolving the key to an organization.
- Authorization rule: an org can only see its own apps, endpoints and events.
- Still in memory.

**Break it:** Use org A's key to read org B's app by guessing its ID. It must return `404` (not `403`) so IDs don't leak. Explain why.

**Done when**
- [ ] Tests prove tenant isolation for every route.
- [ ] ADR-0003 "API keys for machine clients; sessions for the dashboard later".

---

### v0.5 — Trust, but verify
**Videos:** 9. Validations and transformations for backend engineers

**Story:** Customs inspection. Before a parcel enters the building it's checked: right size, real address, nothing forbidden. Some things are tidied on the way in (address in a standard format).

**The pain:** Hookline accepts endpoints with URL `"banana"` and events with empty types.

**Principles to grasp**
- Syntactic validation (shape, types, formats) vs semantic validation (business rules) vs transformation (normalize, trim, lowercase).
- Validate at the boundary, once. Fail fast with every field error, not just the first.
- `400` vs `422` and a consistent field-error format.

**Build**
- Validation for every request: endpoint URL must be absolute `https` (allow `http` only via config later), event type pattern like `invoice.paid`, length limits, required fields.
- Transformations: trim strings, lowercase event types, normalize URLs.
- Error response lists every invalid field.

**Break it:** Write a table of 20 bad inputs and their expected errors; turn it into a test.

**Done when**
- [ ] No handler touches unvalidated input.

---

### v0.6 — Everyone has one job
**Videos:** 10. What are controllers, services, repositories, middlewares and request context?

**Story:** The post office grows departments: the front counter only talks to customers (handlers), the back office applies the rules (services), the archive clerk only stores and fetches (repositories). Every letter passes the same checkpoints on the way in (middleware), and carries a routing slip with its tracking number and owner (request context).

**The pain:** Handlers now mix HTTP parsing, auth, validation, business rules and storage. Changing storage would mean rewriting everything.

**Principles to grasp**
- Separation of concerns: handler (HTTP in/out), service (business logic, no HTTP types), repository (data access behind an interface).
- Middleware as a chain of wrappers: request ID, structured access log, panic recovery, auth, timeouts. Order matters.
- `context.Context` carries deadlines, cancellation, and request-scoped values (request ID, org ID). It is not a bag for dependencies.
- Dependency injection by constructor; no globals.

**Build**
- Refactor into `internal/http` (handlers, middleware), `internal/app` (services), `internal/store` (repository interfaces + in-memory implementation).
- Middleware chain: recover → request ID → logger → auth → handler.
- Unit-test services with a fake repository.

**Break it:** Add a handler that panics. Confirm recovery middleware returns `500` with the request ID and the server keeps running.

**Done when**
- [ ] No `net/http` import inside `internal/app`.
- [ ] **Gate B review** by the supervisor.

---

## Phase 2: A real product (Gate C)

### v0.7 — The service catalogue
**Videos:** 11. Complete REST API design

**Story:** The post office prints a catalogue of services with fixed names, prices and rules, so every customer can use it without asking the clerk.

**The pain:** Routes grew ad hoc. List endpoints return everything at once. Retrying a request can create duplicate events.

**Principles to grasp**
- Resources and nouns, collections vs items, sub-resources, correct use of methods and status codes.
- Pagination (cursor vs offset and why cursor wins here), filtering, sorting.
- Idempotency: `Idempotency-Key` header on `POST /events` so client retries never duplicate.
- `PATCH` vs `PUT`, partial updates, consistent error format (consider RFC 9457 problem details).

**Build**
- Final resource model: `apps`, `endpoints` (url, subscribed event types, enabled), `event-types`, `messages` (an accepted event), `attempts` (each delivery try).
- Cursor pagination on every list. Filters: `?event_type=`, `?status=`, `?before=/after=`.
- Idempotency keys stored per org for 24h.
- Hand-written `docs/api.md` (it becomes OpenAPI in v1.7).

**Break it:** Send the same event with the same idempotency key 50 times in parallel. Exactly one message must exist.

---

### v0.8 — Memory that survives a restart
**Videos:** 12. Mastering databases with Postgres

**Story:** Until now the post office kept records on a whiteboard. Tonight the whiteboard gets wiped. Time for a fireproof archive with indexed filing cabinets.

**The pain:** Restart the process and every app, key and event is gone.

**Principles to grasp**
- Relational modeling, primary/foreign keys, constraints as the last line of defense.
- Migrations as versioned code. Transactions and isolation levels. Indexes and `EXPLAIN ANALYZE`.
- Connection pools. Why `SELECT *` and N+1 queries hurt.

**Build**
- Postgres via docker-compose. Migrations with `goose` (or `golang-migrate`). Queries with `pgx` + `sqlc`.
- Postgres implementations of the repository interfaces from v0.6. The services do not change: that is the payoff of v0.6.
- Unique constraint backs the idempotency rule from v0.7.
- Integration tests against a real Postgres (testcontainers).

**Break it:** Insert 1 million messages, run the "list messages for app, newest first" query with and without an index, and record both `EXPLAIN ANALYZE` outputs.

**Done when**
- [ ] ADR-0004 "Postgres as the system of record".

---

### v0.9 — Don't ask twice
**Videos:** 13. Caching, the secret behind it all

**Story:** The clerk keeps sticky notes of the questions asked a hundred times a day ("is this ID card valid?"), instead of walking to the archive each time. The danger: an out-of-date sticky note.

**The pain:** Every request hits Postgres twice before doing any work (API key lookup, endpoint lookup). Measure it first.

**Principles to grasp**
- Cache levels (in-process, distributed, CDN/HTTP). Cache-aside, write-through, TTLs.
- Invalidation, stampede (thundering herd) and how to prevent it (singleflight, jitter).
- Cache only what you've measured. A cache is a copy, never the source of truth.

**Build**
- Benchmark first (simple `hey` or `k6` script), record p50/p95.
- In-process LRU for API key → org lookups with short TTL and `singleflight`.
- Redis (or Valkey) cache for endpoint configs, invalidated on update.
- Benchmark again; put both numbers in the learning log.

**Break it:** Revoke an API key and measure how long it still works. Decide if that window is acceptable and record the decision in an ADR.

---

### v0.10 — Do it later
**Videos:** 14. Task queues and background jobs

**Story:** The counter clerk stops walking letters to their destinations. They drop each letter in the outbox tray, hand the customer a receipt, and a team of couriers takes letters from the tray in the back.

**The pain:** If Hookline delivered inside the request, a slow customer endpoint would make the sender wait, and a crash would lose the delivery.

**Principles to grasp**
- Producer/consumer, at-least-once vs at-most-once vs "exactly once" (and why exactly-once delivery isn't possible over a network).
- Visibility timeouts / leases, acknowledgements, job states.
- The transactional outbox: write the message and its delivery jobs in the same transaction.

**Build**
- `POST /events` → in one transaction: insert message, insert one delivery job per subscribed endpoint. Return `202 Accepted`.
- A Postgres-backed queue using `SELECT ... FOR UPDATE SKIP LOCKED`, built by hand first to understand it.
- `cmd/worker`: claims jobs, does a plain HTTP POST of the payload to the endpoint with a timeout, records an `attempt` row.
- A tiny "receiver" test server (`cmd/sink`) to see deliveries arrive.

**Break it:** `kill -9` the worker mid-delivery. Prove the job is retried after its lease expires and no message is lost.

**Done when**
- [ ] End-to-end demo: create app → endpoint → send event → see it arrive at the sink → see the attempt in the API.
- [ ] ADR-0005 "Postgres queue now; revisit dedicated broker at scale".
- [ ] **Gate C review.** Hookline is now a product.

---

## Phase 3: Production-ready (Gate D → v1.0)

### v0.11 — Find the needle
**Videos:** 15. Full-text search using Elasticsearch

**Story:** The archive has millions of letters. A customer asks "find every letter that mentions order 8812 that failed last week." The filing cabinets sort by date, not by contents. You need a card catalogue.

**The pain:** Customers debugging deliveries need to search payloads and failure reasons. `LIKE '%8812%'` on a big table is a full scan.

**Principles to grasp**
- Inverted indexes, tokenization, analyzers, relevance scoring.
- Search index as a secondary, eventually consistent copy of the data. Keeping it in sync (outbox or change data capture).
- When Postgres full-text search is enough and when it isn't.

**Build**
- Step 1: Postgres `tsvector` + GIN index on messages. Measure.
- Step 2: Elasticsearch (or OpenSearch) index fed from an outbox table by an indexer worker.
- `GET /v1/apps/{app_id}/messages/search?q=` with filters.

**Break it:** Stop the indexer, send events, restart it. Confirm it catches up and nothing is skipped or duplicated.

---

### v0.12 — When things break
**Videos:** 16. Error handling and building fault-tolerant systems

**Story:** Roads close, addresses move, a recipient refuses the letter. A good post office has rules: try again later, wait longer each time, stop sending couriers to an address that's been closed all week, and keep a "return to sender" shelf.

**The pain:** Hookline tries once. Real endpoints fail constantly.

**Principles to grasp**
- Error classes: client errors, transient vs permanent failures, bugs. Wrapping errors with context in Go.
- Retries with exponential backoff and jitter. Timeouts on every outbound call. Retry budgets.
- Circuit breakers, bulkheads, dead-letter queues. Failing gracefully.

**Build**
- Retry schedule (for example 5s, 30s, 2m, 15m, 1h, 6h, 24h) with jitter; status per message: `pending → succeeded | retrying → failed`.
- Which responses retry (5xx, 429, timeouts, connection errors) and which don't (most 4xx). Respect `Retry-After`.
- Per-endpoint circuit breaker; auto-disable an endpoint after N days of failure, with a notification event.
- Manual "retry now" and "replay message" API endpoints.

**Break it:** Point an endpoint at a sink that fails 70% of the time and one that hangs forever. Show the retry timeline and that the hanging one can't block workers.

---

### v0.13 — Same code, different worlds
**Videos:** 17. Production-grade configuration management

**Story:** The same post office manual is used in every town, but each town has its own opening hours and depot address. The rules are written down, checked on opening day, and keys to the safe are never written in the manual.

**The pain:** Ports, DB URLs, retry schedules and timeouts are hard-coded.

**Principles to grasp**
- Config vs code vs secrets. Sources and precedence (defaults → file → env → flags).
- Validate all config at startup and crash early with a clear message.
- Secret handling: never logged, never committed, rotated.

**Build**
- A typed `config` package loaded once in `main`, passed down explicitly.
- `.env.example`, per-environment overrides, a startup log line that prints config with secrets redacted.

**Break it:** Start with a missing DB URL and an invalid timeout. It must refuse to start and say exactly why.

---

### v0.14 — Turn the lights on
**Videos:** 18. Logging, monitoring and observability

**Story:** Every letter gets a tracking number. Every room has a counter on the wall (letters in, letters out, letters stuck). When a customer says "my letter never arrived," you can follow its path room by room.

**The pain:** When a delivery fails in production, there's no way to know why without guessing.

**Principles to grasp**
- Logs (events), metrics (numbers over time), traces (one request's path across services). What each is for.
- Structured logging with `slog`, log levels, correlation IDs.
- RED and USE methods. SLIs, SLOs, alerting on symptoms not causes.

**Build**
- `slog` JSON logs with request ID and org ID on every line.
- Prometheus metrics: HTTP RED metrics, queue depth, delivery latency, success rate per endpoint.
- OpenTelemetry tracing that follows one event from `POST /events` through the queue into the worker and out to the endpoint.
- Grafana dashboard JSON committed to the repo. One SLO: "99% of messages get a first delivery attempt within 5 seconds."

**Break it:** Inject a 3s delay into the DB. Find it using only the dashboard and a trace.

---

### v0.15 — Leave without breaking anything
**Videos:** 19. Graceful shutdown

**Story:** Closing time. You lock the front door to new customers, let the ones inside finish, wait for couriers already out on the road to come back, then turn off the lights. You don't just pull the fire alarm.

**The pain:** Every deploy kills workers mid-delivery and drops in-flight HTTP requests.

**Principles to grasp**
- Process signals (`SIGTERM`, `SIGINT`, `SIGKILL`), the orchestrator's grace period.
- Stop accepting new work → drain in-flight work → close resources in reverse order.
- Readiness vs liveness checks.

**Build**
- API: `signal.NotifyContext` → readiness goes false → `server.Shutdown(ctx)` with a deadline → close DB pool.
- Worker: stop claiming jobs, finish or release in-flight jobs before the deadline.
- `/readyz` and `/livez`.

**Break it:** Run a load test and restart the API and worker in the middle. Zero failed client requests, zero lost deliveries.

---

### v0.16 — Lock the doors
**Videos:** 20. Backend security: everything you need to know

**Story:** Thieves, forgers and pranksters show up. Someone registers the address of the post office's own safe as their "delivery address." Someone floods the counter. Someone tries to read other customers' mail.

**The pain:** Hookline sends HTTP requests to any URL a customer types. That makes it a perfect tool for attacking internal networks (SSRF).

**Principles to grasp**
- Threat modeling (STRIDE at a basic level), OWASP API Top 10.
- SSRF and how to block it properly (resolve DNS, then block private/loopback/link-local ranges at connect time, not by string matching).
- Rate limiting (token bucket), secrets encryption at rest, least privilege, audit logs, dependency scanning.

**Build**
- A custom `net.Dialer` `Control` hook that blocks private IP ranges after DNS resolution.
- Per-org rate limits on the API (Redis token bucket) with `429` and `Retry-After`.
- Endpoint signing secrets encrypted at rest. API key rotation. Audit log table.
- `govulncheck` and `gosec` in CI. A threat model in `docs/security.md`.

**Break it:** Try to register `http://169.254.169.254/`, `http://localhost:5432`, and a DNS name that resolves to `127.0.0.1`. All must be blocked.

**Done when**
- [ ] **Gate D review.** Tag `v1.0.0`: Hookline is safe to run in production.

---

## Phase 4: Scale and real time (Gate E)

### v1.1 — More traffic than one building can handle
**Videos:** 21. Backend scaling and performance engineering

**Story:** The town grows into a city. One post office can't keep up. You open more branches, add more couriers, and find which room is the real bottleneck before adding staff.

**The pain:** You don't know Hookline's limits. Guessing isn't engineering.

**Principles to grasp**
- Measure, find the bottleneck, fix it, measure again. Latency vs throughput, percentiles.
- Vertical vs horizontal scaling, stateless services, DB connection limits, read replicas, partitioning.
- Profiling (`pprof`), load testing (`k6`).

**Build**
- A `k6` script for event ingestion and a load profile for delivery.
- Run 2+ API and worker replicas behind a load balancer (nginx or Caddy) in docker-compose.
- Profile, fix the top two bottlenecks, and write a performance report in `docs/perf/v1.1.md` with before/after numbers.

---

### v1.2 — Many couriers, one truth
**Videos:** 22. Concurrency and parallelism: IO bound vs CPU bound

**Story:** Fifty couriers share one tray. Two must never grab the same letter. One customer who sends a million letters must not make everyone else's letters wait. Some couriers mostly wait at doors (IO); the clerk who seals envelopes is always busy (CPU).

**The pain:** One big customer's backlog delays every other customer (noisy neighbor). Some endpoints need events in order.

**Principles to grasp**
- Concurrency vs parallelism, goroutines and the Go scheduler, IO-bound vs CPU-bound work.
- Worker pools, bounded concurrency, backpressure. Race conditions and the race detector.
- Fairness and ordering guarantees (per-endpoint FIFO is a choice with a cost).

**Build**
- Bounded worker pool with per-endpoint concurrency limits and fair scheduling across orgs.
- Optional per-endpoint ordered delivery mode.
- Run all tests with `-race`.

**Break it:** One org sends 500k events while another sends 10. Measure the small org's delivery latency before and after fairness.

---

### v1.3 — The warehouse
**Videos:** 23–24. Object storage

**Story:** Some parcels are huge, and the archive fills with years of old receipts. Big parcels and old records move to a cheap warehouse; the archive keeps a claim ticket pointing to each.

**The pain:** Large payloads and millions of attempt records bloat Postgres.

**Principles to grasp**
- Object storage vs file systems vs databases. Keys, buckets, immutability.
- Presigned URLs, multipart upload, lifecycle and retention policies.

**Build**
- MinIO (S3-compatible) in docker-compose.
- Payloads over a size threshold stored in object storage; Postgres keeps the key.
- Nightly archival of old attempts to object storage; download via presigned URL.

---

### v1.4 — Live wire
**Videos:** 25. Real-time backends

**Story:** Customers stop calling to ask "has my letter arrived yet?" The post office puts up a live departures board instead.

**The pain:** The dashboard has to poll to see new deliveries.

**Principles to grasp**
- Polling, long polling, Server-Sent Events, WebSockets: trade-offs.
- Fan-out across many API instances (Postgres `LISTEN/NOTIFY` or Redis pub/sub).

**Build**
- `GET /v1/apps/{app_id}/stream` using SSE that pushes attempt results live.
- Pub/sub so a worker on any machine can notify viewers connected to any API instance.

**Done when**
- [ ] **Gate E review.**

---

## Phase 5: Professional polish (Gate F → v2.0)

### v1.5 — Prove it
**Videos:** 26. Testing for backend engineers

**Story:** Mystery shoppers. Before the post office opens each day, someone tries every service, including the weird ones, and reports what broke.

**The pain:** Tests were added version by version. Now they need a strategy: what to test at which level, and how to trust the delivery guarantees.

**Build**
- Write `docs/testing.md`: unit, integration, end-to-end, contract, load, chaos. What lives where and why.
- A chaos test that randomly kills workers and the DB connection during a delivery run, then asserts at-least-once delivery for every message.
- CI runs unit + integration tests on every PR; coverage reported.

---

### v1.6 — The franchise manual
**Videos:** 27. The twelve-factor app

**Story:** To open a branch in any city, the post office writes a manual so any branch runs the same way.

**Build**
- Audit Hookline against all 12 factors in `docs/twelve-factor.md`. Fix every gap.
- Multi-stage Dockerfile (distroless), GitHub Actions build and push, deploy to one real host (Fly.io, Railway, or a small VPS).

---

### v1.7 — Contract as code
**Videos:** 28. OpenAPI

**Story:** The printed catalogue from v0.7 becomes a machine-readable one, so customers' software can read it directly.

**Build**
- Convert `docs/api.md` into an OpenAPI 3.1 spec. Validate requests against it in tests.
- Generated reference docs (Redoc or Scalar) and a generated Go client.
- CI check that fails on breaking API changes.

---

### v1.8 — Sealed, signed letters
**Videos:** 29. Webhooks

**Story:** Forged letters are showing up at customers' doors claiming to be from the post office. So every letter gets a wax seal only the real post office can make, with the time it was sealed so old letters can't be re-sent.

**The pain:** Since v0.10, Hookline has been sending plain POSTs. Receivers can't tell a real Hookline delivery from a forgery or a replay.

**Build**
- Follow the Standard Webhooks spec: `webhook-id`, `webhook-timestamp`, `webhook-signature` (HMAC-SHA256).
- Secret rotation with an overlap window.
- A small verification library for receivers (Go, and optionally Node/Python) and a consumer guide: verify, reject old timestamps, dedupe by `webhook-id`, respond fast and process async.
- Hookline as a receiver too: ingest webhooks from a provider like GitHub and fan them out to endpoints.

---

### v2.0 — The direct line
**Videos:** 30. WebSockets

**Story:** Some customers live behind a gate the couriers can't get through (a laptop, a private network). So the post office gives them a direct phone line: they call in, keep the line open, and letters are read to them over it.

**The pain:** Developers can't receive webhooks on `localhost`, and some customers can't expose a public URL.

**Build**
- A WebSocket delivery transport with auth, heartbeats, acknowledgements and resume from last acknowledged message.
- A `hookline listen --forward-to localhost:3000` CLI (like `stripe listen`) that opens the socket and forwards deliveries to a local server.
- Load test: how many concurrent sockets can one instance hold?

**Done when**
- [ ] **Gate F portfolio review.** Tag `v2.0.0`.
- [ ] README shows architecture diagram, demo GIF, perf numbers, SLOs, and links to every ADR.

---

## Beyond v2.0: expansion tracks

The roadmap ends at v2.0, but Hookline keeps growing. The v2.0 code already has separate binaries (API, worker, indexer, WebSocket gateway) and a queue between them, which is what makes these tracks possible. Each track uses the same method: new source material, a brief from the supervisor, versions with a pain, a build and a break-it experiment.

| Track | Versions | What it adds |
|---|---|---|
| Microservices | v3.x | Split API, delivery, search and gateway into separately deployed services. gRPC between them, a real broker (NATS JetStream or Kafka), the saga and outbox patterns across services, service discovery, Kubernetes. Includes an ADR on whether splitting was worth it. |
| AI integration | v4.x | An LLM that explains why deliveries fail and suggests fixes, payload schema inference, semantic search over messages, and an MCP server so AI agents can subscribe to and act on events. Covers prompt/response logging, cost and latency budgets, and evaluation. |
| Platform | v5.x | Multi-region delivery, usage metering and billing, a web dashboard, a public status page. |

Add the source material for a track to `docs/curriculum-map.md` before starting it.
