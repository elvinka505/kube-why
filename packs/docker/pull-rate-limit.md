---
title: toomanyrequests, pull rate limit exceeded
aliases: [toomanyrequests, pullratelimitexceeded, youhavereachedyourpullratelimit]
category: registry
---

# toomanyrequests, pull rate limit exceeded

## What it means
Docker Hub rate-limits anonymous (and, at a higher threshold,
authenticated free-tier) image pulls per six-hour window, per IP address or
account. Once you exceed it, every pull fails with this error until the
window resets, regardless of whether the image itself is public and
otherwise pullable.

## Common causes
1. Anonymous pulls from a shared IP (a CI runner, an office network, a
   Kubernetes node pulling the same base image across many pods) hit the
   limit collectively, not per-user.
2. CI pipelines that pull the same base images on every single run without
   any layer caching, multiplying pull volume unnecessarily.
3. A busy Kubernetes cluster with many nodes all pulling from Docker Hub
   anonymously, effectively sharing one rate-limit bucket per egress IP.

## Check it
```
docker pull hello-world
```
The error message itself usually states the limit and reset window
directly. If unauthenticated pulls fail with this but authenticated ones
don't, that confirms it's specifically about being anonymous, not a
broader account or network issue.

## Fix it
- Authenticate your pulls (`docker login`) even for public images, logged-in
  pulls get a meaningfully higher rate limit than anonymous ones.
- Mirror frequently-used base images to a registry you control (or a pull-
  through cache) instead of hitting Docker Hub directly on every pull.
- In CI, cache image layers between runs so you're not re-pulling identical
  base images from scratch every time.
- For Kubernetes clusters, configure `imagePullSecrets` cluster-wide so
  every node authenticates rather than pulling anonymously and sharing one
  rate-limit bucket across the whole cluster.

## Related
- PullAccessDenied
