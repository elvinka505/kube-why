---
title: HPA can't get metrics (FailedGetResourceMetric)
aliases: [failedgetresourcemetric, hpa-no-metrics, unabletofetchmetrics]
category: hpa
---

# HPA can't get metrics (FailedGetResourceMetric)

## What it means
A `HorizontalPodAutoscaler` needs current metrics (CPU, memory, or a custom
metric) to decide whether to scale, and it can't fetch them. Until it can,
the HPA takes no scaling action at all, it doesn't fail open or closed, it
just does nothing.

## Common causes
1. `metrics-server` isn't installed in the cluster, or is installed but
   unhealthy, this is required for basic CPU/memory-based autoscaling and
   isn't always installed by default.
2. The target pods don't have `resources.requests` set for the metric being
   scaled on, CPU/memory autoscaling is calculated as a percentage of
   requests, no request means no percentage to calculate.
3. For custom/external metrics, the metrics adapter (Prometheus adapter,
   cloud provider's adapter) is down or the metric name doesn't match what
   the HPA is querying for.
4. RBAC permissions are missing for the HPA controller to read the metrics
   API.

## Check it
```
kubectl describe hpa <name>
kubectl top pods
kubectl get apiservices | grep metrics
```
`describe hpa` shows the exact condition and reason. `kubectl top pods`
failing with the same kind of error confirms it's a metrics-server problem
generally, not something specific to the HPA.

## Fix it
- Install or fix `metrics-server` if `kubectl top` doesn't work at all.
- Add `resources.requests` to the target workload if scaling on CPU/memory,
  it's a required input, not optional metadata.
- For custom metrics, verify the metrics adapter is healthy and that the
  metric name in the HPA spec matches exactly what the adapter exposes.
- Check the HPA controller's RBAC if metrics-server and the workload both
  look correct and it's still failing.

## Related
- HPAStuckAtMinReplicas
- OOMKilled
