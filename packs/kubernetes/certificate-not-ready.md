---
title: Certificate stuck not ready (cert-manager)
aliases: [certificatenotready, cert-manager-stuck, acme-challenge-failed]
category: ingress
---

# Certificate stuck not ready (cert-manager)

## What it means
cert-manager created a `Certificate` resource but hasn't gotten a signed
certificate back yet, so the Secret it's supposed to populate stays empty or
stale, and TLS on your Ingress fails or falls back to a default/self-signed
cert.

## Common causes
1. The ACME HTTP-01 challenge can't be reached from outside, usually because
   the Ingress routing the challenge path isn't actually publicly accessible
   yet (DNS not propagated, firewall blocking).
2. The DNS-01 challenge's DNS provider credentials are wrong, or the DNS
   record it creates hasn't propagated within the validation window.
3. Rate limiting from the certificate authority (Let's Encrypt's staging vs.
   production rate limits are a very common trap during testing).
4. The `Issuer`/`ClusterIssuer` referenced by the Certificate doesn't exist,
   is misconfigured, or is itself not Ready.

## Check it
```
kubectl describe certificate <name> -n <namespace>
kubectl describe certificaterequest -n <namespace>
kubectl describe order -n <namespace>
kubectl describe challenge -n <namespace>
```
cert-manager's resources form a chain, Certificate to CertificateRequest to
Order to Challenge, walk down that chain with `describe`, the failure is
usually visible as an Event on the lowest-level resource that's stuck.

## Fix it
- HTTP-01 unreachable: confirm the challenge Ingress is actually publicly
  resolvable and not blocked before cert-manager tries to validate it.
- DNS-01 credential/propagation issues: verify the provider secret is correct
  and consider increasing the propagation wait time for slow DNS providers.
- Rate limited: switch to Let's Encrypt's staging issuer while testing, only
  point at production once you're confident it'll succeed.
- Broken Issuer: fix or recreate the `Issuer`/`ClusterIssuer`, a Certificate
  can't succeed if the thing it depends on isn't Ready either.

## Related
- Ingress502BadGateway
- WebhookCallFailed
