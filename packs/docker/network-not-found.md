---
title: Network not found
aliases: [networknotfound, errorresponsefromdaemonnetworknotfound]
category: networking
---

# Network not found

## What it means
The container references a Docker network (via `--network` or a
`docker-compose` service definition) that doesn't currently exist. Docker
networks aren't created automatically the way they are the first time you
run `docker-compose up`, if the network was removed or never created,
anything referencing it by name fails immediately.

## Common causes
1. A network was removed with `docker network rm` or `docker network prune`,
   but containers or compose files still reference it by name.
2. Running a single service's container manually with `--network
   <project>_default` without having first brought up the full
   `docker-compose` project that creates that network.
3. A `docker-compose` file references an `external: true` network that was
   expected to already exist but was never actually created.
4. Restarting Docker itself (or the host) cleared networks that weren't
   defined to persist.

## Check it
```
docker network ls
docker network inspect <name>
```
Confirm whether the network actually exists under that exact name, this is
almost always a create-order or cleanup problem, not a configuration bug.

## Fix it
- If using `docker-compose`, bring the whole project up with `docker-compose
  up` rather than starting individual containers manually, it creates the
  network automatically as part of that.
- For an `external: true` network in compose, create it explicitly first:
  ```
  docker network create <name>
  ```
- If a prune command removed a network still in use, recreate it, and be
  more targeted with prune in the future (`docker network prune` removes
  all networks not used by at least one container, including ones you
  intended to keep around between runs).

## Related
- PortIsAlreadyAllocated
