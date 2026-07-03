---
title: PIDPressure
aliases: [pidpressure, pid-pressure, too-many-processes]
category: node
---

# PIDPressure

## What it means
The node has hit its limit on the number of process IDs the kernel can hand
out. Once `PIDPressure` is true, new processes, including new containers,
can't start on that node, this shows up as pods failing to start with
confusing, generic-looking errors rather than an obvious "out of PIDs"
message.

## Common causes
1. A specific container is fork-bombing or leaking processes (zombie
   processes not being reaped, common when a container's PID 1 isn't
   handling child reaping correctly).
2. A high density of pods per node, each using more processes/threads than
   expected, collectively exhausting the node's PID limit.
3. No per-pod PID limits configured, so a single misbehaving pod can consume
   the entire node's PID budget.
4. A default PID limit that's simply too low for the workload mix now
   running on that node type.

## Check it
```
kubectl describe node <node>
```
Node access, if you have it: check running process counts directly
(`ps aux | wc -l`, or look at `/proc/loadavg` and cgroup PID counters). The
kubelet's Conditions section confirms `PIDPressure: True` from the cluster
side without needing node access.

## Fix it
- Identify the specific pod leaking or spawning excessive processes, this is
  almost always one workload's problem, not a genuine cluster-wide need for
  more PIDs.
- Ensure containers use a proper init process (`tini`, or the container
  runtime's built-in init) to reap zombie processes correctly instead of
  leaving them to accumulate.
- Set per-pod PID limits (`kubelet`'s `--pod-max-pids` or Kubernetes'
  `PodPidsLimit` feature) so one bad pod can't exhaust the whole node.
- If legitimate workloads need more processes than the default allows, raise
  the node's PID limit rather than leaving it uncapped for everyone.

## Related
- NodeNotReady
- MemoryPressure
