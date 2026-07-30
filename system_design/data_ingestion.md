# DMSP System Design — Data Ingestion Cheat Sheet

> **Source:** imported verbatim from `DMSP_System_Design_CheatSheet.docx`.
> Multi-tenant IoT metering platform — a 5-step interview flow plus 4 timeboxed deep dives.
> Stack: **Go · Kafka · Cassandra · KEDA · Azure/K8s**

## Contents

- [0 · Interview flow & time budget](#0--interview-flow--time-budget)
- [1 · Requirements](#1--requirements)
  - [Functional — five capabilities](#functional--five-capabilities)
  - [Non-functional — five constraints](#non-functional--five-constraints)
  - [Scale — ASK, then derive out loud](#scale--ask-then-derive-out-loud)
- [2 · Core entities](#2--core-entities)
- [3 · API + dataflow](#3--api--dataflow)
  - [Upstream — one sentence, then cross the boundary](#upstream--one-sentence-then-cross-the-boundary)
  - [Write dataflow](#write-dataflow)
  - [Read dataflow](#read-dataflow)
  - [APIs](#apis)
  - [Two load balancers — ingest and read, kept separate](#two-load-balancers--ingest-and-read-kept-separate)
- [4 · High-level design (satisfy the FRs only)](#4--high-level-design-satisfy-the-frs-only)
- [5 · Deep dives — 4 × ~5 min](#5--deep-dives--4--5-min)
  - [DD1 · Ingest burst — Kafka partitioning + KEDA (~5 min)](#dd1--ingest-burst--kafka-partitioning--keda-5-min)
  - [DD2 · Cassandra modelling + tenant isolation (~6 min)](#dd2--cassandra-modelling--tenant-isolation-6-min)
  - [DD3 · Durability + availability (~4 min)](#dd3--durability--availability-4-min)
  - [DD4 · Retention, TTL, TWCS + archival (~5 min)](#dd4--retention-ttl-twcs--archival-5-min)
- [6 · Standby modules — do NOT volunteer](#6--standby-modules--do-not-volunteer)
  - [Why Cassandra — and the honest concession ★ most likely unprompted question](#why-cassandra--and-the-honest-concession--most-likely-unprompted-question)
  - [Aggregation job](#aggregation-job)
  - [Alerting — Flink SLA monitoring ← the section you flagged](#alerting--flink-sla-monitoring--the-section-you-flagged)
  - [UI dashboard query ranges](#ui-dashboard-query-ranges)
  - [Reports](#reports)
  - [Rate limiting — token bucket](#rate-limiting--token-bucket)
  - [Growth — onboarding 300 more customers (700 → 1,000)](#growth--onboarding-300-more-customers-700--1000)
- [7 · Key numbers — memorise these](#7--key-numbers--memorise-these)
- [8 · Delivery reminders](#8--delivery-reminders)

---

**DMSP — System Design Interview Cheat Sheet**

Multi-tenant IoT metering platform  ·  5-step flow + 4 timeboxed deep dives

*Go · Kafka · Cassandra · KEDA · Azure/K8s*

## 0 · Interview flow & time budget

| **#** | **Step** | **Time** | **One-job reminder** |
| --- | --- | --- | --- |
| 1 | Requirements — FR + NFR | 5 min | Say the count first. Derive the scale out loud. |
| 2 | Core entities | 2 min | Close with: the hierarchy IS the partition key. |
| 3 | API + dataflow | 4 min | Contract, THEN walk both journeys to Cassandra. |
| 4 | HLD (happy path only) | 8 min | Tag each move to an FR. Hold all NFR material. |
| 5 | Deep dives × 4 | 20–25 min | Thesis first. Announce the length. Hook at the end. |
| — | Close + offer standby | 2 min | Offer the doors you've prepared. |

#### Five delivery rules that fix the stumbling

- Say the COUNT before the list — "four capabilities: one… two…". This single habit kills most run-ons.

- Anchor every NFR to a number or a named mechanism. Fluency jumps when anchored; filler clusters in unanchored prose.

- Thesis first, detail second — if cut off at 90 seconds, the interviewer already has your answer.

- Announce the length: "let me take two minutes on the partition key." Commits you publicly, signals control.

- Close every dive with a hook offering the next door — they stay inside your prepared territory.

## 1 · Requirements

### Functional — five capabilities

- Ingest telemetry from customer smart meters.

- Normalize varied vendor payloads into a common schema.

- Store readings multi-tenant in Cassandra — queryable by customer, device and time, aggregatable over time.

- Let users log in and read their own data via dashboards, reports and monitoring.

- Monitor device SLA and alert customers immediately.

**Out of scope (say it explicitly): **the third-party DAAS layer, Full text search & Live data tailing

### Non-functional — five constraints

- **Throughput — high and bursty. **~3,500/s average, ~17,500–20,000/s peak (≈5×). Handled by KEDA scaling consumers on lag.

- **Tenant isolation. **One vendor must never read another's data.

- **High availability. **Billing data is revenue-critical → under CAP I favour availability over consistency → Cassandra multi-DC with LOCAL_QUORUM.

- **Durability. **No message lost end to end.

- **Alerting latency. **SLA breaches detected and alerted in under 1 minute.

### Scale — ASK, then derive out loud

**SAY IT  ***"Shall I assume a rough scale? I was running on the order of a million meters."*

- 700 vendors × ~1,500 end customers = ~1.05M meters

- 5-minute reporting interval → 1.05M ÷ 300s = ~3,500 msg/s average

- Meters fire on ALIGNED clock boundaries (:00, :05, :10) — they bunch, they don't smear

- Burst lands in ~60s window → 1.05M ÷ 60 = ~17,500 msg/s peak ≈ 5× average

- Daily volume ~300M messages/day → drives CASSANDRA sizing, not Kafka

**Sensitivity (know which knob moves what): **reporting interval drives the AVERAGE and the daily volume. Burst concentration drives the PEAK. Shortening the interval does NOT raise peak — it makes the same-size burst happen more often.

**WATCH OUT  **Don't oversell the scale. Own the framing: "this isn't an extreme-volume system — the interesting engineering is the burstiness and the multi-tenancy, not raw throughput." A sharp interviewer does this math themselves.

## 2 · Core entities

Three-level hierarchy — Vendor → Relationship → Meter. Per entity, three beats: name it, identify it, say what it belongs to.

| **Entity** | **Identified by** | **Belongs to** | **Note** |
| --- | --- | --- | --- |
| Customer / Vendor | customer_id | — | The tenant. A utility / meter vendor. |
| Relationship | rel_id | Customer | End customer / site where the meter is installed. |
| Device / Meter | device_id | Relationship | The smart meter itself. |
| Reading | customer + rel + device + timestamp | Device | Three flavours: raw meters, derived metrics, KPIs. |
| User | user_id | Customer | Read side. Permissions scope them to their tenant. |

**SAY IT  ***"A customer has many relationships, a relationship has many meters, a meter has many readings — and that customer → relationship → device → timestamp hierarchy is exactly my Cassandra partition key. My data model and my storage layout are the same shape."*

**WATCH OUT  **Do not drop this closer, and do not hedge on the device term. No "or whatever you call it" — say "a Device, a smart meter, identified by device_id."

## 3 · API + dataflow

### Upstream — one sentence, then cross the boundary

**SAY IT  ***"Meters publish over MQTT to a broker, and each vendor runs their own Device Gateway that subscribes to it. That gateway is customer-built, so my system starts at the API it calls."*

**Why that sentence earns its place: **"customer-built" is the justification for the ingest API being a PUBLIC CONTRACT — versioned, per-tenant rate limited, schema validated at MY edge. Context that explains a constraint you own.

### Write dataflow

Devices → MQTT (Mosquitto) → Device Gateway [customer-built: device auth, transpose to HTTP + bearer]

  → Ingress/LB → servsup [validate + normalize to common schema]

  → Kafka [topics: meters / metrics / kpi · key = customer_id + rel_id + device_id]

  → consumer groups [KEDA on lag] → Cassandra

### Read dataflow

User UI → LB → DMSP [Entra token → gRPC auth/check] → Cassandra

  scoped by customer_id FROM THE TOKEN, never the URL · dashboards read pre-aggregated tables

### APIs

**Ingest — POST, bearer token (service-to-service, not device auth):**

- /ingest/meters  ·  /ingest/metrics  ·  /ingest/kpi

**DMSP — GET, MS Entra:**

- /dmsp/devices — the tenant's meter list (bounded, no time param needed)

- /dmsp/meters?from=&to=

- /dmsp/kpi?type=&from=&to=

- /dmsp/reports?type=&from=&to=  → async, returns a report_id

- All: pagination; tenant scope derived from the token; raw queries capped at ~7 days

### Two load balancers — ingest and read, kept separate

**Answer first: two. **Ingest is machine traffic — 700 gateways bursting to 20k/s with service auth. Read is human traffic — a few hundred users on Entra tokens. Shared, an ingest burst exhausts connections and takes the dashboard down at exactly the moment operators need visibility. Same tooling, separate resource pools.

- Ingress also does: TLS termination, per-tenant rate limiting, request size caps, path routing to the three ingest deployments

- Uses /readyz to pull draining pods out of rotation — this is what makes KEDA scale-down safe

## 4 · High-level design (satisfy the FRs only)

**SAY IT  ***"I'll build the high-level design by walking the functional requirements in order, keeping scaling and failure handling for the deep dives."*

- **To satisfy ingestion: **the vendor's gateway posts to my ingest API — that's my trust boundary, so validation happens on my side.

- **To satisfy normalization: **servsup validates and normalizes varied payloads into a common schema, then publishes to one of three Kafka topics — meters, metrics, kpi, matching the three reading flavours. Kafka is there to decouple ingestion from persistence and absorb back-pressure, so a slow write path never pushes back into the customer's gateway.

- **To satisfy storage: **consumer groups subscribe to those topics and write into Cassandra, modelled per query pattern and scoped by customer — so tenant isolation falls out of the data model.

- **To satisfy the read path: **users hit the load balancer into DMSP, which authenticates via MS Entra and validates permissions over gRPC, then reads Cassandra scoped to that customer. Dashboards read pre-aggregated tables maintained by the agg job, so a dashboard query never scans raw readings.

**WATCH OUT  **HOLD BACK in HLD: partition key structure, KEDA config, RF/ISR, TWCS, rebalancing, agg job internals. A component a single request TRAVELS THROUGH belongs in HLD (Kafka does). A PROPERTY of how that path behaves under load belongs in deep dives (KEDA, partition keys do not).

**WATCH OUT  **Do not say "I know a lot of issues in this design." Say: "This satisfies the functional requirements; it doesn't yet address burst, availability or durability." Same content, no self-deprecation.

**SAY IT  ***"That satisfies the functional requirements. It doesn't yet address the burst, availability or durability — shall I start with ingestion scaling, or the Cassandra modelling?"*

## 5 · Deep dives — 4 × ~5 min

Order follows the data: ingest → storage → durability of that path → what happens to the data over time. Each dive picks up where the last ended, so no transition costs you context-setting. It also front-loads your strongest material.

### DD1 · Ingest burst — Kafka partitioning + KEDA (~5 min)

**THESIS  ***Peak is 5× average because meters fire on aligned boundaries, so I scale on Kafka consumer lag rather than CPU.*

#### Topics

- Three topics — meters / metrics / kpi. Different volumes, different consumer logic, independent scaling and independent failure. A metrics spike scales only the metrics group.

#### Partition key = customer_id + rel_id + device_id

- **Month is EXCLUDED** — it changes, and would remap a meter to a different partition on every rollover, breaking the one-meter-one-lane guarantee.

- Order preserved per meter — Kafka only guarantees order WITHIN a partition

- Same meter → same partition → same consumer → same Cassandra partition → efficient batched writes

- High cardinality spreads evenly. customer_id ALONE would hot-spot a large vendor onto one partition

- Composite doesn't reintroduce the hot-spot: Kafka hashes the whole key, so cardinality stays at device level

#### Sizing

- 12 partitions/topic — derived: peak ÷ per-consumer throughput (~5k/s) ≈ 4 consumers at peak, plus headroom

- 3 brokers (RF=3 needs ≥3), min.insync.replicas = 2, leaders spread so no broker carries all reads/writes

#### KEDA

- Scales on CONSUMER GROUP LAG, not CPU — consumers are I/O-bound on Cassandra, so CPU understates how far behind they are

- KEDA scales the DEPLOYMENT, not the group directly — but all pods share a GroupID, so they ARE the consumer group

- maxReplicaCount = partition count (12) — can't parallelize past partitions

- Producers scale on HPA instead — they ARE request-driven, so CPU is the right signal there

- cooldownPeriod prevents flapping; every scale event triggers a rebalance that briefly pauses consumption

- Consumers are NOT HTTP services — long-running pull workers. HTTP only for /healthz, /readyz, /metrics

**WATCH OUT  **Adding partitions later is the painful operation: hash(key) % N changes, so meters get remapped and ordering breaks across the boundary. Provision ahead — 12 carries you to roughly 2,500 customers.

### DD2 · Cassandra modelling + tenant isolation (~6 min)

**THESIS  ***customer_id is my tenancy key, so it's my partition key — tenant isolation falls out of the data model rather than being enforced in application code.*

PRIMARY KEY ((customerid, relid, deviceid, month), timestamp)

  WITH CLUSTERING ORDER BY (timestamp DESC)      -- value is a regular column

- Month bucket → ~8,640 rows/partition at 5-min intervals (12 × 24 × 30). Well within limits.

- Why not a DAY bucket → only 288 rows: 30× more partitions than needed, coordinator overhead on every multi-day query

- DESC clustering → "latest reading" is LIMIT 1 with no scan

#### The rule that shapes everything

**ALL partition-key components are required, with equality. **SELECT projects columns; WHERE locates data. Putting device_id only in SELECT gives Cassandra no way to find the rows.

**Therefore: one table per query shape. **Denormalization by access pattern is normal Cassandra modelling, not a smell.

| **Query the UI asks** | **Table** | **Partition key** | **Clustering** |
| --- | --- | --- | --- |
| One meter, time-series chart | readings_by_meter | customer, rel, device, month | timestamp DESC |
| Fleet grid / latest state | device_latest | customer | rel, device |
| Meters at a site | meters_by_relationship | customer, rel | device |
| Long-range trend | readings_hourly / _daily | customer, rel, device, YEAR | bucket DESC |

#### Cardinality — the false friend

- High partition-key cardinality is GOOD in Cassandra. Keys are HASHED to a token — there is no index that grows. More distinct keys = better distribution.

- The Cassandra failure mode is the OPPOSITE: too few distinct values → hot partitions, unbounded growth.

- InfluxDB differs because it maintains an inverted index over tag values — that's why a million meters is a tag explosion there but a non-issue here.

- ~12M partitions/year at ~430 KB each. Cassandra handles billions.

**WATCH OUT  **The real cardinality trap in Cassandra is SECONDARY INDEXES — a query using one fans out to every node. Never index device_id; create another table instead.

**Cross-month range: **month IN ('2026-01','2026-02') = two KNOWN partitions, two targeted reads. That is NOT a scan. Scanning only happens when you don't supply the partition key.

### DD3 · Durability + availability (~4 min)

**THESIS  ***Billing data means no reading can be lost, so the whole path is at-least-once with idempotent writes.*

- **Producer: **acks=all with min.insync.replicas=2 — the write isn't acknowledged until all in-sync replicas have it. Slowest setting, accepted because a lost reading is lost revenue.

- **Consumer: **commits the offset AFTER the Cassandra write succeeds. A crash between the two re-delivers rather than loses.

- **Idempotent writes **keyed on (customer, rel, device, timestamp) — a replay overwrites the same row instead of double-counting.

*That single fact — idempotency — is what lets BOTH sides run at-least-once without corrupting billing data. It ties the producer, the consumer and the storage model into one story.*

- Availability: RF=3 across 2 DCs, LOCAL_QUORUM. Survives a node loss with no downtime; survives a DC loss.

- Under CAP: chose A over C. A briefly stale dashboard is fine; a dashboard that's DOWN is not.

- Production nuance to volunteer: commit in batches (every N messages / few seconds) to cut overhead — the cost is a slightly larger replay window on crash.

### DD4 · Retention, TTL, TWCS + archival (~5 min)

**THESIS  ***At 300M rows/day, hot retention is bounded by node count, not disk cost.*

**The multiplier: **RF 3 × 2 datacenters × ~2 for compaction headroom ≈ 12× your logical data size on disk.

| **Hot retention** | **Rows** | **Logical** | **Disk needed** | **Nodes @ 2 TB** |
| --- | --- | --- | --- | --- |
| 90 days | 27B | 1.4 TB | 16 TB | ~8 |
| 6 months | 55B | 2.7 TB | 33 TB | ~16 |
| 1 year | 110B | 5.5 TB | 66 TB | ~33 |
| 2 years | 219B | 11 TB | 131 TB | ~66 |
| 5 years | 548B | 27 TB | 329 TB | ~164 |

- Constraint is OPERATIONAL, not financial — nodes are kept at 1–2 TB because beyond that repair, compaction and bootstrapping a replacement take unacceptably long.

- Practical ceiling ≈ 1 year. Beyond that, archive.

#### TTL + TWCS — why tombstones never accumulate

compaction = {'class':'TimeWindowCompactionStrategy', 'compaction_window_unit':'DAYS', ...}

default_time_to_live = 31536000   -- 1 year        gc_grace_seconds = 10800

**The one fact everything rests on: **SSTables are immutable, and Cassandra can delete an SSTable only when EVERY row inside it is dead. One live row keeps the whole file.

| **Strategy** | **Groups files by** | **What happens on expiry** |
| --- | --- | --- |
| STCS (default) | size | Mixed ages in one file → live rows keep it alive → tombstones linger for the full TTL span |
| LCS | levels / overlap | Same age-mixing, PLUS heavy rewrite amplification. Worst choice at 300M writes/day. |
| TWCS | WRITE TIME | All rows in a window share the TTL → all expire together → whole SSTable dropped as a FILE DELETE (rm). No compaction, no I/O. |

- TWCS groups by WHEN you wrote, not WHAT you wrote — a day's writes across hundreds of different customers land in the same file simply because they arrived together.

- Tombstones are still CREATED — they just never accumulate and never get read, because the file is dropped before any query touches it.

- Align the compaction window to the partition bucket (monthly) so a month-scoped read touches one window, not thirty. Cost: up to a month of expired data held on disk.

- gc_grace_seconds → lower to a few hours. TTL expiry happens independently on every replica, so there's no delete to propagate. VALID ONLY if you never issue explicit DELETEs on that table.

- Your reads only touch recent months, so the expiring data lives in files no query opens.

**WATCH OUT  **NEVER bulk-delete (DELETE ... WHERE month < X). That generates real range tombstones with real gc_grace cost. TTL + TWCS is the mechanism; deletes are the anti-pattern.

#### Archival to blob — archive at WRITE time, not delete time

- A SECOND consumer group on the same topics writes Parquet to Azure Blob. It never touches Cassandra.

- Why not a batch job reading Cassandra: every archived row would be a coordinator read competing with dashboard serving. The data is already flowing past in Kafka — catching a copy is free.

- Streaming is right HERE (unlike aggregation) because archiving is STATELESS: buffer a batch, write a file, commit. No accumulation, no window-completeness problem.

**Two-stage layout:**

- Stage 1 — streaming archiver writes hourly files, partitioned signal/date/hour (landing zone)

- Stage 2 — monthly compaction job rewrites into signal=/customer=/year=/month= sorted by device+timestamp, ~400 MB per customer-month (curated zone)

- customer and month go in the PATH (partition pruning). device_id stays a COLUMN — a million device folders is the small-file problem; Parquet row-group stats give the filtering anyway.

- Verify row counts per meter-day between Cassandra and blob BEFORE TTL expires, with a backfill path.

- Kafka retention must cover the SLOWEST consumer group — 7 days. Alert on per-group lag.

## 6 · Standby modules — do NOT volunteer

Each is 60–90 seconds. Fire only when the interviewer opens that door.

### Why Cassandra — and the honest concession ★ most likely unprompted question

**Why Cassandra:**

- Optimised for very high write throughput (~300M writes/day), LSM append-only, no read-before-write

- Read pattern is predictable and key-based — one meter, one month, a known partition

- Multi-DC active-active replication — masterless, no failover step, survives a whole DC

- Prioritises availability under CAP, which matches the NFR

- High-cardinality primary keys hash and distribute evenly across nodes

**Why not InfluxDB:**

- OSS version lacks clustering

- Poor fit for very high cardinality — millions of meters cause tag explosion in its inverted index

**Why not ClickHouse (CORRECTED — see warning):**

- Small frequent writes aren't its strength — it wants large batched inserts

- Single-row point lookups by key aren't its strength

- Multi-region active-active replication is more work than Cassandra's

**WATCH OUT  **Your note said "not analytical processing" — that is BACKWARDS and would be a serious slip. ClickHouse IS the analytical engine; that is its whole point, and Dash0 runs it. Never say ClickHouse isn't analytical.

**The honest concession — LEAD with it, don't wait to be pushed:**

- Materialized views into AggregatingMergeTree compute rollups ON INSERT — that replaces my entire agg job, including its scheduling, idempotency and backfill machinery

- Columnar compression is ~5–10× better — that rewrites the storage math, turning a 33-node year into something far smaller

- Ad-hoc GROUP BY across any dimension without pre-modelling — in Cassandra every new query shape is a new table and a new write path

**SAY IT  ***"The Cassandra decision was sound for tenant-scoped serving and write throughput; the gap was ad-hoc aggregation, and that's where I'd revisit it — I'd either run ClickHouse alongside for the analytical side, or reconsider whether two stores are worth it at this volume."*

### Aggregation job

- Scheduled Go job in DMSP. Hourly, incremental — processes only the last window, not whole months.

- Enumerates meters from meters_by_relationship — the raw table CAN'T enumerate, it needs the full partition key.

- Per hour computes: avg, min, max, sum, count, last. Daily is built FROM hourly (24 rows, not 288).

- Rollup tables bucketed by YEAR, not month → a 2-month custom range is ONE partition read. That coarser bucket is deliberate.

- Recompute-and-overwrite → idempotent → late-arriving data self-corrects on the next run with zero special-case code.

**Reading type changes the aggregation — a domain point worth raising unprompted:**

| **Type** | **Example** | **Behaviour** | **Correct rollup** |
| --- | --- | --- | --- |
| Instantaneous | voltage, power, flow rate | goes up and down | avg / min / max — SUM is meaningless |
| Cumulative register | kWh, water volume | only increases | last − first = consumption for the period |

- Cumulative registers ROLL OVER (999999 → 0) and RESET on meter replacement — both look like a huge negative delta. Handle in the agg job, not in ingestion.

**Why batch and not a streaming consumer group:**

- Streaming holds running totals in memory, per meter, per open bucket — state that must be protected

- Pod dies → partial totals gone → replay from an older offset or accept a wrong bucket

- Never knows when a bucket is COMPLETE — a 10:07 reading arriving at 11:03 is a guess either way

- Late data becomes code you must write; with batch it's handled by doing nothing

- Backfill past Kafka retention is impossible — you'd write a batch job anyway

**Why not Flink for rollups: **Flink genuinely solves all of the above (checkpointed state, event-time windows, watermarks, allowed lateness). But the benefit it buys is FRESHNESS, and rollups only serve long-range trends — the live views read raw tables and device_latest and never touch rollups. Against that: a stateful JVM cluster in an all-Go shop, plus you'd still need the batch path for multi-year reports.

**Cassandra has NO automatic rollup or downsampling. **It's a storage engine, not a compute engine. Materialized views only re-key (no GROUP BY). Aggregate functions run at query time on the coordinator. Counters aren't idempotent. That's WHY the agg job exists — a deliberate Cassandra tradeoff, not an oversight.

### Alerting — Flink SLA monitoring ← the section you flagged

**THESIS  ***Flink is a poor fit for the rollups and a good fit for the alerting path — freshness requirement and state requirement together decide batch vs stream.*

**Why Flink here: **must fire in seconds because someone gets paged; detection needs per-entity state and timers; and ABSENCE is not cheaply pollable — a batch job would query a million meters against Cassandra every few minutes and still be minutes late.

**Topology: **a THIRD consumer group on the same topics, own offsets, so it can't affect ingestion or archiving. Then keyBy(meter) into three detectors.

| **#** | **Detector** | **State?** | **How it works** |
| --- | --- | --- | --- |
| 1 | Silent meter | stateful | State = last_seen + alerted_flag + timer. Each reading resets the timer to now+15min. Timer fires unreset → EMIT. The alert is produced by ABSENCE — this is what polling can't do. |
| 2 | Threshold breach | stateless | Value out of range on arrival → emit immediately. |
| 3 | Fleet degraded | windowed | 1-min window counting distinct meters per customer; below ~95% of expected → ONE customer-level alert instead of 1,500 individual ones. |

- Checkpointing snapshots all keyed state AND pending timers to durable storage — a restart resumes watching a million meters rather than forgetting them.

#### Alert idempotency — three layers

- **Flink: emit on TRANSITION, not on condition. **The alerted_flag gives one 'firing' event when it goes silent and one 'resolved' when readings resume — nothing in between. Without it, a meter dead for a week produces ~670 alerts.

- **Kafka: key = device_id + alert_type. **Firing and resolved stay ordered in one partition; log compaction makes the topic hold current alert state. Include a fingerprint hash for replay dedup.

- **Alertmanager fingerprints on labels; OpsGenie dedupes on alias. **Even duplicates that slip through collapse server-side.

- Tradeoff: pure edge-triggering means a LOST notification never retries. Add a low-frequency heartbeat re-emit (30–60 min) while still firing — one notification, but the pipeline stays self-healing.

- Add hysteresis so a flapping meter doesn't fire/resolve repeatedly — require the condition to hold for N checks.

#### Why Kafka between Flink and the alerting stack (not a direct HTTP push)

- Flink's exactly-once covers state and transactional sinks, NOT HTTP side effects — a checkpoint restart would re-push every alert already sent

- Alert storm is the realistic failure mode: a gateway dies and 1,500 meters go silent in one window. The topic absorbs that; 1,500 synchronous HTTP calls would backpressure a stateful streaming job

- You need alert HISTORY for per-customer SLA reports — a second consumer writing to Cassandra. If alerts only existed as HTTP calls into OpsGenie, that history lives in a tool you can't query

- Same Kafka cluster — no new one needed. Alerts are a handful per minute vs 20k/s ingest; 3–6 partitions. Different topics, so no loop.

- Flink DETECTS only. Alertmanager / OpsGenie / Slack / incident.io keep routing, dedup and escalation.

**Location data: **NOT in telemetry. A meter is a fixed installation, so coordinates come from the device REGISTRY at commissioning — which is why absence of readings doesn't block enrichment. Enrich in the NOTIFIER (alert rate) not in Flink (20k/s reading rate). If the lookup fails, send the alert anyway with the device ID — never drop an alert because enrichment failed.

### UI dashboard query ranges

**The rule: range determines granularity. **A chart is ~1,000 pixels wide, so there's no point returning more than ~1,000 points. The server resolves range → granularity → table.

| **UI filter** | **Raw points** | **Served from** |
| --- | --- | --- |
| Last 1 hour (default) | 12 | RAW — so the main screen doesn't depend on the agg job at all |
| Last 3 / 6 / 12 hours | 36 / 72 / 144 | RAW |
| Last 24 hours | 288 | RAW |
| Last 7 days | 2,016 | hourly rollup → 168 points |
| Per month | 8,640 | hourly rollup → 720 points |
| Custom, 2 months | ~17,000 | hourly/daily rollup — bucketed by YEAR, so ONE partition |

- "Latest data" is a SEPARATE query shape. Single meter → clustering DESC + LIMIT 1. Fleet grid → device_latest table, PK customer, upserted on every write: one partition read returns the whole fleet.

- Freshness is bounded by the METERS, not the pipeline — they report every 5 minutes, so "latest" means up to 5 minutes old BY DESIGN. Say this to preempt "is this real-time?"

- Cross-month range = two reads merged. Normal and cheap.

- Cap raw queries at ~7 days at the API and redirect to rollups; push large exports to async.

### Reports

- ALWAYS async: POST /dmsp/reports → 202 + report_id → queue → worker → stream rows to a file in blob → pre-signed time-limited URL → UI polls or user is emailed.

- Why: a fleet-wide multi-year report is thousands of partition reads — cheap individually, far too long for an HTTP request. A report is a batch job with a UI in front of it, not a query.

- STREAM, don't buffer — memory stays flat regardless of range.

- Rate-limit report generation per tenant; it's the heaviest read path you have.

**If reports need RAW data (not rollups): **read the Parquet archive from blob, merging with Cassandra only for the recent hot window. 5 years raw in Cassandra would be ~600B rows / ~100 TB replicated on hot storage read a few times a year — wrong tier. The archive isn't a compromise; it keeps the heaviest read OFF the cluster serving dashboards.

### Rate limiting — token bucket

- TOKEN BUCKET, not leaky bucket: traffic is legitimately bursty, and a leaky bucket smooths to a constant rate — it would reject traffic that's normal for this system.

- Not fixed-window either: a tenant can send a full quota at the end of one window and again at the start of the next = 2× the limit right at your peak.

- At MY edge, not the customer's gateway — that's their code; a buggy gateway replaying a backlog is exactly what this protects against.

- Go: x/time/rate is the standard token bucket but IN-PROCESS — multi-replica means per-pod limits multiply and shift as KEDA scales. Per-tenant needs shared state: go-redis/redis_rate (GCRA), pipelined, fail-open.

- Limit on MESSAGES as well as requests — gateways post batches, so one request with 50k readings passes a request-based limiter untouched. Cap batch size and payload too.

- Return 429 with Retry-After and document it, so a well-behaved gateway backs off.

| **Tier** | **Devices** | **Avg msg/s** | **Burst msg/s** | **Refill** | **Capacity** |
| --- | --- | --- | --- | --- | --- |
| Small | 200 | 0.7 | 3.3 | 5/sec | 400 |
| Standard | 1,500 | 5 | 25 | 20/sec | 2,500 |
| Large | 5,000 | 17 | 83 | 70/sec | 8,000 |
| Enterprise | 10,000 | 33 | 167 | 135/sec | 15,000 |

**Formula (derive any tier live): **refill ≈ (devices ÷ 300) × 4  ·  capacity ≈ devices × 1.5. Refill gives ~4× headroom over sustained; capacity absorbs a full reporting cycle plus 50%.

- Global ceiling ~40k msg/s against a ~20k peak.

- Start LOG-ONLY for a week and set thresholds from observed p99 × 2 — shipping hard limits on estimates is how you page yourself at 3am for a legitimate customer.

### Growth — onboarding 300 more customers (700 → 1,000)

1.5M devices · ~5,000/s average · ~25,000/s peak · ~432M/day — a 43% increase.

| **Layer** | **Current** | **After 1,000** | **Change?** |
| --- | --- | --- | --- |
| Per-tenant rate limit | 20/sec, cap 2,500 | same | NO — tenant size is constant |
| Global rate limit | 40k/sec | 60k/sec | YES — bump |
| Kafka partitions | 12 | 12 | NO |
| KEDA maxReplicas | 12 | 12 | NO |
| Consumers at peak | ~4 | ~5 | automatic (KEDA) |
| Kafka storage | 423 GB/broker | 605 GB/broker | watch |
| Cassandra nodes (90d) | ~8 | ~12 | YES — add nodes |
| Cassandra nodes (1yr) | ~33 | ~47 | YES — add nodes |
| Flink keyed state | 1.05M meters | 1.5M meters | check heap |

**The insight to lead with: **per-tenant limits DON'T change, because tenant size is constant. Only aggregate things move. That's why tenant-provisioned limits are the right design — onboarding is a config row, not a re-tuning exercise.

**Kafka headroom: **25k/s at ~200 bytes ≈ 5 MB/s. Twelve partitions at a conservative 10 MB/s each ≈ 120 MB/s ≈ 600,000 msg/s. You're using ~4% — Kafka could absorb 10× growth without touching partition count. The binding constraint is downstream: consumer throughput into Cassandra, and Cassandra capacity.

## 7 · Key numbers — memorise these

| **Metric** | **Value** | **Where it comes from** |
| --- | --- | --- |
| Tenants (vendors) | 700 | given |
| End customers per vendor | ~1,500 | given |
| Total meters | ~1.05M | 700 × 1,500 |
| Reporting interval | 5 min | given |
| Average throughput | ~3,500 msg/s | 1.05M ÷ 300s |
| Peak throughput | ~17,500–20,000 msg/s | 1.05M ÷ 60s burst window |
| Peak : average | ~5× | aligned clock boundaries |
| Daily volume | ~300M messages | 1.05M × 288 |
| Kafka topics | 3 | meters / metrics / kpi |
| Partitions per topic | 12 | peak ÷ ~5k per consumer + headroom |
| Brokers | 3 | RF=3 minimum |
| Consumers at peak | ~4 | 17,500 ÷ ~5,000 |
| Rows per Cassandra partition | ~8,640 | 12/hr × 24 × 30 (month bucket) |
| Partitions per year | ~12M | 1.05M meters × 12 months |
| Partition size | ~430 KB | 8,640 rows × ~50 bytes |
| Storage multiplier | ~12× | RF3 × 2 DC × 2 compaction headroom |
| 1 year hot | ~66 TB, ~33 nodes | 110B rows × 50 bytes × 12 |
| 90 days hot | ~16 TB, ~8 nodes | 27B rows × 50 bytes × 12 |
| Per-tenant average | 5 msg/s | 1,500 ÷ 300s |
| Per-tenant burst | 25 msg/s | 1,500 ÷ 60s |
| Per-tenant limit | refill 20/s, cap 2,500 | 4× headroom, 1.5× cycle |
| Kafka retention | 7 days | must cover slowest consumer group |
| Parquet file target | ~400 MB / customer-month | 13M rows compressed |

## 8 · Delivery reminders

Personal patterns observed across practice passes — the content is there, these are the gaps that cost points.

- Say the COUNT before any list. "Functionally, four capabilities: one… two…" Filler drops immediately when you're walking a numbered rail.

- Anchor every NFR to a number or a named mechanism. Fluency measurably improves in the numbers sections and degrades in unanchored prose — so scaffold the prose.

- COMMIT on component names. No "or something", no "whatever you call it". Uncertainty markers appear exactly where a phrase isn't rehearsed — so rehearse those phrases.

- Enunciate: KEDA, DMSP, Entra, servsup, bearer token, TWCS.

- Don't apologise for the HLD. "It doesn't yet address the NFRs" — never "I know a lot of issues in this design."

- Don't stop at the API in step 3 — say the words "write dataflow:" and "read dataflow:" out loud and walk both to Cassandra.

- Land the entities closer: the hierarchy IS the partition key.

- Never let the burst number fall out. "~3,500 average, ~17,500 peak, because meters fire on aligned boundaries."

- Never let tenant isolation fall out of the NFR list.

- When a deep-dive detail leaks early it comes out shakier — hold it, then deliver it deliberately.

- End every dive with a hook: "that's X — I can go into Y or Z, which is more useful?"
