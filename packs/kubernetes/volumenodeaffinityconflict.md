---
title: VolumeNodeAffinityConflict
aliases: [volumenodeaffinityconflict, node-affinity-conflict]
category: storage
---

# VolumeNodeAffinityConflict

## What it means
The PersistentVolume the pod's PVC is bound to is only accessible from nodes
in a specific zone (its `nodeAffinity`), but the pod got scheduled to a node
outside that zone, or can't be scheduled at all because no available node
satisfies both the pod's own constraints and the volume's zone requirement.

## Common causes
1. The StorageClass doesn't use `volumeBindingMode: WaitForFirstConsumer`, so
   the volume was provisioned in a zone before the scheduler had a chance to
   pick a node, sometimes landing in a zone with no suitable node.
2. The pod has its own `nodeSelector`/affinity rules that conflict with the
   zone the volume was already provisioned in.
3. A node in the correct zone was removed or scaled down after the volume was
   created, leaving the volume stranded with no eligible node left.
4. A statically pre-created PV has a hardcoded zone that doesn't match where
   your nodes actually run.

## Check it
```
kubectl describe pod <pod>
kubectl get pv <pv-name> -o yaml | grep -A10 nodeAffinity
kubectl get nodes --show-labels | grep topology.kubernetes.io/zone
```
Compare the PV's required zone against the zones your actual nodes are in,
the mismatch is usually immediately obvious once both are side by side.

## Fix it
- Set `volumeBindingMode: WaitForFirstConsumer` on the StorageClass so future
  volumes are provisioned in the same zone as the node the pod actually gets
  scheduled to, this is the real long-term fix.
- For an already-stranded volume, either add capacity in the required zone or
  recreate the PVC so a new volume gets provisioned correctly.
- If a pod's own affinity rules are the conflict, align them with where
  storage is actually available, or vice versa.

## Related
- PVCPending
- FailedScheduling
