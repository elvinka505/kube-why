---
title: ContainerCreating stuck
aliases: [containercreating, stuck-containercreating, container-creating]
category: pod
---

# Stuck in ContainerCreating

## What it means
The pod has been scheduled to a node, but the container hasn't actually started
yet. This state is normal for a few seconds, it's a problem when it lasts
minutes. The delay is almost always something the kubelet is waiting on before
it can start the container: an image pull, a volume mount, or a network setup
step.

## Common causes
1. A large image still pulling (check for this before assuming something is
   broken).
2. A volume can't be mounted: a PVC isn't bound yet, an NFS/CSI backend is slow
   or unreachable, or a Secret/ConfigMap volume references something missing.
3. CNI/network plugin issue on the node, the sandbox can't get an IP.
4. The node itself is unhealthy or overloaded and slow to do anything.

## Check it
```
kubectl describe pod <pod>
```
Look at the Events section in order, they're timestamped and show exactly what
the kubelet is waiting on: `Pulling image`, `FailedMount`, `FailedCreatePodSandBox`.
Whichever event is the most recent and repeating is your answer.
```
kubectl get events --field-selector involvedObject.name=<pod> --sort-by='.lastTimestamp'
```

## Fix it
- Slow image pull: check image size, consider a smaller base image or a
  registry closer to your nodes.
- FailedMount: check PVC status (`kubectl get pvc`) and confirm the underlying
  storage backend is healthy.
- CNI/network errors: check the CNI plugin's pods/logs in `kube-system`, this is
  often a node-level or cluster-level problem, not specific to your pod.
- If only one node shows this behavior, cordon it and let pods reschedule
  elsewhere while you investigate the node separately.

## Related
- FailedMount
- Pending
