---
title: Forbidden
aliases: [forbidden, rbac-forbidden, 403]
category: rbac
---

# Forbidden

## What it means
The API server understood who you are (authentication passed) but refused the
action because your identity has no RBAC permission for it (authorization
failed). This applies equally to a human running `kubectl` and to a
ServiceAccount a pod is using to talk to the API.

## Common causes
1. No `Role`/`ClusterRole` grants the verb (`get`, `list`, `create`, `delete`,
   ...) on that resource to your user, group, or ServiceAccount.
2. The `RoleBinding`/`ClusterRoleBinding` exists but points at the wrong
   subject (wrong ServiceAccount name, wrong namespace, typo'd username).
3. A namespaced `Role` was used when a `ClusterRole` was actually needed (or
   vice versa), so the permission doesn't apply where you're operating.
4. The ServiceAccount is correct, but the pod spec never actually references
   it, so it's silently running as `default`, which has no extra permissions.

## Check it
```
kubectl auth can-i <verb> <resource> --as=<user-or-sa> -n <namespace>
kubectl describe rolebinding,clusterrolebinding -n <namespace>
```
`auth can-i` tells you definitively whether the permission exists before you
go hunting through bindings by hand. The exact error message also names the
verb, resource, and identity it rejected, read it closely, it's specific.

## Fix it
- Add a `RoleBinding` (namespace-scoped) or `ClusterRoleBinding`
  (cluster-wide) granting the missing verb/resource to the correct subject.
- If a pod's ServiceAccount is the problem, confirm `spec.serviceAccountName`
  in the pod/deployment actually references the intended ServiceAccount, not
  the namespace default.
- Prefer the narrowest `Role` that grants exactly what's needed over a broad
  `ClusterRole`, the debugging cost of finding a missing permission later is
  much lower than the blast radius of an overly broad one now.

## Related
- Unauthorized
- AdmissionWebhookDenied
