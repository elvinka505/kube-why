---
title: permission denied on a bind-mounted volume
aliases: [permissiondeniedbindmount, eaccespermissiondeniedvolume, bindmountpermissiondenied]
category: storage
---

# Permission denied on a bind-mounted volume

## What it means
The application inside the container can see the bind-mounted host
directory, but can't read or write specific files in it, because the
UID/GID the container process runs as doesn't match what the host
filesystem's permissions expect. Bind mounts don't remap ownership, the
container sees the exact same UID numbers the host does.

## Common causes
1. The container runs as a UID (often root, `0`, or a fixed application UID
   baked into the image) that doesn't match the UID owning the files on the
   host side of the mount.
2. Files were created on the host by one user, then a container running as
   a different UID tries to write to them.
3. Files were created *by* a previous container run as root, and now a
   later run as a non-root user can't touch them.
4. SELinux or similar mandatory access control on the host blocking the
   container's access entirely, independent of standard UNIX permissions.

## Check it
```
ls -ln <host-path>
docker inspect <container> --format '{{.Config.User}}'
```
`ls -ln` shows the numeric UID/GID owning the files on the host, compare it
directly against the UID the container actually runs as, a mismatch here is
almost always the root cause.

## Fix it
- Align the container's UID with the host directory's ownership, either by
  running the container with `--user <uid>:<gid>` matching the host files,
  or by `chown`-ing the host directory to match the container's expected
  UID.
- Avoid mixing root-created and non-root-created files in the same
  bind-mounted directory across different runs, pick one UID convention and
  stick to it for that volume.
- If SELinux is enforcing, consider the `:z`/`:Z` mount flags to relabel the
  volume appropriately, understand the difference (`:z` shared, `:Z`
  private) before choosing.

## Related
- ReadOnlyFileSystem
