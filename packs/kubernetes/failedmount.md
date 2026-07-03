---
title: FailedMount
aliases: [failedmount, mountvolume-setup-failed, failed-mount]
category: pod
---

# FailedMount / MountVolume.SetUp failed

## What it means
The kubelet can't attach or mount a volume the pod needs, so the container
can't start. This shows up as the pod hanging in ContainerCreating with a
FailedMount event, rather than as a distinct pod status.

## Common causes
1. The referenced PVC doesn't exist, isn't bound, or is bound to a PV in a
   different availability zone than the node the pod was scheduled to.
2. A Secret or ConfigMap volume references a resource that doesn't exist, same
   root issue as CreateContainerConfigError but for volumes instead of env vars.
3. The CSI driver for the storage backend is down, misconfigured, or missing
   required permissions.
4. The volume is already attached to another node (common with
   `ReadWriteOnce` volumes when a pod reschedules faster than the old volume
   detaches).

## Check it
```
kubectl describe pod <pod>
kubectl get pvc <pvc-name>
```
The Events section names the exact volume and the specific failure, `timeout
expired waiting for volumes to attach`, `volume is already exclusively attached
to one node`, or a missing Secret/ConfigMap name. Match the message to the
cause list above.

## Fix it
- Zone mismatch: check that your StorageClass uses `volumeBindingMode:
  WaitForFirstConsumer` so the volume is provisioned in the same zone the pod
  actually gets scheduled to, rather than being bound too early.
- Missing Secret/ConfigMap: same fix as CreateContainerConfigError, create it or
  fix the reference.
- Stuck ReadWriteOnce attachment: confirm the old pod is fully terminated
  before expecting the new one to mount, this usually resolves itself within a
  minute or two, force-deleting the old pod can speed it up if the node it was
  on is actually gone.
- CSI driver issues: check the CSI driver's own pods/logs in `kube-system`, this
  is a cluster-level problem, not something fixable from the pod spec.

## Related
- ContainerCreating stuck
- CreateContainerConfigError
