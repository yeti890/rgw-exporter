# Multi-endpoint and multi-realm clusters

Ceph deployments may expose **multiple RGW endpoints**, for example:

- multi-site configurations,
- multiple RGW zones,
- multiple RGW realms.

RGW Usage Exporter fully supports such setups.

## Deployment model

The exporter is **endpoint-scoped**.

For clusters with multiple RGW endpoints, the recommended approach is:

- deploy **one exporter instance per RGW endpoint**;
- configure each instance with:
  - its own `RGW_ENDPOINT` (Admin API),
  - corresponding `PUB_ENDPOINT` (public S3 endpoint);
- scrape all exporter instances from Prometheus.

This approach provides:
- clear isolation between endpoints,
- predictable performance,
- simple horizontal scaling.

## RGW users and realms

RGW users are **realm-scoped**.

If your cluster uses multiple realms, a dedicated exporter user must be created **in each realm**.

It is recommended to use the same user name (e.g. `rgw-exporter`) in all realms for consistency.

## Creating exporter user in a specific realm

Create the user in a given realm:

```bash
radosgw-admin user create \
  --uid="rgw-exporter" \
  --display-name="RGW Usage Exporter" \
  --rgw-realm=<realm-name>
````

Assign required capabilities:

```bash
radosgw-admin caps add \
  --uid="rgw-exporter" \
  --caps="metadata=read;usage=read;info=read;buckets=read;users=read" \
  --rgw-realm=<realm-name>
```

Repeat these steps for **each RGW realm** used by the exporter.


## Prometheus considerations

In multi-endpoint setups:

* each exporter instance exposes its own `/metrics`;
* metrics are distinguished by the `endpoint` label;
* Prometheus can aggregate metrics across endpoints using standard label matching.

This allows:

* per-endpoint dashboards,
* per-cluster aggregation,
* flexible multi-region views.

## Summary

For multi-endpoint Ceph RGW clusters:

* deploy one exporter per endpoint,
* create exporter user per realm,
* keep configuration explicit and isolated.

This model scales well and is suitable for large production environments.
