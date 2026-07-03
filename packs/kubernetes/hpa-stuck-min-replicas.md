---
title: HPA stuck at minimum replicas under load
aliases: [hpastuckatminreplicas, hpa-not-scaling]
category: hpa
---

# HPA stuck at minimum replicas under load

## What it means
The app is clearly under load, but the `HorizontalPodAutoscaler` isn't adding
pods. Unlike FailedGetResourceMetric, metrics are working here, the HPA is
seeing data, it's just not deciding to scale up.

## Common causes
1. The metric that's actually under pressure (memory, request latency, queue
   depth) isn't the metric the HPA is configured to watch (often it's set to
   watch CPU, but the real bottleneck is something else entirely).
2. `maxReplicas` is already reached, the HPA is scaling correctly, it's just
   hit its ceiling.
3. A stabilization window (`behavior.scaleUp.stabilizationWindowSeconds`) is
   delaying the scale-up longer than expected, waiting for a sustained trend
   rather than reacting to a spike.
4. The target average utilization threshold is set high enough that current
   load, while real, hasn't technically crossed it yet.

## Check it
```
kubectl describe hpa <name>
kubectl top pods -l <selector>
```
Compare the actual current metric value shown in `describe hpa` against the
target threshold, if current is below target, the HPA is behaving correctly,
your assumption about what's under load is what needs revisiting.

## Fix it
- If the wrong metric is being watched, either add the right metric (custom
  or external) or lower the threshold on the existing one to something that
  actually reflects real pressure.
- If `maxReplicas` is the ceiling, raise it, assuming the cluster has room to
  schedule more pods.
- If a stabilization window is the cause and the delay isn't acceptable for
  this workload, shorten it, understanding that makes scaling more reactive
  to short spikes too.
- Consider scaling on more than one metric if a single metric doesn't reliably
  represent load for this workload.

## Related
- HPAFailedGetResourceMetric
- Pending
