---
title: RunContainerError
aliases: [runcontainererror, failed-to-start-container]
category: container-runtime
---

# RunContainerError

## What it means
The image was pulled successfully, but the container runtime (containerd,
CRI-O) failed to actually start it. This sits between a successful image pull
and the application ever running, the failure is at the runtime/OS level, not
inside your application code yet.

## Common causes
1. The container's command/entrypoint doesn't exist in the image, or isn't
   executable (missing execute permission, wrong shebang, built for the
   wrong CPU architecture entirely).
2. A mounted volume can't be attached in the way the container expects
   (permission mismatch, path conflict with something already in the image).
3. The container requests a Linux capability, security profile, or resource
   the node's runtime/kernel doesn't support or allow.
4. `securityContext` settings that conflict with each other or with what the
   image expects (running as a non-root user that doesn't exist inside the
   image, for example).

## Check it
```
kubectl describe pod <pod>
kubectl logs <pod>
```
The Events section usually has the runtime's specific error message, "exec
format error" points to an architecture mismatch, "permission denied"
points to executable bits or security context, "no such file or directory"
points to a wrong entrypoint path.

## Fix it
- Architecture mismatch: rebuild or pull the image for the correct
  architecture (a common trap when building on Apple Silicon and deploying
  to amd64 nodes, or vice versa).
- Wrong entrypoint path: verify the command exists at that exact path inside
  the image, `docker run --entrypoint sh <image>` locally to check directly.
- Permission issues: align the image's expected user with the pod's
  `securityContext`, or adjust the image if it was built assuming root.
- Volume/mount conflicts: check for path collisions between mounted volumes
  and files already present in the image at the same path.

## Related
- CrashLoopBackOff
- ContainerCreating stuck
