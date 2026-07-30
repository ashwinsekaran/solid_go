# Metrics & Monitoring System Design — Observability Platform Cheat Sheet

> **Scenario:** designing a Dash0-like observability platform that collects **logs, metrics, and traces** from many servers.
> **Core idea:** one unified ingestion pipeline, but a *different storage system per telemetry type* — because logs, metrics, and traces have very different read/write patterns.
> Stack: **OpenTelemetry · Kafka · ClickHouse · Prometheus-compatible TSDB · Tempo · (optional) Flink**

## Contents

- [0 · Framing & high-level architecture](#0--framing--high-level-architecture)
- [1 · Requirements](#1--requirements)
  - [Functional requirements](#functional-requirements)
  - [Non-functional requirements](#non-functional-requirements)
- [2 · Scale estimate](#2--scale-estimate)
- [3 · Core entities](#3--core-entities)
- [4 · API surface](#4--api-surface)
  - [Ingestion APIs (OTLP)](#ingestion-apis-otlp)
  - [Query APIs](#query-apis)
  - [Correlation APIs](#correlation-apis)
- [5 · Architecture](#5--architecture)
  - [Edge — OpenTelemetry Collector agents](#edge--opentelemetry-collector-agents)
  - [Regional ingestion gateways](#regional-ingestion-gateways)
  - [Kafka — durability & buffering](#kafka--durability--buffering)
  - [Per-signal consumer pipelines](#per-signal-consumer-pipelines)
  - [Optional stream processing (Flink)](#optional-stream-processing-flink)
- [6 · Data flow — step by step](#6--data-flow--step-by-step)
- [7 · Cross-signal correlation](#7--cross-signal-correlation)
- [8 · Partitioning & tenant isolation](#8--partitioning--tenant-isolation)
- [9 · Retention policies](#9--retention-policies)
- [10 · Capacity estimates](#10--capacity-estimates)
- [11 · Bottlenecks & failure handling](#11--bottlenecks--failure-handling)
- [12 · MVP recommendation](#12--mvp-recommendation)
- [13 · Key tradeoff](#13--key-tradeoff)
- [14 · Interview script (5–7 min)](#14--interview-script-57-min)

---

## 0 · Framing & high-level architecture

**The one-sentence framing:** build a *unified ingestion pipeline*, but use *different storage systems per telemetry type* — because logs, metrics, and traces have very different read/write patterns.

**High-level flow:**

Apps / Servers → OTEL Collectors → Ingestion Gateway → Kafka → processing / routing → specialized stores → Query APIs / UI

**The system has four jobs:**

1. **Collect** telemetry from customer workloads.
2. **Ingest** it reliably at high throughput.
3. **Route** and optionally **enrich** it.
4. **Store** it in backends optimized for each signal, and expose a single UI for querying and correlation.

**Why split storage:** logs need high-throughput search, metrics need time-series aggregation, and traces need lookup by trace ID and span relationships. Forcing all three into one engine usually creates performance or cost problems later.

## 1 · Requirements

### Functional requirements

- Ingest **logs, metrics, and traces** from many servers and Kubernetes clusters.
- **Logs:** near real-time search.
- **Metrics:** time-series dashboards and alert queries.
- **Traces:** distributed trace lookup by trace ID, service, duration, or error.
- **Correlation:** let users jump from a metric spike → related traces → related logs.

### Non-functional requirements

- Very high **write throughput** and **horizontal scalability**.
- **Tenant isolation** — a hot tenant must not crush everyone else.
- **Durability** — accepted telemetry is not lost.
- **Low-latency queries** on recent data.
- **Cost-efficient retention.**
- **Regional ingestion**, backpressure handling, and graceful degradation if one backend is slow.

## 2 · Scale estimate

Assume a rough scale and derive out loud:

| Assumption | Value |
| --- | --- |
| Customers | 1,000 |
| Hosts per customer (avg) | ~500 |
| Total hosts | ~500,000 |
| Per-host logs | 5 KB/s |
| Per-host metrics | 2 KB/s |
| Per-host traces | 1 KB/s |
| Per-host combined | ~8 KB/s |
| **Global ingest (before replication)** | **~4 GB/s** |

**Takeaway:** 4 GB/s is large enough that you definitely want a **buffered pipeline** and **partitioned storage** — not direct writes to a backend.

## 3 · Core entities

| Entity | Key fields | Notes |
| --- | --- | --- |
| **Log record** | timestamp, tenant_id, service, env, region, host, severity, trace_id / span_id (if present), message, structured attributes | Full-text / token search on `message` |
| **Metric sample** | timestamp, tenant_id, metric_name, value, labels, type (counter / gauge) | Label-based queries + rollups |
| **Trace (span)** | tenant_id, trace_id, span_id, parent_span_id, service, operation, duration, status, start_time, attributes | Tree of spans joined by `trace_id` |

**The correlation glue:** every signal carries shared metadata — `tenant_id`, `service`, `env`, `region`, and `trace_id` wherever possible. That shared metadata is what makes cross-signal navigation work.

## 4 · API surface

Split into **ingestion** and **query** APIs.

### Ingestion APIs (OTLP)

OTLP-based rather than custom JSON — these are hit by OTEL collectors, not browsers.

- `POST /otlp/v1/logs`
- `POST /otlp/v1/metrics`
- `POST /otlp/v1/traces`

### Query APIs

- `GET /logs/search?service=&env=&q=&start=&end=&cursor=`
- `GET /metrics/query?expr=&start=&end=&step=`
- `GET /traces/search?service=&operation=&minDuration=&status=&start=&end=`
- `GET /traces/{traceId}`

### Correlation APIs

- `GET /logs/by-trace/{traceId}`
- `GET /services/{service}/overview`

## 5 · Architecture

### Edge — OpenTelemetry Collector agents

On each customer host or cluster, run an **OpenTelemetry Collector** agent. It receives logs, metrics, and traces from apps or sidecars, then **batches, compresses, retries**, and **adds local metadata** (host, k8s pod, region, service name). Apps never talk to a central backend directly for every event.

### Regional ingestion gateways

Collectors send to a fleet of **regional ingestion gateways**. These are **stateless** and horizontally scalable. They handle **auth, tenant validation, rate limiting, basic schema checks**, then push accepted telemetry into Kafka.

### Kafka — durability & buffering

Kafka is the **durability and buffering** layer. It decouples ingestion from storage: if one downstream backend is slow, ingest doesn't fail — Kafka absorbs the burst and consumers lag temporarily. Multiple consumer groups read independently.

### Per-signal consumer pipelines

From Kafka, split into three consumer pipelines:

| Signal | Store (default choice) | Indexed / queried by |
| --- | --- | --- |
| **Logs** | ClickHouse (or OpenSearch) | Partitioned by time + tenant; filtered by service, severity, env. ClickHouse preferred for scale & compression. |
| **Metrics** | Prometheus-compatible TSDB — Mimir / Cortex / VictoriaMetrics | Label-based queries, rollups, downsampling. |
| **Traces** | Tempo (or a ClickHouse-based trace store) | trace_id, service, operation, status, duration. |

### Optional stream processing (Flink)

**Optional for MVP.** Between Kafka and the sinks, a stream processor (e.g. Flink) can normalize schemas, enrich records with tenant metadata, derive service graphs, correlate logs with trace IDs, compute error-rate streams, or apply central sampling. Add it **later**, only if you need centralized enrichment or derived telemetry — not on day one.

## 6 · Data flow — step by step

1. A service emits logs, metrics, and traces.
2. The local **OTEL collector** receives them, batches, compresses, attaches metadata, and forwards upstream.
3. The **ingestion gateway** authenticates the tenant, enforces quotas, and publishes records to Kafka topics — `logs`, `metrics`, `traces`.
4. **Kafka** stores the stream durably; multiple consumer groups read independently.
5. The **logs consumer** reads the logs topic → writes to the log store.
6. The **metrics consumer** reads the metrics topic → writes to the TSDB.
7. The **traces consumer** reads the traces topic → writes to the trace store.
8. The **UI never queries Kafka directly** — it calls query services:
   - logs query service → ClickHouse / OpenSearch
   - metrics query service → TSDB
   - traces query service → Tempo / trace DB

## 7 · Cross-signal correlation

The most important enabler is **consistent metadata** on every record: `tenant_id`, `service`, `env`, `region`, and `trace_id` where possible.

That shared metadata powers workflows like:

- "Show me logs related to **this trace**."
- "Show me traces during **this CPU spike**."
- Alert fires on a metric → open related traces → jump from a span to matching logs.

## 8 · Partitioning & tenant isolation

- Partition Kafka by **`tenant_id` + time bucket** or **`tenant_id` + service**, depending on skew.
- **Tenant isolation matters a lot** — a hot tenant should not create hotspots or starve everyone else.
- Rate limit an over-quota tenant rather than destabilizing the whole platform.

## 9 · Retention policies

Each signal gets a different policy:

| Signal | Hot retention | Long-term |
| --- | --- | --- |
| **Logs** | 7–30 days hot (expensive) | Archive to object storage |
| **Metrics** | Raw samples for a shorter period | Downsampled rollups kept longer |
| **Traces** | Sampled | Keep **all error traces**, only a percentage of healthy ones |

## 10 · Capacity estimates

- Global ingest ~**4 GB/s**; Kafka must handle that **plus replication** (RF=3 multiplies broker network + disk load).
- If one partition safely handles ~**50 MB/s**, then 4 GB/s → **~80 partitions minimum**. In practice provision **several hundred** for headroom and to distribute hot tenants.
- **Logs ≈ 60% of ingest** → ~**2.4 GB/s** → roughly **~207 TB/day raw**. Compression matters a lot; ClickHouse reduces this significantly depending on data shape.

## 11 · Bottlenecks & failure handling

- **Backend slows down** → Kafka absorbs the burst; consumer lag grows temporarily instead of dropping data.
- **Collector can't reach ingestion** → it buffers locally and retries.
- **Tenant exceeds quota** → rate limit that tenant, don't destabilize the platform.
- **Trace volume explodes** → apply dynamic sampling upstream.

## 12 · MVP recommendation

> **MVP:** OTEL collectors for edge batching · Kafka for durable fanout · **ClickHouse** for logs · a **Prometheus-compatible TSDB** for metrics · **Tempo** for traces. Add stream processing (Flink) **later**, only if centralized enrichment or derived telemetry is needed.

This lands well because it's practical and avoids over-engineering.

## 13 · Key tradeoff

**Operational complexity vs. flexibility.** Separate storage backends increase system complexity, but they're the right call because each telemetry type has very different workload characteristics. Forcing logs, metrics, and traces into one storage engine usually creates performance or cost problems later.

## 14 · Interview script (5–7 min)

> I'd design this as a **unified telemetry ingestion platform** for logs, metrics, and traces — but I would **not** store all three the same way, because their access patterns are very different. Logs need high-throughput search, metrics need time-series aggregation, and traces need lookup by trace ID and span relationships.
>
> At a high level: **Applications/Servers → OTEL Collectors → Ingestion Gateway → Kafka → separate processing pipelines → specialized storage backends → query APIs/UI.**
>
> **Functional requirements:** ingest logs, metrics, and traces from many customer workloads; near real-time log search; metrics dashboards and alert queries; trace search and trace detail views; and correlation across all three signals.
>
> **Non-functional requirements:** high write throughput, durability, horizontal scalability, tenant isolation, and low-latency queries for recent data.
>
> I'd **start at the edge.** On each VM, container host, or Kubernetes cluster I'd run an **OpenTelemetry Collector** — local buffering, batching, compression, retries, and enrichment with metadata like host, service, region, environment, pod, and tenant. Applications shouldn't talk directly to a central backend for every event.
>
> Those collectors send to **regional ingestion gateways** — stateless and horizontally scalable. Their job is authentication, tenant validation, quota enforcement, basic schema validation, and writing accepted data into **Kafka**. Kafka is the key buffer: it decouples ingestion from storage, so if one backend becomes slow we don't immediately lose data or reject all writes.
>
> From Kafka I'd **split by telemetry type.** The **logs** topic → a log search store, defaulting to **ClickHouse** for efficient ingest, compression, and analytical filtering. The **metrics** topic → a **time-series backend** (Mimir, Cortex, or VictoriaMetrics) for label-based queries, rollups, and downsampling. The **traces** topic → a **tracing backend** (Tempo or a trace-oriented store), indexed by trace ID, service, operation, duration, and error status.
>
> **Data flow:** an app emits telemetry; the OTEL collector receives, batches, enriches, and forwards it; the gateway authenticates the tenant and writes to Kafka; independent consumer groups read from Kafka and write into their own backends. The **UI never queries Kafka directly** — query services read from the appropriate backend depending on whether the user is searching logs, running a metrics query, or opening a trace.
>
> The most important thing for **cross-signal correlation** is consistent metadata — every log, metric, and trace carries `tenant_id`, `service`, `env`, `region`, and ideally `trace_id`. That lets the UI support: alert fires on a metric spike → user opens related traces → jumps from a span to matching logs.
>
> **APIs:** ingestion uses OTLP endpoints — `/otlp/v1/logs`, `/otlp/v1/metrics`, `/otlp/v1/traces`. Reads expose `/logs/search`, `/metrics/query`, `/traces/search`, and `/traces/{traceId}`.
>
> **Scale:** ~500,000 hosts, each ~8 KB/s combined ≈ **4 GB/s globally before replication.** At that scale Kafka needs many partitions and strong tenant isolation — I'd partition primarily by tenant, then by service or time bucket, so one large tenant doesn't create hotspots.
>
> **Retention** per signal: logs 7–30 days hot then archive to object storage; metrics raw for a shorter time plus downsampled rollups longer; traces sampled — keep all error traces and a percentage of healthy ones.
>
> **Failure handling:** if a backend slows down, Kafka absorbs the pressure and consumers lag temporarily; if the collector loses connectivity, it buffers locally and retries; if a tenant exceeds quota, we rate limit that tenant instead of impacting the whole system.
>
> For **MVP** I'd keep it simple: OTEL Collectors, ingestion gateways, Kafka, ClickHouse for logs, a Prometheus-compatible TSDB for metrics, and Tempo for traces. I'd only add **Flink** later for centralized enrichment, correlation, anomaly detection, or derived telemetry — optional, not required on day one.
>
> The main **tradeoff** is operational complexity versus flexibility. Separate backends add complexity, but each telemetry type has very different workload characteristics — forcing them into one engine usually causes performance or cost problems later.
