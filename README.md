# solid_go

A hands-on Go learning repo covering SOLID design principles, concurrency patterns, networking, and database integrations — all demonstrated through a running invoice domain example.

---

## Concepts at a Glance

| # | Concept | Folder | Description |
|---|---------|--------|-------------|
| 1 | **Single Responsibility Principle** | [s_single_responsibility](solid_principles_go/s_single_responsibility) | Each function does exactly one job — print, save, or email an invoice |
| 2 | **Open / Closed Principle** | [o_open_closed](solid_principles_go/o_open_closed) | Add new discount strategies without touching existing code |
| 3 | **Liskov Substitution Principle** | [l_liskov](solid_principles_go/l_liskov) | Any `InvoiceExporter` implementation can replace another without breaking callers |
| 4 | **Interface Segregation Principle** | [i_interface_segregation](solid_principles_go/i_interface_segregation) | Split fat interfaces into `InvoicePrinter` and `InvoiceExporter` so types only implement what they need |
| 5 | **Dependency Inversion Principle** | [d_dependency_inversion](solid_principles_go/d_dependency_inversion) | `InvoiceService` depends on an `InvoiceRepository` interface, not a concrete MySQL or MongoDB struct |
| 6 | **Goroutine Patterns** | [goroutine_patterns](./goroutine_patterns) | Fan-in (merge multiple channels), fan-out worker pool, and multi-stage pipeline |
| 7 | **RWMutex + Singleflight** | [RWMutexandsingleflight](./RWMutexandsingleflight) | Concurrent-safe cache with read/write locks; singleflight collapses duplicate in-flight DB calls into one |
| 8 | **Worker Pool (HTTP)** | [worker_pool_api_processer](./worker_pool_api_processer) | Fixed-size goroutine pool processes HTTP-submitted jobs; request blocks until its result is ready |
| 9 | **Rate Limiting (Token Bucket)** | [ratelimiting](./ratelimiting) | Token-bucket algorithm with mutex-protected state controls request throughput |
| 10 | **gRPC** | [grpc](./grpc) | Unary, server-streaming, client-streaming, and bidirectional-streaming RPCs over Protocol Buffers |
| 11 | **REST API** | [rest](./rest) | Layered HTTP server: handler → use-case → repository, with graceful shutdown |
| 12 | **GraphQL** | [graphql](db/graphql) | Schema-first GraphQL API using gqlgen with query and mutation resolvers |
| 13 | **PostgreSQL** | [postgres](db/postgres) | `database/sql` with JSONB columns, parameterised queries, and nested JSON marshalling |
| 14 | **MongoDB** | [mongo](db/mongo) | CRUD + compound index creation + `explain` plan analysis via the official Go driver |
| 15 | **Batch Flusher** | [flusher](./flusher) | HTTP handlers buffer events in-memory; a ticker-driven goroutine flushes them to the DB in batches, with a final flush on graceful shutdown |
| 16 | **Rate Limit Middleware** | [ratelimit_middleware](./ratelimit_middleware) | Per-user token-bucket rate limiter applied as an httprouter middleware; each user gets their own bucket keyed by `X-User-ID` |
| 17 | **URL Shortener** | [url_shortner](./url_shortner) | Short-code → full URL redirect service backed by an RWMutex cache and singleflight to prevent cache-stampede DB calls |
| 18 | **Fan-in/Fan-out (Prime Finder)** | [fanin_fanout_prime_number](./fanin_fanout_prime_number) | CPU-bound prime search fanned out across `NumCPU` workers, then fanned back into a single result stream |
| 19 | **Observability & Metrics** | [observability](observability/observability.md) | Cheat-sheet notes: percentiles/histograms, RED & USE methods, SLI/SLO/SLA, MTTA/MTTR/MTBF, cardinality, and alerting philosophy |
| 20 | **OpenTelemetry (OTel)** | [otel](otel/otel.md) | Cheat-sheet notes: the three signals, data model, Collector architecture, sampling strategies, and migration story |
| 21 | **System Design — Data Ingestion** | [data_ingestion](system_design/data_ingestion.md) | Cheat-sheet notes: DMSP multi-tenant IoT metering platform — requirements, Kafka/KEDA ingest burst, Cassandra modelling, durability, retention/TWCS, and interview deep dives |

