---
title: DeadlineExceeded (Job)
aliases: [jobdeadlineexceeded, activedeadlineexceeded]
category: job
---

# DeadlineExceeded (Job)

## What it means
A `Job` set `spec.activeDeadlineSeconds`, and the job's total running time
(across all retries combined) exceeded it, so Kubernetes killed it and marked
it Failed regardless of `backoffLimit`. This is a wall-clock ceiling, separate
from the retry-count ceiling.

## Common causes
1. The task is simply slower than the deadline allows, either it grew over
   time or the deadline was set optimistically.
2. Repeated retries (each one taking real time) add up past the deadline even
   though no single attempt was unusually slow.
3. The job is waiting on a slow or stuck dependency for most of its runtime
   rather than doing real work, and that wait time counts against the
   deadline too.
4. `activeDeadlineSeconds` was copy-pasted from a different, faster job
   template and never adjusted for this job's actual workload.

## Check it
```
kubectl describe job <name>
kubectl get job <name> -o jsonpath='{.spec.activeDeadlineSeconds}'
```
Compare the configured deadline against how long the job's pods actually ran
before being killed, `describe job` shows the Failed condition with a
`DeadlineExceeded` reason and timing.

## Fix it
- Raise `activeDeadlineSeconds` to reflect the job's real runtime, with
  headroom, not the bare minimum observed once.
- If retries are what's pushing it over, investigate why attempts are
  failing (see BackoffLimitExceeded) rather than only extending the deadline.
- If the job is mostly waiting on a dependency, fix that dependency's
  availability rather than giving the job more time to keep waiting.

## Related
- BackoffLimitExceeded
- CronJobStuck
