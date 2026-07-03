---
title: Namespace stuck Terminating
aliases: [namespacestuckterminating, namespace-wont-delete]
category: namespace
---

# Namespace stuck Terminating

## What it means
You deleted a Namespace, and every resource inside it should be garbage
collected before the Namespace object itself disappears. If one of those
resources has a `finalizer` that never clears, the whole Namespace hangs in
`Terminating` indefinitely, sometimes for days.

## Common causes
1. A custom resource (installed by an operator) has a finalizer, and the
   operator/controller responsible for removing it is gone, crashed, or was
   uninstalled before the resources it manages.
2. An API resource type that used to exist in the namespace (from a removed
   CRD or deprecated API version) can no longer be enumerated by the API
   server, so cleanup silently can't proceed.
3. A finalizer is waiting on an external cloud resource (a load balancer, a
   volume) to be deleted, and that deletion is failing or stuck.

## Check it
```
kubectl get namespace <name> -o yaml
kubectl api-resources --verbs=list --namespaced -o name | xargs -n1 -I{} kubectl get {} -n <name> --ignore-not-found
```
The first command shows the Namespace's own status and conditions. The second
is the real diagnostic step, it lists every resource type still present in
the namespace, the finalizer-blocked resource is usually near the bottom of
that list, something that looks orphaned or came from a removed operator.

## Fix it
- Find the specific stuck resource and its finalizer, then find and fix the
  controller that's supposed to remove that finalizer, that's the real fix.
- If the owning controller is genuinely gone for good (uninstalled operator,
  removed CRD) and the resource is confirmed safe to abandon, you can edit
  the resource and remove the finalizer manually, understand this skips
  whatever cleanup that finalizer was meant to guarantee, check for orphaned
  cloud resources afterward.
- Don't reach for finalizer removal as a first move, it's a last resort once
  you've confirmed the owning controller truly isn't coming back.

## Related
- Terminating stuck
- AdmissionWebhookDenied
