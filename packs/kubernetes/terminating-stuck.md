---
title: Terminating stuck
aliases: [terminating, stuck-terminating, wont-delete]
category: pod
---

# Pod stuck in Terminating

## What it means
You deleted the pod, Kubernetes marked it for deletion, and it's still there.
The API server sets a `deletionTimestamp` and waits for the kubelet to confirm
the containers actually stopped before removing the object. If that
confirmation never comes, the pod hangs in this state indefinitely.

## Common causes
1. The container isn't handling `SIGTERM`, so Kubernetes waits out the full
   `terminationGracePeriodSeconds` (default 30s) and then force-kills it, which
   should still resolve it, if it doesn't, something else is wrong.
2. The node the pod is on is unreachable (crashed, network partition, or was
   removed from the cluster) so nothing can report the container as stopped.
3. A `finalizer` on the pod is waiting for some external controller to remove
   it, and that controller is stuck, crashed, or gone.
4. A volume can't be unmounted cleanly, blocking the kubelet from finishing
   cleanup.

## Check it
```
kubectl get pod <pod> -o yaml | grep -A3 finalizers
kubectl get node <node-the-pod-was-on>
```
If the node itself is `NotReady` or gone, that's your answer, the kubelet that
was supposed to confirm termination no longer exists. If there's a finalizer
listed, that's what's blocking deletion, and you need to find the controller
responsible for removing it.

## Fix it
- If the node is genuinely gone (deleted from the cloud provider, crashed
  permanently), force delete the pod:
  ```
  kubectl delete pod <pod> --grace-period=0 --force
  ```
  Only do this once you've confirmed the node is actually gone, not just slow,
  force-deleting while the node is alive can lead to two copies of the same
  container running.
- If a finalizer is stuck, find and fix the controller that owns it rather than
  stripping the finalizer manually, removing it by hand can leave external
  resources (cloud load balancers, volumes) orphaned.
- If the app doesn't handle SIGTERM, add a signal handler for graceful shutdown
  so future deletes don't rely on the grace period timing out.

## Related
- NodeNotReady
