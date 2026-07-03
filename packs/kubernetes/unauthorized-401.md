---
title: Unauthorized (401)
aliases: [unauthorized401, unauthorized, invalid-bearer-token]
category: apiserver
---

# Unauthorized (401)

## What it means
The API server rejected the request because it couldn't verify who you are,
this is a step earlier than Forbidden. Forbidden means "I know who you are
and you can't do that," Unauthorized means "I don't know who you are at all."

## Common causes
1. Your client certificate or token has expired, common with short-lived
   cloud provider tokens (EKS/GKE/AKS) that need periodic refresh.
2. The kubeconfig's `user` credentials don't match what the API server
   expects, often after a cluster was recreated and old credentials are
   stale.
3. A ServiceAccount token mounted into a pod expired or was revoked, and the
   pod hasn't restarted to pick up a fresh one (relevant for older, non-
   auto-rotating token setups).
4. System clock skew between client and server causing token
   validation to fail even though the token itself is technically valid.

## Check it
```
kubectl config view --minify
kubectl auth can-i get pods
```
If `auth can-i` itself fails with Unauthorized (not just returns "no"), that
confirms authentication is the problem, not permissions. Check your cloud
provider's specific token-refresh command if using short-lived tokens
(`aws eks get-token`, `gcloud container clusters get-credentials`, etc.).

## Fix it
- Refresh short-lived cloud provider tokens using the appropriate CLI
  command, most managed kubeconfig setups handle this automatically if
  configured correctly, verify the exec plugin in your kubeconfig is intact.
- Regenerate/update kubeconfig credentials from the source of truth if the
  cluster was recreated.
- For pods, restart them to pick up a freshly mounted ServiceAccount token if
  using an older non-rotating setup, modern clusters auto-rotate this by
  default.
- Check system clock sync (`ntp`/`chrony`) if this is happening consistently
  and token refresh doesn't fix it.

## Related
- Forbidden
- ConnectionRefusedAPIServer
