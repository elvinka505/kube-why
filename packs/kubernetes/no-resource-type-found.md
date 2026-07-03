---
title: The server doesn't have a resource type
aliases: [noresourcetypefound, doesnt-have-a-resource-type, unknown-resource-type]
category: apiserver
---

# The server doesn't have a resource type

## What it means
You ran `kubectl get`/`apply` against a resource kind the API server doesn't
recognize. Either the CRD that defines it was never installed, or you're
pointed at the wrong cluster, or the API version you're using has been
removed.

## Common causes
1. The CRD (Custom Resource Definition) that registers this type was never
   installed, or was installed in a different cluster than the one your
   kubeconfig currently points at.
2. A typo in the resource kind or its shortname.
3. The manifest uses an `apiVersion` that was removed in the current cluster
   version, Kubernetes does deprecate and remove old API versions on a
   schedule, this bites people during cluster upgrades specifically.
4. You're targeting the wrong context entirely, the resource exists, just
   not in the cluster you're currently connected to.

## Check it
```
kubectl config current-context
kubectl api-resources | grep -i <resource-name>
kubectl get crd | grep -i <resource-name>
```
Confirm your context first, this single mistake accounts for a lot of
otherwise-confusing "resource not found" reports. Then check whether the CRD
is actually installed at all.

## Fix it
- Wrong cluster/context: switch to the correct one.
- Missing CRD: install whatever operator/chart is supposed to provide it
  before applying resources of that kind.
- Removed API version: update the manifest to the current `apiVersion` for
  this resource kind, check the Kubernetes deprecation guide for your
  cluster's version for the replacement.
- Typo: fix the resource kind/shortname in the command or manifest.

## Related
- ConnectionRefusedAPIServer
