---
title: Admission webhook denied the request
aliases: [admissionwebhookdenied, webhook-denied]
category: webhook
---

# Admission webhook denied the request

## What it means
A `ValidatingWebhookConfiguration` or `MutatingWebhookConfiguration` is
intercepting requests to the API server (policy engines like OPA/Gatekeeper,
Kyverno, or a custom admission controller), and it explicitly rejected this
one. The rejection is the webhook doing its job, the message just needs
decoding.

## Common causes
1. The resource genuinely violates a policy (missing required labels, a
   disallowed image registry, a security context that isn't permitted), and
   the webhook is correctly blocking it.
2. The policy itself is misconfigured or too broad, rejecting things it
   shouldn't (a common outcome right after a new policy rollout).
3. The webhook's own logic has a bug and is denying valid requests
   incorrectly.
4. The request is fine, but it's being evaluated against the wrong policy
   scope (a policy meant for one namespace accidentally applies cluster-wide).

## Check it
```
kubectl get validatingwebhookconfigurations,mutatingwebhookconfigurations
```
The denial message returned by `kubectl apply` almost always includes which
webhook rejected it and why, read that message first, it's usually specific
enough to point straight at the violated policy.

## Fix it
- If the policy is correct and your manifest violates it, fix the manifest,
  add the required label, use an approved image registry, adjust the
  security context.
- If the policy is wrong or too broad, that's a policy change, coordinate
  with whoever owns cluster policy rather than working around it per-request.
- If you suspect a bug in the webhook itself, check its logs directly, most
  policy engines log the full evaluation, not just the final verdict.
- For urgent unblocks, check whether the policy engine supports a scoped
  exception/exemption mechanism rather than disabling the webhook entirely.

## Related
- WebhookCallFailed
- Forbidden
