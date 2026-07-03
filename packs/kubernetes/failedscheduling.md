---
title: FailedScheduling
aliases: [failedscheduling, insufficientresources, insufficient-cpu, insufficient-memory]
category: scheduling
---

# FailedScheduling

## What it means
This is the specific event the scheduler emits when it can't find a node for a
pod. It's closely related to a pod being stuck Pending, this is the reason
behind it, shown as an Event rather than a pod status.

## Common causes
1. Every node lacks enough free CPU or memory to satisfy the pod's requests
   (`Insufficient cpu` / `Insufficient memory` in the message).
2. Taints on all available nodes that the pod doesn't tolerate.
3. `nodeAffinity` or `nodeSelector` rules out every node in the cluster.
4. Topology spread constraints or pod anti-affinity rules make placement
   impossible given how pods are currently distributed.
5. A `PersistentVolumeClaim` the pod needs is stuck unbound.

## Check it
```
kubectl describe pod <pod>
```
The message is specific and literal, e.g. `0/8 nodes are available: 5 Insufficient
memory, 3 node(s) didn't match node affinity/selector`. It counts nodes against
each individual reason, read it as a checklist rather than a single cause.

## Fix it
- Insufficient resources: right-size the pod's requests if they're inflated
  relative to actual usage, or add cluster capacity.
- Taint/toleration mismatch: add the required toleration to the pod, or confirm
  the taint is actually still needed on those nodes.
- Affinity/selector: verify the labels you're selecting for actually exist on
  some node (`kubectl get nodes --show-labels`), a renamed or removed label is a
  common silent cause.
- Anti-affinity/topology spread: relax the constraint or add more nodes across
  the zones/hosts it's trying to spread across.

## Related
- Pending
- FailedMount