---

## SOLID Principles

SOLID is an acronym for five object-oriented design principles that make software easier to maintain and extend.

### S — Single Responsibility
> A type or function should have one reason to change.

`Invoice` is a plain data struct. Printing, saving, and emailing are separate functions. If the email format changes, only `SendEmail` needs updating.

### O — Open / Closed
> Open for extension, closed for modification.

`PrintFinalAmount` accepts any `DiscountStrategy`. Adding a `SeasonalDiscount` requires zero changes to existing code — just a new type that satisfies the interface.

### L — Liskov Substitution
> Subtypes must be usable wherever their base type is expected.

`PDFExporter` and `CSVExporter` both satisfy `InvoiceExporter`. `ProcessExport` doesn't know or care which one it receives.

### I — Interface Segregation
> No type should be forced to implement methods it doesn't use.

`BasicReporter` only prints — it implements `InvoicePrinter`. `FileExporter` only exports — it implements `InvoiceExporter`. `FullExporter` does both. No type carries dead methods.

### D — Dependency Inversion
> High-level modules should depend on abstractions, not concrete types.

`InvoiceService` holds an `InvoiceRepository` interface. Swapping the backing store from MySQL to MongoDB requires no change to the service.

---

## Concurrency Patterns

### Fan-in
Multiple goroutines produce values on separate channels. A `merge` function reads all of them and forwards to a single output channel, using a `WaitGroup` to close the output only after every producer finishes.

### Fan-out / Worker Pool
Work is distributed across N goroutines reading from a shared `jobCh`. Results go to a `resCh`. A separate goroutine closes `resCh` once all workers finish (via `WaitGroup`). Context cancellation propagates via `select`.

### Pipeline
Goroutines are chained: each stage reads from the previous stage's channel, transforms the value, and writes to the next stage's channel. Closing a stage's output channel signals the downstream stage to stop.

### RWMutex
`sync.RWMutex` allows many concurrent readers (`RLock`) but only one writer at a time (`Lock`). Used here to protect a shared in-memory cache under high read load.

### Singleflight
`golang.org/x/sync/singleflight` ensures that when many goroutines request the same key simultaneously, only **one** real DB call is made. All waiters share the result. `atomic.Int64` tracks how many actual DB calls happened.

### Token Bucket Rate Limiter
A bucket holds up to `capacity` tokens. Each request consumes one token (`Allow()`). Tokens are refilled at `refillRate` per `refillDuration`. A `sync.Mutex` serialises concurrent access to the mutable token count.

### Fan-in / Fan-out (Prime Finder)
A single `repeatFunc` goroutine produces an unbounded stream of random integers. Since primality testing (naive trial division) is CPU-bound, the stream is fanned out to one `primeFinder` goroutine per `runtime.NumCPU()`, each reading from the same shared channel. Their outputs are fanned back into one stream via `fanIn`, which uses a `WaitGroup` to close the merged channel only after every worker's channel has closed. A `take` helper caps consumption to the first N primes found.

---

## Networking

### gRPC
- **Unary**: one request → one response (`CreateInvoice`)
- **Server streaming**: one request → stream of responses (`ListInvoices`)
- **Client streaming**: stream of requests → one response (`UploadInvoices`)
- **Bidirectional streaming**: concurrent send/receive on both sides (`SyncInvoices`)

Service definitions live in `.proto` files; generated Go code is in `grpc/proto/`.

### REST API
Three-layer architecture:
- **handler** — HTTP parsing and response writing
- **use-case** — business logic (thin closures here)
- **repo** — data access

