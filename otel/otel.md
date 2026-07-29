# OpenTelemetry — Cheat Sheet 

## Day 1 — What OTel is & the signals
- **OTel** = vendor-neutral standard to **generate + collect** telemetry; solves fragmentation / vendor lock-in.
- Two things: **spec/data model (OTLP)** + **SDKs & Collector**.
- **Three signals:** traces (one request's journey), metrics (aggregated numbers over time), logs (timestamped event records).
- **Trace stitching:** shared `trace_id` + per-hop `span_id` + `parent_span_id` (says *where in the tree*).
- **Context propagation:** inject on send / extract on receive — HTTP `traceparent` header, Kafka record headers.
- **Metrics do NOT propagate** (unlike traces); each service emits its own, correlated by shared attributes.
- **Exemplar** = sample `trace_id` on a metric point → jump from a metric spike to one real trace.
- **IoT reality:** device can't run an SDK → trace starts at the **ingest gateway**; gateway translates payload → OTel signals and stamps `device.id` on all three.

## Day 2 — Data model
- **Span** = one unit of work: `span_id`, `parent_span_id`, `trace_id`, name, start/end (→ duration = latency), **status (OK/ERROR → which hop failed)**, attributes, events.
- **Span context** = the traveling slice: `trace_id` + `span_id` + **trace flags** (sampling decision). Serialized into `traceparent`.
- **Attribute** = key-value tag on a span; *varies per request* (`tenant_id`, `device.id`, `db.system`).
- **Resource** = who/where is emitting; *constant per process* (`service.name`, `k8s.pod.name`, `cloud.region`).
- **Anchor:** resource = *which pod*, attribute = *which request*.
- **Cardinality:** metrics = one **time series per label combo** → high-cardinality label (`device.id`, millions) = explosion / OOM ("cardinality bomb"). Traces = discrete records → `device.id` is fine & useful.
- **OTLP** = standard *envelope* (wire format). **Semantic conventions** = standard *attribute names* (`db.system`, `http.request.method`) → why backends are swappable with no code change.

## Day 3 — Architecture
- Flow: **instrument → OTLP → Collector → backend.**
- **Auto-instrumentation** = framework boundaries (HTTP, Kafka, Cassandra) with no tracing code. **Manual** = your business logic. Use **both**.
- **Go nuance:** compiled, no runtime patching → uses wrapper libs (`otelhttp`, `otelsarama`); more explicit than Java's agent.
- **OTLP transports:** gRPC (binary, default) + HTTP/protobuf.
- **Why a Collector (vs app→backend direct):** decoupling/vendor-switch, offloading (batch/retry/compress), buffering when backend down, central processing (PII/filter/sample), normalization.
- **Deployment:** **agent** (sidecar/daemonset, near service) + **gateway** (central pool, heavy processing). Real setups use **both**.

## Day 4 — Inside the Collector
- **Receivers → Processors → Exporters**, wired into **pipelines** (one per signal).
- **Receivers:** OTLP receiver; **Prometheus receiver** = migration bridge (scrapes old metrics → OTLP).
- **Processors:** batch, memory limiter, attribute/resource (add/strip PII), filter (drop noise), **tail sampling**.
- **Head sampling** = random % up front (blind, can miss errors). **Tail sampling** = buffer whole trace, then keep errored/slow, drop boring. Better at IoT scale.
- **Tail sampling → gateway only:** needs the *complete* trace; a per-pod agent sees only its own spans. Also stateful/memory-heavy.
- **Exporters:** OTLP exporter → backend. Swap backend = change exporter config (no code). Multiple exporters = dual-export during migration.

## Day 5 — Advantages / Disadvantages / tc1 story
- **Advantages:** vendor neutrality; one standard for all 3 signals; cross-signal correlation; broad ecosystem/auto-instrumentation; Collector as control point.
- **Disadvantages:** Collector = infra you run/scale/monitor; **logging least mature signal**; cardinality footguns; instrumentation/sampling overhead; complexity + spec churn.
- **tc1 → OTel story:** *before* = fragmented Prometheus + Jaeger + Alertmanager; *migration* = Collector + Prometheus receiver, migrate signal-by-signal, dual-export; *scale judgment* = tail sampling in gateway, `device.id` never a metric label (use exemplars).
- **tc1 mapping:** Prometheus→OTLP metrics · Jaeger→OTLP traces · device.id→attribute join key · Kafka context→record headers · edge collector→agent · multi-DC→resource attributes · rate-limit/batch→processors.

---
**Killer one-liners:** trace_id says *which trace*, parent_span_id says *where in the tree* · exemplar = metric→trace bridge · resource = which pod / attribute = which request · OTLP = envelope / semantic conventions = language · tail sampling needs the whole trace → gateway only.
