---
title: Resource has been modified (409 conflict)
aliases: [resourceversionconflict, object-has-been-modified, 409-conflict]
category: apiserver
---

# Resource has been modified (409 conflict)

## What it means
Kubernetes uses optimistic concurrency control: every object has a
`resourceVersion`, and an update must include the version you last read. If
someone (or something) else modified the object between your read and your
write, your update is rejected rather than silently overwriting their change.
This is the system working correctly, not a bug.

## Common causes
1. Two people (or a person and a controller/operator) are editing the same
   resource at nearly the same time.
2. A `kubectl edit` session sat open too long while something else (a
   controller reconciling the object, an autoscaler) updated the resource in
   the meantime.
3. A CI/CD pipeline applies manifests based on a stale cached copy of the
   resource rather than fetching the current version first.
4. Client-side retry logic reuses an old cached object across retries instead
   of re-fetching before each attempt.

## Check it
```
kubectl get <resource> <name> -o yaml
```
There's not much to diagnose beyond confirming this is genuinely a
concurrent-edit situation, check who or what else might be touching this
resource (a controller's reconcile loop is a common invisible actor here).

## Fix it
- Simplest fix: re-fetch the resource and reapply your change on top of the
  current version, this is expected, not exceptional, behavior.
- For `kubectl edit`, just retry the edit, it'll show you the current state
  again.
- In automation, always fetch-then-modify-then-update in the same operation
  rather than applying a manifest generated earlier from a stale read,
  or use server-side apply, which handles this more gracefully for
  declarative workflows.
- If this happens constantly on the same resource, something is reconciling
  it very frequently, worth understanding what and why before just retrying
  in a loop.

## Related
- FieldImmutable
