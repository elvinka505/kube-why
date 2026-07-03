---
title: DNS resolution failing inside the cluster
aliases: [dnsresolutionfailed, corednsfailure, cannot-resolve-service]
category: networking
---

# DNS resolution failing inside the cluster

## What it means
A pod can't resolve another Service's DNS name (`myservice.mynamespace.svc.cluster.local`
or its short form). This is almost always a CoreDNS problem, not a problem
with the Service or pod you're actually trying to reach, DNS is a shared
dependency for the whole cluster.

## Common causes
1. CoreDNS pods are crashing, unready, or under-provisioned for the query
   volume (common right after a cluster scales up quickly).
2. The pod's `dnsPolicy` was overridden to something that skips cluster DNS
   (`Default` uses the node's resolver instead of CoreDNS, easy to set by
   accident).
3. A `NetworkPolicy` blocks egress to the CoreDNS pods on port 53.
4. You're using the short name (`myservice`) from a different namespace than
   the Service lives in, without the namespace suffix, which only resolves
   within the same namespace.

## Check it
```
kubectl get pods -n kube-system -l k8s-app=kube-dns
kubectl exec <pod> -- nslookup <service-name>.<namespace>.svc.cluster.local
kubectl logs -n kube-system -l k8s-app=kube-dns
```
Confirm CoreDNS pods are healthy first, that rules out or confirms the most
common cause immediately. Then test resolution with the fully-qualified name
to rule out a namespace-scoping mistake.

## Fix it
- CoreDNS unhealthy: treat it like any other crashing/OOMKilled workload,
  check its logs and resource limits, it's frequently under-resourced by
  default for larger clusters.
- Wrong `dnsPolicy`: set it back to `ClusterFirst` (the default) unless you
  have a specific reason not to.
- NetworkPolicy blocking DNS: explicitly allow egress to the `kube-system`
  namespace on UDP/TCP 53 in your policy.
- Namespace-scoping mistake: use the fully-qualified name
  (`service.namespace.svc.cluster.local`) when calling across namespaces.

## Related
- ServiceUnreachable
- NetworkPolicyBlocked
