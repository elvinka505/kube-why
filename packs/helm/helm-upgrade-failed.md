---
title: Helm upgrade failed
aliases: [helmupgradefailed, upgrade-failed-has-no-deployed-releases]
category: helm
---

# Helm upgrade failed

## What it means
Helm attempted to reconcile your chart's rendered manifests against the
cluster and something in that process rejected it, this is a broader bucket
than the stuck-lock case, the release isn't locked, an actual upgrade attempt
ran and failed.

## Common causes
1. A rendered manifest is invalid or violates an admission webhook/policy,
   the underlying `kubectl apply` equivalent failed, and Helm surfaces that
   as an upgrade failure.
2. A resource being modified has an immutable field (see FieldImmutable),
   Helm can't apply an in-place change to something Kubernetes doesn't allow
   changing.
3. "has no deployed releases" specifically means Helm has no known-good prior
   revision to compare against or roll back to, common on a first install
   that failed, there's nothing to fall back on yet.
4. A hook (pre-upgrade job, etc.) failed, and by default Helm treats a failed
   hook as blocking the whole upgrade.

## Check it
```
helm status <release> -n <namespace>
helm get manifest <release> -n <namespace>
kubectl describe <resource-type> <name> -n <namespace>
```
The Helm CLI's own error output usually names the specific resource and
underlying Kubernetes error, that's the fastest path in, `describe` on that
specific resource fills in the rest.

## Fix it
- Immutable field errors: this usually means the field genuinely requires
  deleting and recreating the resource, understand the consequences (for a
  StatefulSet's volumeClaimTemplates, for example) before doing that.
- Failed hooks: check the hook Job's own pod logs directly, it failed for the
  same reasons any Job fails, fix that underlying cause.
- "No deployed releases": if this is truly a fresh install that failed
  partway, clean up any partially-created resources and reinstall rather than
  trying to upgrade into a release with no baseline.
- Policy/webhook rejections: same fix as AdmissionWebhookDenied, resolve the
  actual policy violation in the chart's templates.

## Related
- HelmPendingUpgrade
- FieldImmutable
- AdmissionWebhookDenied
