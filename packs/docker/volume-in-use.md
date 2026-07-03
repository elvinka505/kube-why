---
title: volume is in use
aliases: [volumeisinuse, remove1or2volumeinuse]
category: storage
---

# volume is in use

## What it means
You tried to remove a volume that's still referenced by at least one
container, including stopped ones. Docker refuses the removal rather than
silently detaching it, to avoid a running (or restartable) container
suddenly losing its data out from under it.

## Common causes
1. A stopped container (not actively running, just not removed yet) still
   references the volume, `docker volume rm` checks references regardless
   of whether the container is currently running.
2. A `docker-compose` project was torn down with `docker-compose down`
   without the `-v` flag, which intentionally leaves volumes behind, then a
   later manual `docker volume rm` fails because a lingering container
   reference still exists.
3. Multiple containers share the same named volume, and only some of them
   were removed.

## Check it
```
docker ps -a --filter "volume=<volume-name>"
```
This lists every container, running or stopped, that references the
volume. That's the actual blocker, remove or confirm those containers
first before the volume itself can be removed.

## Fix it
- Remove the referencing containers first if they're genuinely done with:
  ```
  docker rm <container>
  docker volume rm <volume-name>
  ```
- For `docker-compose` projects, `docker-compose down -v` removes
  associated volumes as part of teardown, rather than needing a manual
  follow-up removal.
- If you're unsure whether a volume's data is still needed, back it up
  before removing anything, volume removal is not reversible.

## Related
- ContainerNameConflict
