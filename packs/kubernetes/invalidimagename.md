---
title: InvalidImageName
aliases: [invalidimagename, invalid-image-name]
category: pod
---

# InvalidImageName

## What it means
The string in `spec.containers[].image` isn't a syntactically valid image
reference at all. This fails before Kubernetes even tries to contact a
registry, it's a parsing error, not a network or auth error.

## Common causes
1. A templating variable didn't get substituted (Helm/Kustomize placeholder like
   `${IMAGE_TAG}` left literally in the manifest).
2. Invalid characters in the tag (uppercase letters are not allowed in most tag
   contexts, spaces, or stray quotes from a shell script generating the YAML).
3. A malformed registry path, like a double slash or a missing tag separator.

## Check it
```
kubectl describe pod <pod>
```
The Events message quotes the exact invalid string, which usually makes the
problem obvious at a glance, look closely for unresolved template syntax or
stray characters.

## Fix it
- If it's an unresolved template variable, check your CI/CD pipeline step that's
  supposed to substitute it, the substitution step likely didn't run or ran
  against the wrong file.
- If it's a manual typo, fix the image string directly.
- Validate rendered manifests before applying them in CI (`helm template`,
  `kustomize build`) so this gets caught before it reaches the cluster.

## Related
- ImagePullBackOff
- CreateContainerConfigError
