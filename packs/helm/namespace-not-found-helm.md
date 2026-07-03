---
title: namespaces not found, create namespace
aliases: [namespacesnotfoundhelm, createnamespaceflagmissing]
category: install
---

# Namespaces "x" not found

## What it means
Unlike `kubectl apply`, `helm install` does not create the target namespace
automatically unless you explicitly tell it to. If the namespace doesn't
already exist, the install fails immediately trying to create the first
namespaced resource in it.

## Common causes
1. Installing into a namespace that hasn't been created yet, assuming Helm
   behaves like `kubectl apply -f` (which also doesn't auto-create
   namespaces, but this assumption is common regardless).
2. A CI/CD pipeline that provisions namespaces in a separate step that
   didn't run, was skipped, or ran against the wrong environment.
3. A typo in the target namespace name, so it's genuinely a different,
   nonexistent namespace than the one that was actually meant to be
   created earlier.

## Check it
```
kubectl get namespace <namespace>
```
Confirm directly whether the namespace exists at all before assuming
anything else is wrong, this is a one-command check that resolves the
ambiguity immediately.

## Fix it
- Add `--create-namespace` to the install command if you want Helm to
  create it as part of the install:
  ```
  helm install <release> <chart> -n <namespace> --create-namespace
  ```
- If your pipeline provisions namespaces as a separate deliberate step
  (common when namespaces carry labels, quotas, or RBAC that Helm
  shouldn't own), fix that step rather than relying on `--create-namespace`
  to paper over a broken pipeline order.

## Related
- ReleaseNotFound
