---
title: warning, hook failed
aliases: [warninghookfailed, preinstallhookfailed, prehookfailed]
category: release
---

# Hook failed

## What it means
Helm hooks (`pre-install`, `pre-upgrade`, `post-install`, and others) are
just Kubernetes Jobs (usually) that Helm runs at specific points in a
release's lifecycle. If a hook's Job fails, Helm treats the whole
install/upgrade as failed by default, even if every other resource in the
chart would have applied cleanly.

## Common causes
1. The hook Job itself fails for an ordinary reason, a bad image, a missing
   ConfigMap/Secret it depends on, a bug in the script it runs, this is a
   normal Job failure, just happening at a moment that blocks your release.
2. The hook depends on something not ready yet (a database migration hook
   running before the database itself is reachable).
3. The hook's `helm.sh/hook-delete-policy` isn't set the way you expect, so
   a failed hook Job's pod lingers and its logs get harder to find on a
   retry.
4. A timing assumption that held on a previous chart version no longer
   holds after some other change in ordering.

## Check it
```
kubectl get jobs -l "app.kubernetes.io/managed-by=Helm" -n <namespace>
kubectl logs job/<hook-job-name> -n <namespace>
```
Debug the hook exactly like you would any other failing Kubernetes Job (see
`BackoffLimitExceeded` in the kubernetes pack), the fact that Helm triggered
it doesn't change how you diagnose it.

## Fix it
- Fix the underlying cause in the hook's Job the same way you'd fix any
  other failing Job, bad image, missing dependency, script bug.
- If the hook depends on something not guaranteed ready yet, add retry
  logic inside the hook script itself rather than assuming ordering alone
  guarantees readiness.
- Set an explicit `helm.sh/hook-delete-policy` (e.g.
  `before-hook-creation,hook-succeeded`) so retries don't accumulate stale
  failed Jobs that make debugging harder over time.

## Related
- BackoffLimitExceeded
- HelmUpgradeFailed
