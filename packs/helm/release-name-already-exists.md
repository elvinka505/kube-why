---
title: release name already exists / cannot re-use a name that is still in use
aliases: [releasenamealreadyexists, cannotreuseanamethatisstillinuse]
category: install
---

# Release name already exists

## What it means
You tried to `helm install` using a release name that's already taken in
that namespace, including one that failed partway through a previous
install and was never cleaned up. Helm tracks release state by name, a
failed release still counts as "in use" until it's uninstalled or the name
is freed up.

## Common causes
1. A previous `helm install` with the same name failed partway through and
   was never uninstalled, the release object still exists in a `failed`
   state.
2. Re-running an install script or CI pipeline against the same environment
   without checking whether the release already exists.
3. Two different teams or pipelines unintentionally using the same release
   name in the same namespace.
4. Confusing `helm install` with `helm upgrade`, install always expects the
   name to be new, upgrade is what you want for an existing release.

## Check it
```
helm list -n <namespace> --all
helm status <release> -n <namespace>
```
`--all` matters, it includes failed and uninstalled-but-not-purged
releases that a plain `helm list` wouldn't show, which is exactly the state
that causes this error.

## Fix it
- If you meant to update an existing release, use `helm upgrade` (or
  `helm upgrade --install` to handle both cases in one command) instead of
  `helm install`.
- If a previous failed install left a stuck release behind, uninstall it
  first, then retry:
  ```
  helm uninstall <release> -n <namespace>
  ```
- In CI, prefer `helm upgrade --install` over a plain `install` for
  idempotent pipelines that might run against an environment more than
  once.

## Related
- HelmPendingUpgrade
- HelmReleaseNotFound
