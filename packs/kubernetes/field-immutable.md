---
title: Field is immutable
aliases: [fieldimmutable, immutable-field]
category: apiserver
---

# Field is immutable

## What it means
Some fields on Kubernetes resources can't be changed after creation by
design, the API server rejects any update touching them. This is most often
hit trying to `kubectl apply` a manifest that changes something structural,
not from hand-editing a live object.

## Common causes
1. Changing a StatefulSet's `volumeClaimTemplates` after creation, storage
   binding for existing pods can't be altered in place.
2. Changing a Job's `spec.selector` or `spec.template` (Jobs are meant to be
   immutable once created, a changed Job should be a new Job).
3. Changing a Service's `spec.clusterIP` or switching certain `type` values
   in ways that aren't supported as an in-place transition.
4. Changing a PersistentVolume's `spec` fields that define the underlying
   storage backend.

## Check it
```
kubectl apply -f <manifest> --dry-run=server
```
The dry-run flag surfaces this exact error without actually attempting the
change, useful for catching it in CI before a real deploy fails partway
through. The error message names the specific immutable field.

## Fix it
- If the field genuinely needs to change, the resource typically needs to be
  deleted and recreated, understand what that means for anything depending
  on it first (a StatefulSet's PVCs aren't deleted automatically, for
  example, which can be useful or dangerous depending on intent).
- For Jobs specifically, don't try to update one in place, create a new Job
  instead, that's the intended pattern.
- Add `--dry-run=server` to your deploy pipeline's validation step so this
  gets caught during review rather than mid-rollout.

## Related
- ResourceVersionConflict
- HelmUpgradeFailed
