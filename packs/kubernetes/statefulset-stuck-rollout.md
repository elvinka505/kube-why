---
title: StatefulSet rollout stuck
aliases: [statefulsetstuckrollout, statefulset-not-updating]
category: statefulset
---

# StatefulSet rollout stuck

## What it means
StatefulSets roll out one pod at a time, in order, and won't move to the next
pod until the current one is Running and Ready. If pod N never becomes ready,
every pod after it in the ordinal sequence just waits, indefinitely, with no
error saying "stuck," it looks like nothing is happening.

## Common causes
1. Pod N is failing for one of the usual pod-level reasons (CrashLoopBackOff,
   FailedMount, ReadinessProbeFailed), the StatefulSet's rollout is only ever
   as fast as its least healthy pod.
2. `podManagementPolicy` is the default `OrderedReady`, which is correct for
   most stateful workloads but means one bad pod genuinely blocks all the
   others, this isn't a bug, it's how ordering guarantees work.
3. `updateStrategy.rollingUpdate.partition` is set higher than expected,
   which intentionally holds back updates to lower-ordinal pods, easy to
   forget it's set after a canary test.
4. A `PersistentVolumeClaim` for a specific ordinal is stuck (see PVCPending),
   blocking that specific pod from ever starting.

## Check it
```
kubectl rollout status statefulset/<name>
kubectl get pods -l <selector> -o wide
kubectl describe pod <name>-<stuck-ordinal>
```
Identify the specific stuck ordinal first, then debug that one pod like any
other pod issue, the StatefulSet-level view won't tell you why, the pod-level
Events will.

## Fix it
- Fix the underlying pod issue on the blocking ordinal, once it's Ready, the
  rollout proceeds automatically to the next one.
- Check `partition` isn't intentionally holding the rollout back from an
  earlier canary or staged rollout you forgot about.
- If a specific ordinal's PVC is the blocker, resolve it like any other PVC
  issue, StatefulSet PVCs are per-ordinal and don't get recreated
  automatically if deleted.

## Related
- CrashLoopBackOff
- PVCPending
- ReadinessProbeFailed
