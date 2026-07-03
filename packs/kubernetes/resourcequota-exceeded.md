---
title: ResourceQuota exceeded
aliases: [resourcequotaexceeded, quota-exceeded, forbidden-quota]
category: quota
---

# ResourceQuota exceeded

## What it means
The namespace has a `ResourceQuota` capping total resource usage (CPU,
memory, object counts), and your request would push usage past that cap. The
API server rejects the request outright, it never gets scheduled, it's
refused at creation time.

## Common causes
1. The namespace's total pods (across everyone deploying to it) have
   collectively hit the CPU/memory quota ceiling, and your new pod's requests
   would exceed it, even if your specific request looks reasonable in
   isolation.
2. An object-count quota (max Pods, max PVCs, max Services) has been reached,
   unrelated to compute resources entirely.
3. Every pod in the namespace is required to specify `resources.requests`
   when a quota is active, and yours doesn't, so it's rejected on that
   technicality rather than on the numbers themselves.
4. A quota was tightened after workloads were already running near the old
   ceiling, and a routine redeploy now tips it over.

## Check it
```
kubectl describe resourcequota -n <namespace>
kubectl describe quota <name> -n <namespace>
```
This shows used vs. hard limit for every tracked resource side by side, the
one at or near its hard limit is your answer, and the rejection message
usually names it directly too.

## Fix it
- If usage is legitimately near the limit, request a quota increase from
  whoever owns namespace capacity planning, don't just keep retrying.
- If your pod is missing `resources.requests`, add it, this is often the
  actual blocker even when there's plenty of room left in the quota.
- Check for abandoned/unused resources counting against the quota (old
  PVCs, completed Jobs that were never cleaned up) before assuming you need
  more quota rather than less waste.

## Related
- LimitRangeRejected
- Pending
