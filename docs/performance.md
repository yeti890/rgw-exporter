# Performance considerations

This document explains why **RGW Usage Exporter** performs well on large Ceph RGW clusters.

---

## Design goals

- Support **100k+ buckets**
- Handle **millions of objects**
- Avoid Prometheus/Grafana performance degradation
- Maintain low and predictable memory usage

---

## Key performance decisions

### 1. No PromQL joins

All expensive joins are avoided by design.

- Bucket metrics already contain `uid`
- User usage is pre-aggregated inside the exporter

This shifts work from:
- Prometheus query time  
to:
- exporter collection time (controlled, predictable)

---

### 2. Background collection model

- Collection is done on fixed intervals
- `/metrics` endpoint is lightweight
- Scrapes do not trigger RGW API calls

Result:
- scrape latency stays low,
- Prometheus timeouts avoided.

---

### 3. Minimal allocations

- Static metric descriptors
- Reused data structures
- No per-scrape metric reconstruction

This significantly reduces:
- GC pressure
- memory fragmentation
- CPU spikes

---

### 4. Controlled cardinality

Cardinality drivers are well understood:
- users
- buckets
- categories

The exporter:
- avoids dynamic labels,
- avoids free-form metadata,
- avoids tenant-level explosion.

---

### 5. Aggregations done once

Expensive operations (sums, totals, percentages) are:
- computed once per collection cycle,
- exported as ready-to-use metrics.

This makes dashboards cheap and fast.

---

## Observability of the exporter itself

Collector execution time is exported:

- `radosgw_usage_collector_usage_duration_seconds`
- `radosgw_usage_collector_buckets_duration_seconds`
- `radosgw_usage_collector_users_duration_seconds`

This allows:
- alerting on slow RGW responses,
- capacity planning for exporter instances.

---

## Practical results

In production environments this design allows:

- stable scrape times even with high cardinality,
- responsive Grafana dashboards,
- long retention periods in Prometheus/VictoriaMetrics.

---

## Summary

The exporter is optimized for:
- predictable performance,
- operational safety,
- long-term monitoring at scale.

It intentionally trades flexibility for stability.
