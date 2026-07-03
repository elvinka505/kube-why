---
title: release, not found
aliases: [releasenotfound, helmreleasenotfound]
category: release
---

# Release: not found

## What it means
You ran a command (`upgrade`, `rollback`, `status`, `uninstall`) against a
release name that Helm has no record of in the namespace you're targeting.
Helm releases are scoped to a specific namespace, a release that exists
elsewhere isn't visible unless you're pointed at the right one.

## Common causes
1. Targeting the wrong namespace, the release exists, just not where you're
   currently looking (`helm` defaults to the namespace in your current
   kubeconfig context if `-n` isn't specified).
2. A typo in the release name.
3. The release was already uninstalled (intentionally or by someone else)
   before this command ran.
4. Confusing a Kubernetes resource's name with the Helm release name, they
   aren't always the same string depending on how the chart was installed.

## Check it
```
helm list --all-namespaces
```
This searches across every namespace at once, immediately telling you
whether the release exists somewhere else, was already removed, or never
existed under that exact name at all.

## Fix it
- Point the command at the correct namespace with `-n <namespace>` once
  you've found where the release actually lives.
- Fix the release name if it was a typo.
- If it was already uninstalled and you meant to update it, that's now a
  fresh `helm install`, not an `upgrade`, since there's nothing existing to
  upgrade.

## Related
- ReleaseNameAlreadyExists
