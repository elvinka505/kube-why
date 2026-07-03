---
title: CronJob not running or stuck
aliases: [cronjobstuck, missedschedule, cronjob-not-running]
category: job
---

# CronJob not running (or stuck)

## What it means
A `CronJob` is supposed to create a new `Job` on schedule, but either it
isn't creating one at all, or the previous Job is still running and blocking
the next one, depending on `concurrencyPolicy`.

## Common causes
1. `concurrencyPolicy: Forbid` (the default-ish safe choice) skips a new run
   entirely if the previous Job hasn't finished, and that previous Job is
   hung rather than genuinely still working.
2. The CronJob controller missed its scheduling window (`startingDeadlineSeconds`
   too tight, or the controller itself was down briefly) and, depending on
   config, just skips that run instead of running it late.
3. The schedule expression itself is wrong, most often a timezone
   misunderstanding, cron schedules run in UTC by default unless you've set
   `spec.timeZone`.
4. `successfulJobsHistoryLimit`/`failedJobsHistoryLimit` are set to 0, which
   makes it look like nothing ran because the evidence gets cleaned up
   immediately.

## Check it
```
kubectl get cronjob <name>
kubectl get jobs --selector=job-name --sort-by=.metadata.creationTimestamp
kubectl describe cronjob <name>
```
Check the CronJob's `LAST SCHEDULE` time against what you expected. Then look
at the Jobs it actually created (or didn't), a hung previous Job blocking a
`Forbid`-policy CronJob is the most common real-world cause.

## Fix it
- Hung previous Job blocking new runs: fix or manually clean up the stuck
  Job, and add `activeDeadlineSeconds` to the Job template so this can't
  happen silently again.
- Timezone confusion: set `spec.timeZone` explicitly instead of assuming
  local time.
- Missed windows: widen `startingDeadlineSeconds` if runs close together in
  time are acceptable, or investigate why the controller itself was delayed.
- If you need to debug history, temporarily raise the history limits so
  finished Jobs stick around long enough to inspect.

## Related
- BackoffLimitExceeded
- JobDeadlineExceeded
