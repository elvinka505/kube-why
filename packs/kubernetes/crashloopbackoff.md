---
title: CrashLoopBackOff
aliases: [crashloop, crash-loop-backoff, crashloopbackoff]
category: pod
---

# CrashLoopBackOff

## What it means
Your container starts, exits, and Kubernetes keeps restarting it with a growing
delay between attempts. CrashLoopBackOff is not the bug itself, it's Kubernetes
telling you the bug is happening repeatedly. The real error is almost always in
the previous container's logs.

## Common causes
1. The app crashes on startup (bad config, missing env var, unhandled exception).
2. The container's main process exits immediately because it isn't meant to run
   in the foreground (common with images that expect a shell wrapper).
3. A liveness probe is too aggressive and kills the container before it finishes
   booting.
4. The app depends on something not ready yet (database, another service) and
   has no retry logic, so it just dies.

## Check it
```
kubectl logs <pod> --previous
kubectl describe pod <pod>
```
`--previous` matters. By the time you look, the current container may not have
crashed yet, but the one before it did, and that's the one with the real error.
Check the Events section at the bottom of `describe pod` too, it often names the
exact reason (OOMKilled, Error, non-zero exit code).

## Fix it
- Read the actual crash log first, don't guess.
- If the process exits instantly with no log output at all, check the image's
  entrypoint/cmd, it may need `-f` or a foreground flag.
- If it's a liveness probe killing a slow-starting app, add `initialDelaySeconds`
  or switch to a `startupProbe`.
- If it depends on another service, add retry/backoff logic in the app, or an
  init container that waits for the dependency.

## Related
- OOMKilled
- ImagePullBackOff
- ReadinessProbeFailed
