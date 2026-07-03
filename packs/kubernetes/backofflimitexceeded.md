---
title: BackoffLimitExceeded
aliases: [backofflimitexceeded, job-failed-backoff]
category: job
---

# BackoffLimitExceeded

## What it means
A `Job`'s pod failed enough times to hit `spec.backoffLimit` (default 6), so
the Job controller marks it permanently Failed and stops retrying. Unlike a
Deployment, a Job doesn't keep trying forever, this is a deliberate ceiling.

## Common causes
1. The underlying task genuinely fails every time (bad input, bug in the
   job's code), so retries were never going to help.
2. The job depends on something not ready yet (database, another service)
   and fails faster than that dependency becomes available, burning through
   retries before conditions improve.
3. `backoffLimit` is set too low for a task that has occasional, expected
   transient failures.
4. The job's pod is being OOMKilled or hitting a resource limit every attempt,
   which looks like a logic failure but is actually a sizing problem.

## Check it
```
kubectl describe job <name>
kubectl get pods -l job-name=<name>
kubectl logs <pod> --previous
```
List the pods the Job created, there are usually several, one per failed
attempt. Check logs across a few of them, not just the last one, to see if the
failure is identical every time (points to a deterministic bug) or varies
(points to a flaky dependency or resource issue).

## Fix it
- Deterministic failure: fix the actual bug, retries won't save a task that
  always fails the same way.
- Dependency timing: add an init container or startup retry logic so the job
  doesn't fail immediately, or raise `backoffLimit` combined with
  `activeDeadlineSeconds` to allow more attempts within a bounded time.
- Resource-related failures: check for OOMKilled specifically, size the job's
  resource requests/limits to what it actually needs.
- For jobs expected to have occasional transient failures, raise
  `backoffLimit` to a number that reflects that reality, not the default.

## Related
- OOMKilled
- JobDeadlineExceeded
