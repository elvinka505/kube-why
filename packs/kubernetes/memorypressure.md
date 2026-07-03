---
title: MemoryPressure
aliases: [memorypressure, memory-pressure]
category: node
---

# MemoryPressure

## What it means
The kubelet has detected the node's available memory has dropped below its
eviction threshold, this is a node-wide condition, not a single container
hitting its own limit (that's OOMKilled). Once `MemoryPressure` is true, the
kubelet starts evicting pods to bring usage back down.

## Common causes
1. Pods on the node collectively use more memory than the node has, often
   because many pods have no memory `requests` set, so the scheduler had no
   real signal to avoid overcommitting the node.
2. One or a few pods are consuming far more memory than their `requests`
   suggested, memory `limits` are enforced per-container, but a node can
   still fill up if many pods each use up to a generous limit simultaneously.
3. System-level processes (kubelet, container runtime, monitoring agents)
   are using more memory than budgeted for, leaving less available than
   expected for workloads.
4. A node type with too little memory for the pod density being scheduled
   onto it.

## Check it
```
kubectl describe node <node>
kubectl top node <node>
kubectl top pods --all-namespaces --field-selector spec.nodeName=<node>
```
`describe node` shows the `MemoryPressure` condition directly. `top pods`
filtered to that node shows you which pods are actually consuming the most,
that's usually where to look first.

## Fix it
- Set `resources.requests.memory` on every pod so the scheduler can actually
  reason about node capacity instead of overcommitting.
- Identify and fix or resize the specific pods consuming disproportionate
  memory relative to their requests.
- Reserve memory for system processes (`kube-reserved`/`system-reserved`) if
  it's consistently tighter than expected across the fleet.
- If workload density genuinely exceeds what these nodes can handle, use a
  larger node type or spread workloads across more nodes.

## Related
- OOMKilled
- Evicted
- DiskPressure
