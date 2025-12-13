![Go](https://img.shields.io/badge/go-1.22+-00ADD8?logo=go&logoColor=white)
![Prometheus](https://img.shields.io/badge/Prometheus-exporter-E6522C?logo=prometheus&logoColor=white)
![Ceph](https://img.shields.io/badge/Ceph-RGW-2A5ADA)
![Grafana](https://img.shields.io/badge/Grafana-dashboard-F46800?logo=grafana&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green)
![Container](https://img.shields.io/badge/container-Podman%20%7C%20Docker-blue)

<p align="center">
  <img src="./rgw-exporter-mini.png" alt="RGW Exporter logo" width="180" />
</p>

# RGW Usage Exporter for Ceph

High-performance Prometheus exporter for **Ceph RGW (RADOS Gateway / S3)** built for **large production clusters**.

Collects RGW usage, bucket stats, users and quotas via **RGW Admin API** — optimized to avoid expensive PromQL joins.  
**Bucket metrics include `uid`**, so Grafana tables stay fast even with high cardinality.

## Features

- 🚀 **Built for scale**
  - designed for 100k+ buckets / millions of objects
  - low allocations and GC pressure
  - no PromQL `join` / `group_left`

- 📦 **Bucket metrics (with uid)**
  - logical size, actual size, objects
  - quotas (enabled / max size / max objects)
  - index shards and objects per shard
  - quota usage percent

- 👤 **User metrics**
  - suspended flag
  - user quota and user bucket quota
  - used size (sum of owned bucket sizes)
  - quota usage percent
  - buckets per user

- 📊 **Cluster aggregates**
  - buckets/users/objects totals
  - total logical/actual size
  - total configured user/bucket quotas (oversell analysis)

- 📈 **Grafana-ready**
  - dashboard JSON (import-ready)

## Prerequisites

A working **Ceph RGW (RADOS Gateway)** setup is required.

### Create a dedicated RGW user (recommended)

The exporter requires a read-only RGW user to access the Admin API.

```bash
radosgw-admin user create \
  --uid="rgw-exporter" \
  --display-name="RGW Usage Exporter"
```
```bash
radosgw-admin caps add \
  --uid="rgw-exporter" \
  --caps="metadata=read;usage=read;info=read;buckets=read;users=read"
```

### RGW usage logging

RGW usage logging must be enabled to collect statistics.

See detailed configuration and production recommendations here:
**[docs/prerequisites.md](docs/prerequisites.md)**

## Multi-endpoint (multi-RGW) clusters

The exporter supports Ceph clusters with **multiple RGW endpoints** (e.g. multi-site, multi-zone or multi-realm setups).

The recommended deployment model is:

- run **one exporter instance per RGW endpoint**;
- configure each instance with its own `RGW_ENDPOINT` and `PUB_ENDPOINT`;
- use consistent labels (`cluster`, `region`, `endpoint`) for aggregation in Prometheus.

For security and isolation reasons, a dedicated `rgw-exporter` user must exist **in each RGW realm** used by the exporter.

See detailed multi-endpoint setup instructions here:   **[docs/multisite.md](docs/multisite.md)**

## Quick start (Docker)
```bash
docker run -d \
  --name rgw-exporter \
  --network host \
  -e ACCESS_KEY=xxxx \
  -e SECRET_KEY=yyyy \
  -e RGW_ENDPOINT=https://rgw-admin:443 \
  -e PUB_ENDPOINT=s3.example.com \
  -e REGION=DC1 \
  -e CLUSTER_NAME=SRV-01 \
  -e USERS_COLLECTOR_ENABLE=true \
  -e LISTEN_IP=0.0.0.0 \
  -e LISTEN_PORT=9240 \
  -e USAGE_COLLECTOR_INTERVAL=30 \
  -e BUCKETS_COLLECTOR_INTERVAL=300 \
  -e USERS_COLLECTOR_INTERVAL=600 \
  -e RGW_CONNECTION_TIMEOUT=300 \
  -e START_DELAY=0 \
  -e INSECURE=true \
  -e SKIP_WITHOUT_BUCKET=false \
  docker.io/yeti89/rgw-exporter:latest
```

## Quick start (Podman)
```bash
podman run -d \
  --name rgw-exporter \
  --network host \
  -e ACCESS_KEY=xxxx \
  -e SECRET_KEY=yyyy \
  -e RGW_ENDPOINT=https://rgw-admin:443 \
  -e PUB_ENDPOINT=s3.example.com \
  -e REGION=DC1 \
  -e CLUSTER_NAME=SRV-01 \
  -e USERS_COLLECTOR_ENABLE=true \
  -e LISTEN_IP=0.0.0.0 \
  -e LISTEN_PORT=9240 \
  -e USAGE_COLLECTOR_INTERVAL=30 \
  -e BUCKETS_COLLECTOR_INTERVAL=300 \
  -e USERS_COLLECTOR_INTERVAL=600 \
  -e RGW_CONNECTION_TIMEOUT=300 \
  -e START_DELAY=0 \
  -e INSECURE=true \
  -e SKIP_WITHOUT_BUCKET=false \
  docker.io/yeti89/rgw-exporter:latest
```

## Configuration (ENV)
| Variable                     | Description                                   |
| ---------------------------- | --------------------------------------------- |
| `ACCESS_KEY`                 | RGW admin access key                          |
| `SECRET_KEY`                 | RGW admin secret key                          |
| `RGW_ENDPOINT`               | Internal RGW Admin endpoint                   |
| `PUB_ENDPOINT`               | Public S3 endpoint (used as label `endpoint`) |
| `REGION`                     | Region/DC/zone label                          |
| `CLUSTER_NAME`               | Cluster label                                 |
| `LISTEN_IP`                  | Listen IP for `/metrics`                      |
| `LISTEN_PORT`                | Listen port (default `9240`)                  |
| `USAGE_COLLECTOR_INTERVAL`   | Usage collection interval (sec)               |
| `BUCKETS_COLLECTOR_INTERVAL` | Buckets collection interval (sec)             |
| `USERS_COLLECTOR_INTERVAL`   | Users collection interval (sec)               |
| `USERS_COLLECTOR_ENABLE`     | `true` / `false`                              |
| `RGW_CONNECTION_TIMEOUT`     | RGW request timeout                           |
| `START_DELAY`                | Startup delay                                 |
| `INSECURE`                   | Disable TLS verification                      |
| `SKIP_WITHOUT_BUCKET`        | Skip entries without bucket                   |

## Metrics
See full metrics reference: [docs/metrics.md](docs/metrics.md)

Metrics endpoint:
```bash
curl http://<host>:9240/metrics
```

## Grafana dashboard
This exporter is designed to be used with a dedicated Grafana dashboard.

The dashboard is currently being polished and prepared for public release on [**grafana.com**](https://grafana.com/grafana/dashboards/).
It is already used in production and will be published after final cleanup and documentation.

## PromQL examples

Buckets owned by a user:
```
radosgw_usage_bucket_size{uid="user1"}
```
Top users by used size:
```
topk(10, radosgw_usage_user_used_size_bytes)
```
Buckets with too many objects per shard:
```
radosgw_usage_bucket_objects_per_shard > 500000
```

## License
 - [MIT](./LICENSE)

## Contributing
Issues and PRs are welcome.

<p align="center">
  <img src="./logo-rgw-exporter.png" alt="RGW Exporter logo" width="500" />
</p>