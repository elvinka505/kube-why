---
title: Ingress falling through to default backend 404
aliases: [defaultbackend404, ingress404, no-matching-ingress]
category: ingress
---

# Ingress falling through to default backend (404)

## What it means
The request reached the ingress controller, but no Ingress rule matched it,
so the controller served its default backend, usually a generic 404 page. The
controller is working correctly, your routing rule just doesn't cover this
request.

## Common causes
1. The `host` in the Ingress rule doesn't match the `Host` header the client
   actually sent (a trailing dot, wrong subdomain, or testing against an IP
   directly instead of the hostname).
2. The `path` doesn't match, especially with `pathType: Exact` when the
   request has extra trailing segments, or vice versa with `Prefix` when you
   expected exact matching.
3. The Ingress was applied in the wrong namespace, or uses an
   `ingressClassName` that doesn't match the controller actually watching for
   it (common in clusters running multiple ingress controllers).
4. DNS points at the right load balancer, but the request's `Host` header
   still doesn't match any rule because of a typo in either place.

## Check it
```
kubectl get ingress -A
kubectl describe ingress <name>
curl -v -H "Host: <expected-host>" http://<controller-ip>/<path>
```
Testing with `curl` and an explicit `Host` header directly against the
controller's IP isolates whether this is a DNS problem or an Ingress rule
problem, that distinction changes where you look next.

## Fix it
- Fix the `host`/`path` in the Ingress rule to match what's actually being
  requested, or fix the client/DNS if the rule was correct all along.
- Confirm `ingressClassName` matches the controller you intend to handle this
  traffic, especially in clusters with more than one controller installed.
- Double check `pathType` semantics (`Exact` vs `Prefix` vs `ImplementationSpecific`)
  match what you actually want matched.

## Related
- Ingress503ServiceUnavailable
- DNSResolutionFailed
