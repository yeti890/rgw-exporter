# Architecture — RGW Usage Exporter

This document describes the high-level architecture of **RGW Usage Exporter**.

The exporter is designed for **large Ceph RGW clusters**, focusing on predictability, low overhead and scalability.

---

## High-level overview
```text
+----------------------------+
|         Prometheus         |
|      (scrape /metrics)     |
+-------------+--------------+
              |
              v
+----------------------------+
|     RGW Usage Exporter     |
|                            |
| - usage collector          |
| - buckets collector        |
| - users collector          |
|                            |
+-------------+--------------+
              |
              v
+----------------------------+
|     Ceph RGW Admin API     |
| (radosgw-admin equivalent) |
+----------------------------+
```

```mermaid
flowchart TB
  P[Prometheus<br/>(scrape /metrics)] --> E[RGW Usage Exporter<br/>- usage collector<br/>- buckets collector<br/>- users collector]
  E --> A[Ceph RGW Admin API<br/>(radosgw-admin equivalent)]
```

---

## Core principles

### 1. Pull-based metrics
- Prometheus scrapes `/metrics`
- Exporter maintains internal state updated by background collectors
- No heavy work is done during scrape itself

---

### 2. Background collectors

The exporter runs multiple collectors on configurable intervals:

| Collector | Purpose |
|---------|---------|
| Usage collector | RGW ops and traffic statistics |
| Buckets collector | Bucket stats, quotas, shards |
| Users collector | User info and quotas |

Each collector:
- runs independently,
- measures its own execution time,
- updates shared in-memory structures.

---

### 3. No PromQL joins

**Bucket owner (`uid`) is injected into all bucket-level metrics at collection time.**

This removes the need for:
- `group_left`
- metadata joins
- cross-metric lookups in Grafana

Result:
- predictable PromQL performance,
- fast Grafana tables,
- stable memory usage.

---

### 4. Minimal shared state

- Data is stored in compact Go structs
- No metric arrays rebuilt on every scrape
- Prometheus descriptors are static
- Only numeric values are updated

---

### 5. Label strategy

Labels are intentionally limited and consistent:

- `region`
- `cluster`
- `endpoint`
- `uid`
- `bucket`
- `category`

This provides:
- sufficient filtering power,
- controlled cardinality,
- predictable TSDB growth.

---

## Failure model

- RGW API failures affect only the current collection cycle
- Last known values remain exposed
- Collector durations are always exported for observability

---

## Deployment model

- Single exporter instance per RGW endpoint
- Stateless process
- Horizontal scaling supported (per RGW / per zone)

---

## Summary

The architecture favors:
- stability over completeness,
- predictable cost over dynamic joins,
- exporter-side computation over PromQL complexity.

This makes the exporter suitable for **long-term production monitoring**.