Uses `httprouter` for routing, `envconfig` for config, and a graceful shutdown pattern (`signal.Notify` + `http.Server.Shutdown`).

### GraphQL
Schema-first with `gqlgen`. The schema declares types, queries, and mutations; `schema.resolvers.go` contains the implementations. An LRU query cache and Automatic Persisted Queries (APQ) are enabled for performance.

---

## Databases

### PostgreSQL
Uses the standard `database/sql` package with the `lib/pq` driver. Demonstrates:
- DDL execution (`CREATE TABLE IF NOT EXISTS`)
- JSONB column storage with `json.Marshal` / `json.Unmarshal`
- JSONB path operators (`->>`, `#>>`) in parameterised queries

### MongoDB
Uses the official `go.mongodb.org/mongo-driver`. Demonstrates:
- CRUD operations (`InsertOne`, `FindOne`, `UpdateOne`, `DeleteOne`, `InsertMany`)
- Compound index creation
- `explain` command to compare index vs. collection scan execution plans

---

## Patterns

### Batch Flusher
High-throughput writes can overwhelm a database if every HTTP request issues its own `INSERT`. The flusher pattern solves this by buffering events in memory and writing them in bulk on a timer:
- HTTP handler acquires a mutex, appends the event to a slice, and returns immediately (no DB call per request).
- A background goroutine ticks every N seconds, locks the buffer, drains it to the DB, then clears the slice.
- On graceful shutdown, the `quit` channel is closed which triggers one final flush so no buffered events are lost.

### Rate Limit Middleware
Rather than embedding rate-limiting logic inside each handler, a middleware function wraps any `httprouter.Handle` and intercepts requests before they reach business logic:
- `RateLimiterStore` holds a per-user `rateLimiter` keyed by `X-User-ID`.
- Each user gets a token bucket with capacity 10, refilled fully every second.
- The middleware returns `429 Too Many Requests` when the bucket is empty, leaving the inner handler untouched.
- A single store-level mutex guards the map; fine for moderate concurrency — shard if needed.

### URL Shortener
A redirect service that resolves short codes to full URLs with two layers of protection against thundering-herd DB load:
- **RWMutex cache**: multiple goroutines can read the in-memory map simultaneously (`RLock`); writes are exclusive (`Lock`). This keeps redirects fast under high concurrency.
- **Singleflight**: on a cache miss, `group.Do` ensures only one `GetFromDB` call runs per short code at a time. All other goroutines racing for the same key block and share the single result, preventing a storm of identical DB queries.

---

## Observability

Two study cheat-sheets that pair with the tracing and metrics used throughout this repo. These are reference notes rather than runnable code.

### Observability & Metrics
[observability/observability.md](observability/observability.md) — the fundamentals of measuring a system:
- **Percentiles & histograms** — why averages hide tail pain, and how cumulative buckets let you query any percentile later.
- **RED** (Rate, Errors, Duration) for request-driven services and **USE** (Utilization, Saturation, Errors) for resources.
- **SLI / SLO / SLA** and error budgets; **MTTA / MTTR / MTBF** incident metrics.
- **Histogram vs Summary**, **cardinality** footguns, **metrics vs logs vs traces**, distributed tracing, and alerting philosophy (symptom vs cause, burn-rate alerts).

### OpenTelemetry (OTel)
[otel/otel.md](otel/otel.md) — the vendor-neutral standard for generating and collecting telemetry:
- **Three signals** (traces, metrics, logs) and how a `trace_id` + `span_id` + `parent_span_id` stitch a request's journey together.
- **Data model** — spans, span context, attributes vs resources, and cardinality trade-offs; **OTLP** wire format and semantic conventions.
- **Collector architecture** — receivers → processors → exporters, agent vs gateway deployment, and head vs tail sampling.
- **Migration story** — moving a fragmented Prometheus + Jaeger stack to a Collector-centric OTel pipeline.
