---
title: LimitRange rejected the pod
aliases: [limitrangerejected, limitrange-violation]
category: quota
---

# LimitRange rejected the pod

## What it means
A `LimitRange` in the namespace sets min/max/default resource bounds per
container, and your pod's `resources` fall outside them, so the API server
rejects it at creation time. This is different from a ResourceQuota, a
LimitRange constrains individual containers, a Quota constrains the namespace
total.

## Common causes
1. The container requests more CPU/memory than the LimitRange's `max` allows.
2. The container requests less than the LimitRange's `min`, which is easy to
   violate accidentally if requests were copy-pasted from a smaller reference
   workload.
3. No `resources` were specified at all, and the LimitRange's default
   conflicts with something else in the same manifest (a `limits` without a
   matching `requests`, for example).
4. A `maxLimitRequestRatio` constraint is set, capping how much larger
   `limits` can be relative to `requests`, and the manifest exceeds that
   ratio.

## Check it
```
kubectl describe limitrange -n <namespace>
kubectl get pod <pod> -o yaml | grep -A6 resources
```
`describe limitrange` shows the exact min/max/default/ratio bounds
configured for the namespace, compare your container's actual
requests/limits against them directly.

## Fix it
- Adjust the container's `resources.requests`/`resources.limits` to fall
  within the namespace's configured bounds.
- If the bounds themselves seem wrong for legitimate workloads in this
  namespace, that's a conversation with whoever owns the LimitRange, not
  something to work around per-deployment.
- If you're relying on LimitRange defaults, verify what they actually are,
  don't assume, they vary per namespace and are easy to forget about until
  something like this happens.

## Related
- ResourceQuotaExceeded
- OOMKilled
