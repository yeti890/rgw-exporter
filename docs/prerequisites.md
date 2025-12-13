# Prerequisites — RGW Usage Exporter

This document describes the required Ceph RGW configuration and permissions needed to run **RGW Usage Exporter** in production environments.

---

## Ceph RGW requirements

A working **Ceph cluster with RADOS Gateway (RGW)** is required.

The exporter relies on RGW usage statistics, which are **disabled by default** and must be explicitly enabled.

---

## RGW usage logging configuration

The following example shows a typical production configuration.

```ini
rgw enable usage log = true
rgw usage log flush threshold = 1024
rgw usage log tick interval = 30
rgw usage max shards = 32
rgw usage max user shards = 8
````

### Notes

* These values are examples and should be tuned according to:

  * number of buckets,
  * number of users,
  * RGW load characteristics.
* Incorrect tuning may result in large objects in the RGW usage log pool.
* Always review the official Ceph documentation before applying changes.

---

## RGW Admin API

The exporter uses the **RGW Admin API**.

Ensure that:

* the admin entry point is enabled,
* the Admin API is accessible.

```ini
rgw admin entry = admin
rgw enable apis = s3, admin
```

---

## RGW exporter user

It is strongly recommended to use a **dedicated RGW user with read-only permissions**.

### Create user

```bash
radosgw-admin user create \
  --uid="rgw-exporter" \
  --display-name="RGW Usage Exporter"
```

### Assign capabilities

```bash
radosgw-admin caps add \
  --uid="rgw-exporter" \
  --caps="metadata=read;usage=read;info=read;buckets=read;users=read"
```

### Required permissions overview

These capabilities allow the exporter to:

* list users and buckets,
* read bucket metadata,
* collect usage statistics,
* read quota configuration.

No write permissions are required.

---

## Load balancer considerations

When using a load balancer (e.g. HAProxy) in front of RGW:

* usage queries may take longer on clusters with:

  * large number of buckets,
  * large number of users.
* backend timeouts must be configured accordingly.

For HAProxy, ensure that `timeout server` is sufficiently high to allow RGW usage queries to complete under peak load.

---

## Security considerations

* The exporter operates in **read-only mode**.
* A dedicated service account with minimal privileges is recommended.
* The exporter does not modify cluster state.

---