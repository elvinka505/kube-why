---
title: Pending
aliases: [pending, unschedulable, stuck-pending]
category: scheduling
---

# Pod stuck in Pending

## What it means
The pod has been accepted by the API server but the scheduler hasn't been able
to place it on any node. Nothing is running yet, this happens before a
container ever starts.

## Common causes
1. No node has enough free CPU/memory to satisfy the pod's `resources.requests`.
2. A `nodeSelector`, `affinity`, or `taint`/`toleration` mismatch rules out every
   available node.
3. Not enough `PersistentVolume`s available to satisfy a `PersistentVolumeClaim`
   the pod needs.
4. The cluster's autoscaler is supposed to add a node but hasn't finished yet,
   or has hit a quota/limit.
5. A `PodDisruptionBudget` or topology spread constraint makes placement
   mathematically impossible with the current node layout.

## Check it
```
kubectl describe pod <pod>
```
The Events section will say exactly why, usually a line like `0/5 nodes are
available: 3 Insufficient cpu, 2 node(s) had taints that the pod didn't
tolerate`. That single line tells you which cause you're dealing with, read it
before guessing.

## Fix it
- Insufficient resources: lower the pod's requests if they're inflated, or add
  capacity (scale the node pool, or let the autoscaler catch up).
- Taint/toleration or affinity mismatch: fix the pod spec to match available
  nodes, or fix the node labels/taints if they're wrong.
- PVC pending: check `kubectl get pvc` and `kubectl describe pvc`, the storage
  class may not support dynamic provisioning, or you're out of volumes in that
  zone.
- Autoscaler lag: check the autoscaler's own logs/events, it may be blocked by a
  cloud quota rather than actually scaling.

## Related
- FailedScheduling
- InsufficientResources
