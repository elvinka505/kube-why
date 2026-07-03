---
title: CreateContainerConfigError
aliases: [createcontainerconfigerror, container-config-error]
category: pod
---

# CreateContainerConfigError

## What it means
The kubelet has the image and is ready to start the container, but the config
it needs to build the container is invalid, usually because it references a
ConfigMap, Secret, or key that doesn't exist. The container never actually
starts.

## Common causes
1. `envFrom` or `env.valueFrom` points to a ConfigMap or Secret that doesn't
   exist in this namespace.
2. The ConfigMap/Secret exists but the specific `key` referenced inside it
   doesn't.
3. The resource exists in a different namespace than the pod (ConfigMaps and
   Secrets are namespace-scoped, references don't cross namespaces).
4. A typo in the resource name in the pod spec.

## Check it
```
kubectl describe pod <pod>
```
The Events section names the missing reference directly, something like
`couldn't find key API_KEY in Secret app/db-creds`. Then confirm it actually
exists:
```
kubectl get configmap <name> -n <namespace>
kubectl get secret <name> -n <namespace> -o yaml
```

## Fix it
- Create the missing ConfigMap/Secret, or fix the name/key referenced in the pod
  spec to match what actually exists.
- If this is happening after a fresh deploy to a new namespace, check whether
  your deployment pipeline creates the ConfigMap/Secret before the workload, or
  assumes it already exists.
- Double check you're not confusing a Secret's key name with its literal value,
  `envFrom`/`valueFrom` need the key name, not the decoded content.

## Related
- CrashLoopBackOff
- InvalidImageName
