---
title: object is being deleted, cannot create resource
aliases: [objectisbeingdeletedcannotcreateresource, blockedbyexistingtermination]
category: install
---

# Object is being deleted, cannot create resource

## What it means
Helm is trying to create a resource with the same name as one that's
currently mid-deletion (stuck `Terminating`, often due to a finalizer, see
the Kubernetes pack's `Terminating stuck`). The API server won't let a new
object be created under a name that's still tied up in an in-progress
deletion.

## Common causes
1. A previous `helm uninstall` (or manual delete) is stuck, a resource has
   a finalizer that hasn't cleared, and a new install/upgrade immediately
   afterward collides with it.
2. Uninstalling and reinstalling a release in quick succession, without
   waiting for the previous uninstall's resources to actually finish
   terminating.
3. A namespace itself is stuck `Terminating` (see the Kubernetes pack), and
   Helm can't create anything inside it.

## Check it
```
kubectl get <kind> <name> -n <namespace> -o yaml | grep -A5 finalizers
```
Same diagnostic as the Kubernetes pack's stuck-terminating entries, this is
that same underlying problem, just encountered through a Helm operation
instead of a direct `kubectl delete`.

## Fix it
- Find and fix whatever's blocking the stuck resource's finalizer (see
  `Terminating stuck` and `Namespace stuck Terminating` in the kubernetes
  pack), the fix is identical, Helm being involved doesn't change the
  underlying cause.
- Wait for a previous uninstall to fully complete before immediately
  reinstalling, `helm uninstall --wait` blocks until resources are actually
  gone rather than returning as soon as deletion is requested.

## Related
- HelmReleaseNotFound
