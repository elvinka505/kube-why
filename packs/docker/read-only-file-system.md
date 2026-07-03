---
title: read-only file system
aliases: [readonlyfilesystem, readonlyfilesystemerror]
category: runtime
---

# read-only file system

## What it means
The application inside the container tried to write to a path that's
mounted read-only, and the write was rejected at the filesystem level. This
is usually intentional hardening (`--read-only` on the container, or a
read-only bind mount), the error means the application wasn't written with
that constraint in mind, not that anything is broken.

## Common causes
1. The container was explicitly run with `--read-only` for security
   hardening, but the application writes logs, cache files, or temp files
   to a path that wasn't given a writable volume.
2. A bind-mounted host directory was mounted `:ro` intentionally, but the
   application inside expects to write to it.
3. The application writes to a location like `/tmp` assuming it's always
   writable, but the container's filesystem layering makes it read-only in
   this specific setup.
4. Copying `--read-only` container configuration from another project
   without auditing what this specific application actually needs to write
   to.

## Check it
```
docker inspect <container> --format '{{.HostConfig.ReadonlyRootfs}}'
docker inspect <container> --format '{{json .Mounts}}'
```
Confirm whether the root filesystem is genuinely read-only, and check each
mount's read/write mode, this tells you exactly which path the application
is trying to write to that it isn't allowed to.

## Fix it
- Add a writable volume or tmpfs mount for the specific path the
  application needs (logs, cache, temp files), rather than disabling
  `--read-only` entirely and losing the hardening benefit:
  ```
  docker run --read-only --tmpfs /tmp ...
  ```
- If a bind mount was meant to be writable, drop the `:ro` suffix.
- If the application's write path is configurable, redirect it to a
  location you've explicitly made writable instead of wherever it defaults
  to.

## Related
- BindMountPermissionDenied
