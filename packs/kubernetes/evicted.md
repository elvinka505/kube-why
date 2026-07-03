---
title: Evicted
aliases: [evicted, node-pressure]
category: pod
---

# Evicted

## What it means
The kubelet on the node proactively killed the pod to reclaim resources before
the node ran out entirely. This is different from OOMKilled: OOMKilled is the
kernel killing one container that individually exceeded its limit, Evicted is
the kubelet killing whole pods because the node itself is under pressure
(memory, disk, or PIDs).

## Common causes
1. Node is low on memory overall (not just this pod's limit, the node's total).
2. Node is low on disk space, often from log files, unpruned images, or a
   runaway container writing to disk.
3. Too many processes/PIDs on the node hitting the kernel's PID limit.
4. The pod has no resource requests set, so it's among the first candidates the
   kubelet picks when it needs to free something up.

## Check it
```
kubectl get pod <pod> -o yaml | grep -A5 status
kubectl describe node <node>
```
The pod status will show a `Reason` field, e.g. `Evicted`, and a `Message`
naming the pressure type: `The node was low on resource: memory` or
`ephemeral-storage`. Check the node's own conditions with `describe node` to
confirm which resource is actually under pressure right now.

## Fix it
- Disk pressure: clean up unused images (`crictl rmi` or let garbage collection
  run), find what's filling the disk (often application logs with no rotation).
- Memory pressure at the node level: this is a capacity problem, not a single
  pod's config, add nodes or reduce total scheduled load.
- Set resource requests on every pod so the scheduler can actually reason about
  capacity instead of overcommitting the node.
- If one specific pod is the actual culprit (log spam, disk-filling job), fix
  that pod rather than just adding capacity to mask it.

## Related
- OOMKilled
- Pending
