---
title: Ingress returning 503 Service Unavailable
aliases: [ingress503, serviceunavailable, 503-service-unavailable]
category: ingress
---

# Ingress returning 503 Service Unavailable

## What it means
The ingress controller has no backend at all to send the request to. This is
one layer earlier than a 502, a 502 means it tried a backend and failed, a 503
usually means it found nothing to try in the first place.

## Common causes
1. The Service referenced by the Ingress has no ready endpoints (the pods
   backing it are down, unready, or the selector matches nothing).
2. The Ingress rule's `path`/`host` doesn't actually match any configured
   backend, so the controller falls through to its default 503 response.
3. The Service name or port in the Ingress spec has a typo and doesn't match
   an existing Service.
4. The ingress controller itself just reloaded its config (new Ingress
   applied) and hasn't finished propagating yet, usually resolves in seconds.

## Check it
```
kubectl describe ingress <name>
kubectl get endpoints <service-name>
kubectl logs -n <ingress-namespace> <ingress-controller-pod>
```
`describe ingress` shows exactly which Service/port each rule points to,
confirm that Service exists and has endpoints. If everything looks correct on
paper, check the controller's logs for a config reload happening right around
the same time as the errors.

## Fix it
- No ready endpoints: this is the same root problem as "Service has no
  endpoints," fix the backing pods.
- Path/host mismatch: correct the Ingress rule to match what clients are
  actually requesting.
- Typo in Service name/port: fix the reference in the Ingress spec.
- If it's transient during a reload, no action needed, if it's persistent
  after a config change, check the controller logs for a rejected/invalid
  Ingress spec that it silently failed to apply.

## Related
- ServiceUnreachable
- Ingress502BadGateway
