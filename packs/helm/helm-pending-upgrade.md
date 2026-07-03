---
title: Helm release stuck pending-upgrade
aliases: [helmpendingupgrade, another-operation-in-progress, helm-stuck]
category: helm
---

# Helm release stuck pending-upgrade

## What it means
Helm tracks release state, and if a previous `install`/`upgrade`/`rollback`
was interrupted (network drop, CI job killed, laptop closed mid-deploy) before
it could mark itself finished, Helm believes an operation is still in
progress. Every subsequent `helm upgrade` fails immediately with "another
operation is in progress" even though nothing is actually running.

## Common causes
1. A CI/CD pipeline job was cancelled or timed out mid-deploy, leaving the
   release marked `pending-upgrade` or `pending-install` in Helm's stored
   state (a Secret in the release's namespace).
2. Someone ran `helm upgrade` and killed the terminal (Ctrl+C, closed laptop)
   before it completed.
3. The underlying Kubernetes resources actually did apply successfully, but
   Helm's own bookkeeping never got the final confirmation, so the mismatch
   is purely in Helm's state, not the cluster's.

## Check it
```
helm status <release> -n <namespace>
helm history <release> -n <namespace>
kubectl get secrets -n <namespace> -l owner=helm,name=<release>
```
`helm history` shows the stuck status directly, usually `pending-upgrade` or
`pending-install` as the most recent revision with no completed follow-up.

## Fix it
- `helm rollback <release> <last-good-revision> -n <namespace>` resolves this
  cleanly in most cases, it forces Helm back to a known-good state.
- If rollback also fails, `helm history` to find the last `deployed` revision,
  then manually mark the stuck revision as failed or superseded is a more
  invasive fix, only do this once you understand what's actually running in
  the cluster versus what Helm believes is running.
- In CI, always run deploys with a reasonable timeout and treat interruption
  as a signal to check release state before retrying, rather than blindly
  re-running `helm upgrade` into a still-locked release.

## Related
- HelmUpgradeFailed
