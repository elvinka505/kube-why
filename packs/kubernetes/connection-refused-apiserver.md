---
title: Unable to connect to the API server
aliases: [connectionrefusedapiserver, unabletoconnecttotheserver, kubectl-connection-refused]
category: apiserver
---

# Unable to connect to the API server

## What it means
`kubectl` (or any client) can't reach the Kubernetes API server at all. This
is happening before authentication or authorization even come into play, the
network connection itself isn't being established.

## Common causes
1. The cluster endpoint in your kubeconfig is wrong, stale, or points at a
   cluster that no longer exists (very common after a cluster was recreated
   with a new control plane endpoint).
2. You're not on the required network (VPN not connected, not on a private
   network the API server is only reachable from).
3. The control plane itself is down or unreachable, rare in managed cloud
   offerings, more common in self-managed clusters during an incident.
4. A local proxy, firewall, or expired VPN session is silently dropping the
   connection.

## Check it
```
kubectl config current-context
kubectl config view --minify
curl -k https://<api-server-endpoint>/healthz
```
Confirm which context and endpoint you're actually pointed at, it's easy to
be debugging the wrong cluster entirely after switching contexts earlier.
The raw `curl` to `/healthz` isolates whether this is a `kubectl`/kubeconfig
problem or a genuine network/control-plane problem.

## Fix it
- Wrong/stale endpoint: refresh your kubeconfig from the source of truth
  (cloud CLI's get-credentials command, or whoever issues cluster access).
- VPN/network access: connect to the required network before retrying.
- Control plane genuinely down: this is an infrastructure incident, not a
  client-side fix, escalate accordingly.
- Local proxy/firewall interference: test with `--insecure-skip-tls-verify`
  temporarily to isolate TLS issues from pure connectivity issues, don't
  leave it set that way.

## Related
- Unauthorized
