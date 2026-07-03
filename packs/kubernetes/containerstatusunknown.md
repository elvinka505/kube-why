---
title: ContainerStatusUnknown
aliases: [containerstatusunknown, unknown-container-state]
category: container-runtime
---

# ContainerStatusUnknown

## What it means
The kubelet can't determine the actual state of a container, it's not
reporting Running, Waiting, or Terminated with any confidence, usually
because the kubelet lost communication with the container runtime or the
node itself became unreachable partway through tracking that container.

## Common causes
1. The node the pod is on became `NotReady` or was lost entirely, the
   kubelet that would normally report accurate status isn't reachable.
2. The container runtime (containerd/CRI-O) itself crashed or restarted,
   losing track of containers it was previously managing.
3. A kubelet restart happened at an inconvenient moment, and it hasn't
   finished re-syncing container state with the runtime yet.
4. Heavy node resource pressure is causing the runtime to become
   unresponsive to status queries without fully crashing.

## Check it
```
kubectl get node <node>
kubectl describe pod <pod>
```
Check the node's status first, if it's `NotReady`, that's very likely your
root cause, this pod's unknown status is a symptom of the node problem, not
an independent issue.

## Fix it
- If the node is `NotReady` or gone, this resolves the same way, let the
  pod reschedule elsewhere once the node situation is confirmed (see
  NodeNotReady), don't chase the individual pod's status directly.
- If the container runtime crashed on an otherwise healthy node, check its
  own logs/status (`systemctl status containerd`, if you have node access),
  restarting the runtime is sometimes necessary if it's genuinely wedged.
- If this is transient during a kubelet restart, it usually self-resolves
  within a short window, avoid taking disruptive action immediately.

## Related
- NodeNotReady
- Terminating stuck
