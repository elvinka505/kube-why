---
title: too many open files
aliases: [toomanyopenfiles, emfiletoomanyopenfiles]
category: runtime
---

# Too many open files

## What it means
The container process hit the maximum number of open file descriptors
allowed, either its own process limit, or the host's system-wide limit
inherited by the container. Once hit, the application can't open new files,
sockets, or connections, and typically starts failing in confusing,
seemingly unrelated ways.

## Common causes
1. The application leaks file descriptors, opening files or sockets without
   closing them, so usage climbs until the limit is hit.
2. A genuinely high-concurrency workload (many simultaneous connections)
   exceeds the default `ulimit` for open files, which is often set
   conservatively.
3. The host's system-wide file descriptor limit is itself too low for the
   combined load of everything running on it, not just this one container.
4. A container restarts frequently under load, and each restart's cleanup
   doesn't fully release descriptors before the next one starts.

## Check it
```
docker exec <container> sh -c 'ulimit -n'
cat /proc/<pid>/limits
```
Check the configured limit inside the container, then compare it to actual
usage in the application, if usage climbs steadily over the container's
uptime rather than starting high, that points to a leak rather than just an
undersized limit.

## Fix it
- Quick mitigation: raise the container's file descriptor limit:
  ```
  docker run --ulimit nofile=65536:65536 ...
  ```
- If it's a genuine leak, this only delays the problem, profile the
  application to find what's not closing file handles or sockets properly.
- If the host's system-wide limit is the actual ceiling, that needs
  raising at the host/kernel level (`/etc/security/limits.conf` or
  equivalent), not just at the container level.

## Related
- ContainerOOMKilled
