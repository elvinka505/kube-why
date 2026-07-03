---
title: NodeNotReady
aliases: [nodenotready, node-not-ready]
category: node
---

# NodeNotReady

## What it means
The kubelet on this node has stopped reporting a healthy status to the API
server, or hasn't reported at all within the expected window. The node
controller marks it `NotReady`, and after a further timeout, starts evicting
its pods so they can be rescheduled elsewhere.

## Common causes
1. The kubelet process crashed or was stopped on the node.
2. Network partition between the node and the control plane, the node is fine,
   it just can't be heard from.
3. The node ran out of disk space or memory badly enough that the kubelet
   itself can't function.
4. The underlying VM/instance was terminated, rebooted, or is unhealthy at the
   infrastructure level.

## Check it
```
kubectl describe node <node>
```
Check the Conditions section, `Ready: False` will have a `Reason` and `Message`
attached, e.g. `KubeletNotPosting` or a specific pressure condition. If you have
infrastructure access, check whether the underlying instance is even running:
```
kubectl get events --field-selector involvedObject.name=<node>
```

## Fix it
- If it's a resource exhaustion issue (disk/memory pressure crashed the
  kubelet), that's the immediate cause, not the root one, find out what filled
  the disk or memory and fix that, then the node should recover on its own once
  space frees up.
- If the instance itself is gone or broken at the infra layer, let the cluster
  autoscaler or your infra tooling replace it rather than trying to resuscitate
  a dead VM.
- If it's a transient network blip, the node often self-recovers, don't
  immediately drain/delete it, confirm it's actually gone first, premature
  deletion while it's just slow to report can cause pod duplication.
- For recurring flakiness on the same node, check kubelet logs on that node
  directly (`journalctl -u kubelet`), the pattern usually points to a specific
  resource or driver issue.

## Related
- Terminating stuck
- Evicted
