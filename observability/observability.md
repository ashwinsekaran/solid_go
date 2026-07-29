# Observability & Metrics — Cheat Sheet

## Day 1: Percentiles & Histograms
- **Average** blends all values — hides outlier pain (one slow request can be masked by many fast ones)
- **Percentile** = value at a specific *sorted position*. p95 = 95% of requests were at or below this value; 5% were slower
- p99 ≠ worst case — the single slowest request is p100/max
- Can't sort millions of live values → use **histograms**: pre-defined **cumulative buckets** (`≤10ms, ≤50ms...`), each a simple counter
- A value increments *every* bucket ≥ its own value (not just one range)
- Percentile computed later via interpolation across bucket boundaries (`histogram_quantile()`)
- **Exact per-record data** (e.g., one order's duration) lives in **logs/traces**, not metrics — histograms lose individual identity

## Day 2: RED Method (request-driven services)
- **R**ate — requests/sec — "is anyone using this?"
- **E**rrors — % failing — "is it broken?"
- **D**uration — latency (p95/p99) — "is it slow?"
- Fits APIs/RPC services; doesn't fit batch jobs, queues, cron

## Day 3: USE Method (resources)
- **U**tilization — % of capacity busy ("how full")
- **S**aturation — is work queuing/backing up ("is it keeping up") — e.g., Kafka lag, pending compactions
- **E**rrors — failed operations at resource layer (disk/network/memory)
- High utilization ≠ high saturation, and vice versa — they're independent signals
- RED = user-facing symptom view; USE = internal resource view — run both together

## Day 4: SLI / SLO / SLA
- **SLI** = actual measured value (e.g., 99.95% success)
- **SLO** = internal target (e.g., 99.9%) — tighter than SLA, gives early-warning buffer
- **SLA** = external/contractual promise (e.g., 99.5%) — breach = refunds/credits/legal
- **Error budget** = allowed failure margin (100% − SLO%); burn it → slow down releases, focus on stability

## Day 5: MTTA / MTTR / MTBF
- **MTTA** = Mean Time To Acknowledge — alert fired → human responds (paging/on-call health)
- **MTTR** = Mean Time To Resolve/Recover/Repair — problem start → fixed (state which definition you mean)
- **MTBF** = Mean Time Between Failures — reliability/frequency, not response speed
- Good MTTR + bad MTBF = fragile system, well-drilled response. Compare **total downtime**, not just one metric

## Day 6: Histogram vs Summary
- **Histogram**: bucket counts, computed server-side at query time. Additive across pods (counts sum, then compute once) — works for fleet-wide percentiles
- **Summary**: percentile computed client-side at observation time. NOT additive — can't average percentiles across pods/instances
- Summary locks in percentiles at instrumentation time (no p99.9 later without redeploy); histogram lets you query any percentile later from stored bucket counts
- Distributed systems → prefer **Histograms**

## Day 7: Cardinality
- Cardinality = # of unique label/tag value combinations for a metric → each combo = a separate stored time series
- Multiplies across labels (e.g., 10 status codes × 5 methods × 4 regions = 200 series)
- **Never label with unbounded values**: user_id, email, order_id, pod_name (ephemeral), IPs, timestamps
- Cause of "cardinality explosion" → memory/storage/query cost blows up, system slows/crashes
- Need per-entity detail? → use **logs**, not metric labels

## Day 8: Metrics vs Logs vs Traces
- **Metrics** — aggregated numbers, cheap, real-time — "is something wrong, roughly how bad?"
- **Logs** — per-event records — "what exactly happened?"
- **Traces** — full request journey across services (spans) — "where did time go, across the system?"
- General investigation: Metrics (narrow) → Logs (detail) → Traces (pinpoint cross-service latency)
- One specific known request: Trace first (if locatable) → Logs → Metrics (to check if isolated or broad)
- **Trace ID** ties all three together (OpenTelemetry = unified standard for this)

## Day 9: Distributed Tracing
- **Span** = one unit of work: name, start time, duration, parent span ID, attributes
- **Trace** = tree of spans, connected by a shared **trace_id**
- **Trace context** (trace_id + span_id) propagated via headers (`traceparent`) across service hops
- Dropped/missing propagation → new trace starts, original journey splits into two disconnected traces
- **Sampling** (can't trace 100% at scale):
  - **Head-based** — decide randomly upfront (cheap, may miss rare issues)
  - **Tail-based** — buffer everything, decide after seeing outcome (always catches errors/slow requests, but expensive/complex)
  - Common: blend both — small % baseline + always keep errors/high-latency

## Day 10: Alerting Philosophy
- **Symptom-based** — alert on user-facing pain (latency, error rate) → **page** immediately
- **Cause-based** — alert on internal conditions (CPU, disk, lag) → **ticket/warn**, not page (avoids alert fatigue)
- Best practice: combine both — e.g., page only if cause (CPU high) *and* symptom (latency/errors) both degrade
- **Most mature approach: alert on error budget burn rate**, not flat thresholds — ties alerts directly to SLO/business impact
  - Fast burn (budget gone in hours) → page
  - Slow burn (budget gone in weeks) → ticket
