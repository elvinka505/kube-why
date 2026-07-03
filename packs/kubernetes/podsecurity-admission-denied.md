---
title: Pod Security admission denied
aliases: [podsecurityadmissiondenied, psadenied, violates-podsecurity]
category: podsecurity
---

# Pod Security admission denied

## What it means
The namespace enforces a Pod Security Standard level (`privileged`,
`baseline`, or `restricted`, set via a namespace label), and your pod's spec
violates it, most often by requesting privileges the `restricted` level
disallows. The built-in admission controller rejects it before it's ever
created.

## Common causes
1. The pod or a container runs as root (no `runAsNonRoot: true`, or no
   explicit non-root `runAsUser`), which `restricted` disallows.
2. A container requests `privileged: true`, added Linux `capabilities`
   beyond what's allowed, or `allowPrivilegeEscalation: true`.
3. `hostNetwork`, `hostPID`, or `hostPath` volumes are used, all disallowed
   under `baseline` and `restricted`.
4. The namespace's enforcement level was tightened (often cluster-wide, by a
   platform team) after workloads already running there were written against
   a looser standard.

## Check it
```
kubectl get namespace <namespace> --show-labels
kubectl describe pod <pod>
```
Check the namespace's `pod-security.kubernetes.io/enforce` label to see which
level applies. The rejection message from `kubectl apply` names the specific
field that violated the policy, it's usually precise enough to fix directly.

## Fix it
- Add the required `securityContext` fields: `runAsNonRoot: true`, drop all
  capabilities and add back only what's genuinely needed, remove
  `privileged`/`allowPrivilegeEscalation` unless there's a real reason for
  them.
- If the workload genuinely needs elevated privileges (a node-level agent,
  for example), it likely belongs in a namespace enforcing `privileged` or
  `baseline`, not `restricted`, that's a placement decision, not something to
  bypass.
- Coordinate with whoever owns the namespace's security label if the level
  was tightened cluster-wide and multiple workloads need updating, this is
  rarely a one-pod fix in practice.

## Related
- Forbidden
- AdmissionWebhookDenied
