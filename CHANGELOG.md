# Changelog

All notable changes to this project will be documented in this file.

The format is based on **Keep a Changelog**, and this project adheres to **Semantic Versioning (SemVer)**.

## [1.1.0] - 2025-12-13

### Changed
- Configuration is now **ENV-only** (recommended for container/systemd deployments).
- Collector ticker loops refactored for clarity (explicit immediate run + periodic ticks).

### Removed
- Removed optional YAML config file support and the `-c` flag (to reduce complexity and avoid ENV/file precedence ambiguity).

### Fixed
- Fixed a potential data race when exporting collector duration metrics (durations are now read under mutex).

### Security
- Dependency cleanup and version pinning to reduce supply-chain risk and improve vulnerability scan results.

## [1.0.0] - 2025-12-12

### Added
- Prometheus exporter for Ceph RGW (RADOS Gateway) usage statistics via RGW Admin API.
- Background collectors with configurable intervals:
  - Usage collector (ops/bytes by category, bucket, uid)
  - Buckets collector (bucket stats, quotas, shards)
  - Users collector (user info, quotas)
- Bucket-level metrics (with `uid` label included in bucket metrics for fast filtering):
  - logical size, actual size, objects
  - quota enabled/max size/max objects
  - number of shards and objects per shard
  - quota usage percent (size-based)
- User-level metrics:
  - suspended flag
  - user quota (enabled/max size/max objects)
  - user bucket quota (enabled/max size/max objects)
  - user used size in bytes (sum of owned bucket sizes)
  - user quota usage percent (size-based)
  - user buckets total
- Cluster aggregates:
  - buckets total, users total, objects total
  - total logical/actual bucket size
  - total configured bucket quotas size
  - total configured user quotas size
- Collector duration metrics for buckets/usage/users collectors.
- Configuration via environment variables and optional YAML config file (`-c`).
- HTTP endpoint `/metrics` for Prometheus scraping.

### Changed
- N/A (initial release)

### Fixed
- N/A (initial release)

### Security
- Supports TLS verification toggle (`INSECURE`) and request timeouts.
