---
title: Webhook call failed (unreachable or timed out)
aliases: [webhookcallfailed, failed-calling-webhook]
category: webhook
---

# Webhook call failed (unreachable or timed out)

## What it means
This is different from a webhook denying a request, here the API server
couldn't even reach the webhook to ask it. Depending on `failurePolicy`, this
either blocks every matching request cluster-wide (`Fail`, the safer but more
disruptive default) or silently skips the check (`Ignore`).

## Common causes
1. The webhook's backing Service has no ready endpoints, its pods are down,
   crashing, or not passing readiness, so there's genuinely nothing to answer
   the call.
2. The webhook's TLS certificate is expired or doesn't match the
   `caBundle` configured in the webhook registration, causing the API server
   to refuse the TLS handshake.
3. A `NetworkPolicy` blocks the API server from reaching the webhook's pods.
4. The webhook is slow enough to exceed `timeoutSeconds`, which the API
   server treats the same as unreachable.

## Check it
```
kubectl get validatingwebhookconfigurations,mutatingwebhookconfigurations -o yaml
kubectl get endpoints -n <webhook-namespace> <webhook-service>
kubectl logs -n <webhook-namespace> <webhook-pod>
```
Confirm the webhook's Service has endpoints first, this is the single most
common cause. If it does, check the webhook pod's own logs for TLS or timeout
errors on incoming requests.

## Fix it
- No ready endpoints: fix the webhook's own pods like any other crashing or
  unready workload, this is high severity, with `failurePolicy: Fail` it can
  block unrelated deployments cluster-wide.
- Expired/mismatched certificate: rotate the webhook's TLS cert and make sure
  the `caBundle` in the webhook registration is updated to match.
- NetworkPolicy blocking the control plane: explicitly allow ingress from the
  API server's egress range to the webhook's pods.
- Chronic slowness: raise `timeoutSeconds` as a stopgap, but investigate why
  the webhook itself is slow, it's a shared dependency for every matching
  request cluster-wide.

## Related
- AdmissionWebhookDenied
- CertificateNotReady
