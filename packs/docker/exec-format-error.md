---
title: exec format error
aliases: [execformaterror, standard_initlinuxgostartcontainerbin, ocirunimeexecfailed]
category: runtime
---

# exec format error

## What it means
The container tried to run a binary built for a different CPU architecture
than the one it's actually running on, most commonly an image built for
arm64 (Apple Silicon) being run on an amd64 host, or vice versa. The kernel
can't execute it at all, this fails before your application code runs.

## Common causes
1. Building an image on an Apple Silicon Mac and pushing/running it on an
   amd64 server (or a CI runner) without building for the right target
   architecture.
2. Pulling a base image that only publishes one architecture, and it isn't
   the one you're running on.
3. A multi-stage Dockerfile copies a binary built in one stage for the
   wrong architecture into the final stage.
4. Using a pre-built binary release (downloaded in the Dockerfile) that
   doesn't match the target architecture.

## Check it
```
docker inspect <image> --format '{{.Architecture}}'
uname -m
```
Compare the image's built architecture against the host's actual
architecture (`uname -m` on the machine actually running the container).
A mismatch here confirms this exact cause.

## Fix it
- Build multi-architecture images with `docker buildx build --platform
  linux/amd64,linux/arm64`, so the same tag works correctly regardless of
  where it's pulled and run.
- If you only need one target, explicitly specify `--platform` at build
  time to match the deployment target rather than relying on your local
  machine's default.
- For downloaded binaries in a Dockerfile, make sure the download URL
  selects the architecture correctly rather than hardcoding one.

## Related
- PullAccessDenied
