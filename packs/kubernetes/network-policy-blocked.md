---
title: Traffic blocked by NetworkPolicy
aliases: [networkpolicyblocked, connection-refused-networkpolicy]
category: networking
---

# Traffic blocked by NetworkPolicy

## What it means
A `NetworkPolicy` in the namespace is restricting traffic, and the connection
you expect to work is being silently dropped or refused. Unlike most errors
here, there's no event or log line that says "NetworkPolicy blocked this,"
the connection just fails as if nothing on the other end were listening.

## Common causes
1. A default-deny policy was added (common security baseline) and the
   specific `Ingress`/`Egress` rule allowing your traffic was never added
   alongside it.
2. The policy's pod selector doesn't match the labels on the pods you expect
   it to apply to, so it's not doing what whoever wrote it intended.
3. DNS egress (port 53 to `kube-system`) wasn't allowed, which breaks
   everything that depends on service discovery, not just the traffic you're
   actually debugging.
4. The policy allows traffic by namespace label, but the target namespace is
   missing that label, this is an easy one to miss since namespaces aren't
   labeled by default.

## Check it
```
kubectl get networkpolicy -n <namespace> -o yaml
kubectl describe networkpolicy <name> -n <namespace>
```
There's no direct "this policy blocked this connection" log in vanilla
Kubernetes, you have to read the policy and reason about it manually, or check
your CNI plugin's own logs/metrics if it exposes policy decisions (Cilium and
Calico both have tooling for this).

## Fix it
- Add an explicit `Ingress`/`Egress` rule permitting the traffic you need,
  don't just delete the default-deny policy, that removes the protection for
  everything else too.
- Fix the pod selector so the policy actually targets the pods it's meant to.
- Always pair a default-deny policy with an explicit DNS-egress allow rule, or
  every workload in the namespace breaks in a confusing way.
- If your CNI supports it, use its policy visualization/audit tooling instead
  of reasoning through YAML by hand, it's much faster for anything beyond a
  couple of rules.

## Related
- DNSResolutionFailed
- ServiceUnreachable
