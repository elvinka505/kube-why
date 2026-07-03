---
title: PersistentVolumeClaim stuck Pending
aliases: [pvcpending, pvc-stuck-pending]
category: storage
---

# PersistentVolumeClaim stuck Pending

## What it means
The PVC is waiting for a matching `PersistentVolume` to bind to, either one
provisioned dynamically by a `StorageClass`, or a pre-created static PV. Until
it binds, any pod referencing this PVC stays stuck in `ContainerCreating` or
`Pending`.

## Common causes
1. The `StorageClass` referenced doesn't exist, or has a typo, dynamic
   provisioning can't start without a valid one.
2. The provisioner backing the StorageClass is down or misconfigured (CSI
   driver pods crashing, missing cloud credentials).
3. The requested size or access mode doesn't match any available static PV,
   if you're not using dynamic provisioning.
4. The cluster (or your cloud account) is out of quota for the underlying
   storage type/zone.

## Check it
```
kubectl describe pvc <name>
kubectl get storageclass
kubectl get pv
```
`describe pvc` Events usually name the exact issue, `waiting for a volume to
be created`, `no persistent volumes available`, or a provisioner-specific
error. Confirm the StorageClass it references actually exists and check the
provisioner's own pods if the class looks fine.

## Fix it
- Missing/misspelled StorageClass: fix the PVC or Deployment spec to
  reference one that actually exists (`kubectl get storageclass`).
- Provisioner down: treat the CSI driver pods like any other crashing
  workload, check their logs directly, credential issues are common here.
- No matching static PV: create one with matching size/access mode, or switch
  to dynamic provisioning if that's available in this cluster.
- Quota exhausted: this is an infrastructure-level fix, request more quota or
  free up unused volumes.

## Related
- FailedMount
- VolumeNodeAffinityConflict
- StorageClassNotFound
