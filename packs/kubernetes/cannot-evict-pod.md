---
title: Cannot evict pod
aliases: [cannot-evict-pod, cannot-evict, violate-the-pods-disruption-budget]
category: deployment
---

# Cannot evict pod

## What it means
The eviction API refused to evict the pod because doing so would violate its
`PodDisruptionBudget`. This is different from `Evicted`, where the kubelet has
already removed the pod because the node itself was under resource pressure.
Here, the pod stays running because Kubernetes won't allow another voluntary
disruption until more replicas become available.

## Common causes
1. You're draining a node (`kubectl drain`) or the cluster autoscaler is
   scaling a node down, but another replica is already unavailable, so the
   PodDisruptionBudget won't allow one more voluntary disruption.
2. Another pod is CrashLooping, Pending, or failing readiness checks,
   exhausting the disruption budget before the drain even starts.
3. The Deployment doesn't have enough replicas to satisfy the PDB's
   `minAvailable` requirement during maintenance.
4. The PodDisruptionBudget is stricter than the application's replica count
   allows (for example, `minAvailable` equals the total replica count).

## Check it
```
kubectl describe pod <pod>
kubectl get pdb
kubectl describe pdb <name>
kubectl get pods -l <selector>
```

The Events section usually contains the exact eviction error. Then compare the
number of Ready pods with the PDB's `minAvailable` or `maxUnavailable`
settings. If one replica is already unavailable, that's often enough to exhaust
the disruption budget and block another eviction.

## Fix it
- If another pod is unhealthy, fix that underlying issue first (it's usually
  CrashLoopBackOff, Pending, or ReadinessProbeFailed). The
  PodDisruptionBudget is just preventing the drain from making availability
  even worse.
- If the Deployment doesn't have enough replicas, scale it up so the PDB can
  still be satisfied during maintenance.
- If the PDB is stricter than necessary, review `minAvailable` or
  `maxUnavailable` and adjust them to match the application's real
  availability requirements.
- For planned maintenance, temporarily relax the PDB, complete the drain,
  then restore the original settings afterwards.

## Related
- MinimumReplicasUnavailable
- Evicted
- Pending