---
title: ImagePullBackOff
aliases: [imagepullbackoff, errimagepull, image-pull-backoff]
category: pod
---

# ImagePullBackOff / ErrImagePull

## What it means
The kubelet on the node can't pull the container image you asked for.
`ErrImagePull` is the first failed attempt, `ImagePullBackOff` is Kubernetes
backing off and retrying after repeated failures. Same root problem, different
stage.

## Common causes
1. Typo in the image name or tag (`myapp:latst` instead of `myapp:latest`).
2. The image or tag genuinely doesn't exist in the registry.
3. The image is in a private registry and the pod has no `imagePullSecrets`, or
   the secret is wrong/expired.
4. The node can't reach the registry at all (network policy, proxy, DNS, or a
   registry outage).
5. Rate limiting from the registry (Docker Hub anonymous pulls are a frequent
   culprit).

## Check it
```
kubectl describe pod <pod>
```
The Events section will usually spell out the exact failure: `manifest unknown`,
`unauthorized`, `no such host`, `429 Too Many Requests`. Read that line before
doing anything else, it tells you which of the causes above you're dealing with.

## Fix it
- Typo or missing tag: fix the manifest, redeploy.
- Private registry: confirm the secret exists in the right namespace and is
  referenced in the pod spec or the ServiceAccount:
  ```
  kubectl get secret <name> -n <namespace>
  kubectl get sa default -n <namespace> -o yaml
  ```
- Network/DNS issue: exec into another pod on the same node and try to reach the
  registry directly to isolate node vs. cluster-wide problems.
- Rate limiting: authenticate your pulls instead of pulling anonymously, or mirror
  the image to a registry you control.

## Related
- CrashLoopBackOff
- InvalidImageName
