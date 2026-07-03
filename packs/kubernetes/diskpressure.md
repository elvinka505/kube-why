---
title: DiskPressure
aliases: [diskpressure, disk-pressure]
category: node
---

# DiskPressure

## What it means
The kubelet has detected the node's available disk space (or inodes) has
dropped below its eviction threshold. Once a node reports `DiskPressure`,
the kubelet starts evicting pods to reclaim space, and the scheduler avoids
placing new pods there.

## Common causes
1. Application logs writing to the node's local filesystem with no rotation,
   this is one of the most common real-world causes.
2. Too many unused container images accumulated on the node, image garbage
   collection lags behind image churn.
3. A container writing large amounts of data to an `emptyDir` volume backed
   by node disk rather than memory.
4. The node's disk is simply undersized for the workload density scheduled
   onto it.

## Check it
```
kubectl describe node <node>
kubectl get pods --all-namespaces --field-selector spec.nodeName=<node> -o wide
```
Check the node's Conditions for `DiskPressure: True`. If you have node
access, `df -h` and `du -sh` on likely directories (container logs, image
storage) pinpoint what's actually consuming the space.

## Fix it
- Fix log rotation for applications writing to local disk, or better, ship
  logs to a centralized backend instead of accumulating them locally.
- Manually trigger or tune image garbage collection thresholds if image
  buildup is the cause.
- Move large scratch-space usage to memory-backed `emptyDir` (with a size
  limit) if the workload can tolerate that, or to a proper volume instead of
  node-local disk.
- If it's a genuine capacity issue, resize the node's disk or the node type
  itself.

## Related
- Evicted
- NodeNotReady
