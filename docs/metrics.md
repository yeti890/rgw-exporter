# Metrics reference — RGW Usage Exporter

This document describes all Prometheus metrics exposed by **RGW Usage Exporter for Ceph**.

All metrics are collected via **Ceph RGW Admin API** and are optimized for large-scale production clusters.

---

## Common labels

Most metrics include the following labels:

| Label | Description |
|------|------------|
| `region` | Region / zone name |
| `cluster` | Ceph cluster name |
| `endpoint` | Public S3 endpoint |
| `uid` | RGW user ID |
| `bucket` | Bucket name |
| `category` | RGW operation category (GET, PUT, LIST, etc.) |

---

## Usage metrics (RGW operations)

### `radosgw_usage_ops_total`
Total number of RGW requests.

Labels: {region, cluster, endpoint, uid, bucket, category}

Type: `counter`

---

### `radosgw_usage_successful_ops_total`
Number of successful RGW requests.

Labels: {region, cluster, endpoint, uid, bucket, category}

Type: `counter`

---

### `radosgw_usage_sent_bytes_total`
Total bytes sent by RGW to clients.

Labels: {region, cluster, endpoint, uid, bucket, category}

Type: `counter`

---

### `radosgw_usage_received_bytes_total`
Total bytes received by RGW from clients.

Labels: {region, cluster, endpoint, uid, bucket, category}

Type: `counter`

---

## Bucket-level metrics

> All bucket-level metrics **include `uid` label** (bucket owner),  
> so no PromQL joins are required.

---

### `radosgw_usage_bucket_size`
Logical bucket size (sum of object sizes).

Labels: {region, cluster, endpoint, bucket, uid}

Type: `gauge`  
Unit: `bytes`

---

### `radosgw_usage_bucket_actual_size`
Actual on-disk bucket size.

Labels: {region, cluster, endpoint, bucket, uid}

Type: `gauge`  
Unit: `bytes`

---

### `radosgw_usage_bucket_objects`
Number of objects in the bucket.

Labels: {region, cluster, endpoint, bucket, uid}

Type: `gauge`

---

### `radosgw_usage_bucket_num_shards`
Number of index shards for the bucket.

Labels: {region, cluster, endpoint, bucket, uid}

Type: `gauge`

---

### `radosgw_usage_bucket_objects_per_shard`
Average number of objects per shard.

Labels: {region, cluster, endpoint, bucket, uid}

Type: `gauge`

---

## Bucket quota metrics

### `radosgw_usage_bucket_quota_enabled`
Bucket quota enabled flag.

Values:
- `1` — enabled
- `0` — disabled

Labels: {region, cluster, endpoint, bucket, uid}

Type: `gauge`

---

### `radosgw_usage_bucket_quota_size`
Bucket quota maximum size.

Labels: {region, cluster, endpoint, bucket, uid}

Type: `gauge`  
Unit: `bytes`

---

### `radosgw_usage_bucket_quota_objects`
Bucket quota maximum number of objects.

Labels: {region, cluster, endpoint, bucket, uid}

Type: `gauge`

---

### `radosgw_usage_bucket_quota_usage_percent`
Bucket quota usage percentage (size-based).

Labels: {region, cluster, endpoint, bucket, uid}

Type: `gauge`  
Unit: `percent`

---

## User-level metrics

### `radosgw_usage_user_suspended`
User suspended flag.

Values:
- `1` — suspended
- `0` — active

Labels: {region, cluster, endpoint, uid}

Type: `gauge`

---

### `radosgw_usage_user_quota_enabled`
User quota enabled flag.

Labels: {region, cluster, endpoint, uid}

Type: `gauge`

---

### `radosgw_usage_user_quota_size_bytes`
User quota maximum size.

Labels: {region, cluster, endpoint, uid}

Type: `gauge`  
Unit: `bytes`

---

### `radosgw_usage_user_quota_objects`
User quota maximum number of objects.

Labels: {region, cluster, endpoint, uid}

Type: `gauge`

---

### `radosgw_usage_user_bucket_quota_enabled`
User bucket quota enabled flag.

Labels: {region, cluster, endpoint, uid}

Type: `gauge`

---

### `radosgw_usage_user_bucket_quota_size_bytes`
User bucket quota maximum size.

Labels: {region, cluster, endpoint, uid}

Type: `gauge`  
Unit: `bytes`

---

### `radosgw_usage_user_bucket_quota_objects`
User bucket quota maximum number of objects.

Labels: {region, cluster, endpoint, uid}

Type: `gauge`

---

### `radosgw_usage_user_used_size_bytes`
Total logical size of all buckets owned by the user.

Labels: {region, cluster, endpoint, uid}

Type: `gauge`  
Unit: `bytes`

---

### `radosgw_usage_user_quota_usage_percent`
User quota usage percentage (size-based).

Labels: {region, cluster, endpoint, uid}

Type: `gauge`  
Unit: `percent`

---

## Cluster-level aggregate metrics

### `radosgw_usage_buckets_total`
Total number of buckets in the cluster.

Labels: {region, cluster, endpoint}

Type: `gauge`

---

### `radosgw_usage_users_total`
Total number of users in the cluster.

Labels: {region, cluster, endpoint}

Type: `gauge`

---

### `radosgw_usage_objects_total`
Total number of objects in the cluster.

Labels: {region, cluster, endpoint}

Type: `gauge`

---

### `radosgw_usage_buckets_size_total_bytes`
Total logical size of all buckets.

Labels: {region, cluster, endpoint}

Type: `gauge`  
Unit: `bytes`

---

### `radosgw_usage_buckets_actual_size_total_bytes`
Total actual on-disk size of all buckets.

Labels: {region, cluster, endpoint}

Type: `gauge`  
Unit: `bytes`

---

### `radosgw_usage_bucket_quotas_size_total_bytes`
Sum of all configured bucket quotas.

Labels: {region, cluster, endpoint}

Type: `gauge`  
Unit: `bytes`

---

### `radosgw_usage_user_quotas_size_total_bytes`
Sum of all configured user quotas.

Labels: {region, cluster, endpoint}

Type: `gauge`  
Unit: `bytes`

---

## Collector performance metrics

### `radosgw_usage_collector_usage_duration_seconds`
Duration of usage collector execution.

Labels: {region, cluster, endpoint}

Type: `gauge`  
Unit: `seconds`

---

### `radosgw_usage_collector_buckets_duration_seconds`
Duration of buckets collector execution.

Labels: {region, cluster, endpoint}

Type: `gauge`  
Unit: `seconds`

---

### `radosgw_usage_collector_users_duration_seconds`
Duration of users collector execution.

Labels: {region, cluster, endpoint}

Type: `gauge`  
Unit: `seconds`

---

## Notes

- Metrics are designed to avoid high-cost PromQL joins.
- All bucket metrics include `uid` for direct filtering.
- Suitable for large Ceph RGW clusters with high cardinality.