---
title: StorageClass not found
aliases: [storageclassnotfound, no-storage-class]
category: storage
---

# StorageClass not found

## What it means
A PVC references a `storageClassName` that doesn't exist in the cluster, so
there's nothing to provision against and the claim can never bind. This is a
simpler, more direct version of PVCPending, the class itself is missing, not
just misconfigured.

## Common causes
1. A typo in `storageClassName` in the PVC or a Helm chart's values.
2. The manifest was copied from another cluster or environment where a
   differently-named StorageClass exists (`gp2` vs `gp3`, `standard` vs
   `standard-rwo`, this varies a lot between cloud providers and even cluster
   versions).
3. No default StorageClass is set, and the PVC omitted `storageClassName`
   expecting one to apply automatically.
4. The StorageClass was deleted (intentionally or not) while PVCs still
   reference it.

## Check it
```
kubectl get storageclass
kubectl get pvc <name> -o jsonpath='{.spec.storageClassName}'
```
List what actually exists, then compare directly against what the PVC is
asking for, this is almost always an exact-match problem, not a subtle
config issue.

## Fix it
- Fix the typo, or update the manifest to reference a StorageClass that
  actually exists in this cluster.
- If you expected a default StorageClass to apply, check whether one is
  actually marked default:
  ```
  kubectl get storageclass -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.metadata.annotations.storageclass\.kubernetes\.io/is-default-class}{"\n"}{end}'
  ```
- When porting manifests between clusters/clouds, don't assume StorageClass
  names match, verify them for the target environment explicitly.

## Related
- PVCPending
