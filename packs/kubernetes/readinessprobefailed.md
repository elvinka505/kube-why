---
title: ReadinessProbeFailed
aliases: [readinessprobefailed, unhealthy, readiness-probe-failed]
category: pod
---

# ReadinessProbeFailed / Unhealthy

## What it means
The container is running, but Kubernetes doesn't consider it ready to receive
traffic because its readiness probe is failing. The pod stays out of Service
endpoints until the probe passes. This is different from a liveness probe
failure, a failed readiness probe doesn't restart the container, it just keeps
it out of rotation.

## Common causes
1. The probe checks the wrong port or path (common after a port change that
   wasn't updated in the probe config).
2. The app takes longer to become ready than `initialDelaySeconds` allows, so
   the first several checks fail before it's actually up.
3. The endpoint the probe hits depends on something not ready yet (database
   connection, cache warmup) and returns a non-2xx status until that finishes.
4. The probe's timeout is too short for a slow but healthy response, especially
   under load.

## Check it
```
kubectl describe pod <pod>
```
Look at the readiness probe's configured port/path in the pod spec and compare
it against what the app actually listens on. Then check Events for the specific
failure, `connection refused` points to a port/path mismatch, a timeout points
to a slow endpoint.
```
kubectl exec <pod> -- curl -sv localhost:<port><path>
```
Running the probe's exact check by hand, from inside the container, tells you
immediately whether it's a config mismatch or a genuinely slow/broken endpoint.

## Fix it
- Port/path mismatch: fix the probe config to match the app.
- Slow startup: increase `initialDelaySeconds`, or better, use a `startupProbe`
  so readiness/liveness only kick in once startup is confirmed complete.
- Dependency not ready: either make the readiness endpoint not depend on
  non-critical dependencies, or accept the pod being out of rotation until it
  genuinely can serve traffic, that's the probe doing its job correctly.
- Slow responses under load: raise the probe's `timeoutSeconds`, or investigate
  why the endpoint is slow in the first place.

## Related
- CrashLoopBackOff
- ProgressDeadlineExceeded
