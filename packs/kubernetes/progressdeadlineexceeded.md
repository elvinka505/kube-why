---
title: ProgressDeadlineExceeded
aliases: [progressdeadlineexceeded, deployment-stuck, rollout-stuck]
category: deployment
---

# ProgressDeadlineExceeded

## What it means
A Deployment's rollout hasn't made progress within `progressDeadlineSeconds`
(default 600). This is the Deployment controller giving up waiting for new pods
to become ready, it doesn't roll back automatically, the broken and working
pods can sit side by side until you intervene.

## Common causes
1. New pods are crashing (CrashLoopBackOff) so they never pass their readiness
   check and the rollout can't proceed.
2. New pods are stuck Pending, usually a resource or scheduling constraint that
   only the new pod spec triggers (a new resource request, a new node
   selector).
3. A readiness probe is misconfigured, pointing at the wrong port or path, so
   healthy pods never report ready.
4. `maxUnavailable`/`maxSurge` settings combined with limited cluster capacity
   mean the rollout can't create enough new pods to proceed.

## Check it
```
kubectl rollout status deployment/<name>
kubectl describe deployment <name>
kubectl get pods -l <selector> --show-labels
```
Find the new ReplicaSet's pods specifically and describe one of them, the
underlying problem is almost always visible in the individual pod's events, the
Deployment-level error just tells you the rollout stalled, not why.

## Fix it
- Diagnose the underlying pod issue first (it's usually one of CrashLoopBackOff,
  Pending, or a bad readiness probe), fix that, and the rollout resumes on its
  own.
- If the new version is genuinely broken, roll back rather than debugging in
  production:
  ```
  kubectl rollout undo deployment/<name>
  ```
- Once fixed, consider tightening `progressDeadlineSeconds` for faster failure
  detection on critical services, the default 10 minutes is a long time to
  serve degraded traffic before anyone's paged.

## Related
- CrashLoopBackOff
- Pending
- ReadinessProbeFailed
