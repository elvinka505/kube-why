---
title: Port is already allocated
aliases: [portisalreadyallocated, bindaddressalreadyinuse, portalreadyinuse]
category: networking
---

# Port is already allocated

## What it means
Docker tried to bind a host port for your container (`-p hostPort:containerPort`)
but something is already listening on that host port, could be another
container, or a completely unrelated process on the machine itself.

## Common causes
1. Another container (often a previous, not-fully-stopped run of the same
   project) is already bound to that host port.
2. A non-Docker process on the host (a local dev server, another service)
   is already using it.
3. A `docker-compose` project was brought up twice, or an old instance from
   a different terminal/session was never brought down.
4. Restarting a container quickly after stopping it, before the OS has
   finished releasing the port from the previous process.

## Check it
```
docker ps --filter "publish=<port>"
lsof -i :<port>          # macOS/Linux, shows any process, not just Docker
```
`docker ps --filter publish=` narrows it to Docker containers specifically,
`lsof` (or `netstat`/`ss` on Linux) tells you if it's actually a non-Docker
process hogging the port, which changes what you need to do next.

## Fix it
- Another container holding it: stop or remove that container, or pick a
  different host port for the new one.
- A non-Docker process: stop that process, or use a different host port for
  the container.
- Restarted too quickly after stopping: wait a moment and retry, or bind
  explicitly to confirm it's actually free first.
- For local development, consider letting Docker pick a random free host
  port (`-p containerPort` without a fixed host side, or `docker port
  <container>` to see what it picked) instead of hardcoding one that
  conflicts across projects.

## Related
- ContainerNameConflict
