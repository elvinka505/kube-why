---
title: OOMKilled
aliases: [oomkilled, oom-killed, oom]
category: pod
---

# OOMKilled

## What it means
The container tried to use more memory than its limit allowed, so the kernel's
OOM killer terminated it. This is not a Kubernetes decision, it's the Linux
kernel enforcing the cgroup memory limit you set (or that a default/LimitRange
set for you).

## Common causes
1. `resources.limits.memory` is set lower than what the app actually needs under
   real load, not just at idle.
2. No memory limit set at all, and a `LimitRange` in the namespace applies a
   default that's too small.
3. A memory leak, the app's usage grows over time until it hits the ceiling.
4. A traffic spike or batch job pushes usage past a limit that was fine on
   average.

## Check it
```
kubectl describe pod <pod>
```
Look for `Last State: Terminated, Reason: OOMKilled` under the container status.
Then compare requested/actual usage:
```
kubectl top pod <pod>
```
If you have metrics history (Prometheus, etc.), graph memory usage over time
leading up to the kill, a slow climb points to a leak, a sudden spike points to
load.

## Fix it
- Quick fix: raise `resources.limits.memory` to a realistic number based on
  observed usage, not a guess.
- Better fix: profile the app under production-like load to find its actual
  working set, then set requests/limits close to that with headroom.
- If it's a leak, this buys you time but doesn't fix the underlying bug, treat
  the limit increase as a mitigation, not a solution.
- Set `requests.memory` close to typical usage too, not just the limit, so the
  scheduler places the pod on a node that can actually sustain it.

## Related
- CrashLoopBackOff
- Evicted
