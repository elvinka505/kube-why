---
title: Service has no endpoints
aliases: [serviceunreachable, no-endpoints, endpointsnotfound]
category: networking
---

# Service has no endpoints

## What it means
A `Service` routes traffic to pods matched by its label selector. If nothing
currently matches that selector, the Service has zero endpoints, and any
traffic sent to it fails, connection refused or timeout, even though the
Service object itself looks perfectly fine.

## Common causes
1. The Service's `selector` labels don't match the labels actually on the
   pods, often after a rename or a copy-pasted manifest with stale labels.
2. The matching pods exist but aren't `Ready`, endpoints only include pods
   that pass their readiness probe.
3. The Service targets the wrong `targetPort`, so even matched pods aren't
   actually listening where traffic is sent.
4. The pods matching the selector were scaled to zero or deleted, and nothing
   noticed because the Service itself doesn't error, it just goes quiet.

## Check it
```
kubectl get endpoints <service-name>
kubectl describe service <service-name>
kubectl get pods -l <same-selector-as-service> -o wide
```
If `get endpoints` shows `<none>`, that's confirmed, the Service has nothing
to send traffic to. Compare the Service's selector directly against the
labels on the pods you expect it to hit.

## Fix it
- Selector/label mismatch: fix whichever one is wrong so they actually match,
  usually the Service selector was written before a pod label got renamed.
- Pods not Ready: fix the underlying readiness probe issue (see
  ReadinessProbeFailed), the Service will populate endpoints automatically
  once pods pass their check.
- Wrong `targetPort`: match it to the port the container actually listens on,
  not the Service's external `port`.
- Add a `kubectl get endpoints` check to your deploy pipeline for critical
  services, so an empty-endpoints state gets caught before someone reports the
  outage.

## Related
- ReadinessProbeFailed
- DNSResolutionFailed
