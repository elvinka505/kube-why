---
title: executable file not found in $PATH
aliases: [executablefilenotfoundinpath, ocirunimeexecfailedexecutablenotfound]
category: runtime
---

# Executable file not found in $PATH

## What it means
You tried to `docker exec` into a running container using a shell
(`bash`, often) that simply doesn't exist inside that image. This is
distinct from `exec format error`, the binary you're asking for isn't
present at all, it's not an architecture mismatch.

## Common causes
1. The image is based on a minimal or distroless base (Alpine, `scratch`,
   Google's distroless images) that either has no shell at all, or has
   `sh` but not `bash` specifically.
2. Habitually running `docker exec -it <container> bash` out of muscle
   memory without checking what shell (if any) the specific image actually
   ships.
3. A multi-stage build's final stage copies only the compiled binary into a
   minimal runtime image, deliberately excluding a shell for size and
   security reasons.

## Check it
```
docker exec <container> which bash
docker exec <container> which sh
docker inspect <image> --format '{{.Os}}/{{.Architecture}}'
```
Check whether any shell exists at all before assuming this is an
architecture problem, minimal images often have no shell by design, this
isn't a bug in the image.

## Fix it
- Try `sh` instead of `bash`, many minimal images include a POSIX shell
  even without full bash:
  ```
  docker exec -it <container> sh
  ```
- For truly shell-less images (`scratch`, distroless), you generally can't
  exec a shell in at all, use `docker cp` to inspect files, or attach a
  debug sidecar/ephemeral container that shares the target's namespaces
  instead.
- If you need routine shell access for debugging, consider a debug variant
  of the image for development that includes a shell, while keeping the
  production image minimal.

## Related
- ExecFormatError
