---
title: Ingress returning 502 Bad Gateway
aliases: [ingress502, badgateway, 502-bad-gateway]
category: ingress
---

# Ingress returning 502 Bad Gateway

## What it means
The ingress controller (nginx, traefik, etc.) accepted the request but
couldn't get a valid response from the backend Service it forwarded to. This
is a controller-to-pod problem, the request never made it back successfully,
distinct from a 503 which usually means the controller couldn't even find a
backend to try.

## Common causes
1. The backend pod crashed or is restarting exactly when the request landed.
2. The Ingress points at the wrong `targetPort`, so it's connecting to a port
   nothing is listening on.
3. The backend is slow enough to exceed the controller's proxy timeout, which
   the controller reports as a gateway failure, not a timeout, from the
   client's point of view.
4. A protocol mismatch, the backend expects HTTPS but the controller is
   connecting over plain HTTP, or vice versa.

## Check it
```
kubectl logs -n <ingress-namespace> <ingress-controller-pod>
kubectl get endpoints <service-name>
```
The ingress controller's own logs almost always show the specific upstream
error (connection refused, timeout, reset). Cross-check that the Service has
healthy endpoints at all, an empty endpoints list produces this same symptom.

## Fix it
- Confirm the Service has ready endpoints (see Service has no endpoints), a
  502 with zero endpoints is really that problem wearing a different name.
- Fix the `targetPort` if it doesn't match what the container listens on.
- If it's a timeout under load, raise the controller's proxy timeout
  annotations, but also look at why the backend is that slow, a longer
  timeout is a mitigation, not a fix.
- Match the protocol/scheme annotation on the Ingress to what the backend
  actually speaks.

## Related
- ServiceUnreachable
- Ingress503ServiceUnavailable
