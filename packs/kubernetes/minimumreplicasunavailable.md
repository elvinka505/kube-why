---
title: MinimumReplicasUnavailable
aliases: [minimumreplicasunavailable, deployment-not-available]
category: deployment
---

# MinimumReplicasUnavailable

## What it means
A Deployment's `Available` condition is False because fewer replicas are
ready than the minimum required (accounting for `maxUnavailable` during a
rollout). This is the Deployment reporting current unhealthiness, distinct
from ProgressDeadlineExceeded, which is specifically about a rollout timing
out, this condition can be true even outside an active rollout.

## Common causes
1. Currently running pods are crashing or failing readiness checks, so
   fewer are Ready than the Deployment expects.
2. A rollout is mid-flight and `maxUnavailable` is temporarily taking pods
   down faster than new ones are becoming ready, expected and transient if
   it resolves on its own.
3. Not enough cluster capacity to run the desired replica count at all, some
   replicas are stuck Pending rather than unhealthy.
4. A `PodDisruptionBudget` combined with node maintenance (drains, upgrades)
   is temporarily reducing available replicas below the minimum.

## Check it
```
kubectl describe deployment <name>
kubectl get pods -l <selector>
kubectl get pdb -n <namespace>
```
Look at the actual pod states, are they CrashLooping, Pending, or just
mid-rollout and progressing normally, that determines which underlying issue
you're actually chasing.

## Fix it
- If pods are crashing, this is really CrashLoopBackOff or
  ReadinessProbeFailed wearing a Deployment-level label, fix the pod issue
  directly.
- If it's a capacity problem, add nodes or reduce other scheduled load.
- If it's a transient rollout dip within expected `maxUnavailable` bounds, no
  action needed, confirm it resolves once the rollout completes.
- If a PDB combined with node drains is the cause, that's often expected
  behavior during planned maintenance, confirm the PDB's `minAvailable` still
  reflects what you actually need before assuming something's broken.

## Related
- ProgressDeadlineExceeded
- CrashLoopBackOff
- Pending
