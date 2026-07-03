---
title: Container name already in use
aliases: [containernameconflict, thecontainernameisalreadyinuse, conflictnamealreadyinuse, isalreadyinusebycontainer]
category: container
---

# Container name already in use

## What it means
You tried to create a container with a name that's already taken by another
container, including a stopped one. Docker container names must be unique
regardless of whether the existing container is running.

## Common causes
1. A previous run of the same container (from a script, `docker-compose`, or
   a manual `docker run --name`) exited or crashed but was never removed.
2. Re-running a script or CI job that doesn't clean up its container from
   the last run before starting a new one.
3. Two different projects or scripts happen to use the same hardcoded
   container name.

## Check it
```
docker ps -a --filter "name=<name>"
```
`-a` matters here, a stopped or exited container still holds the name, a
plain `docker ps` (running only) won't show it, and you'll be confused about
why the name is "already in use" when nothing looks like it's running.

## Fix it
- Remove the stale container if it's genuinely done with:
  ```
  docker rm <name>
  ```
- For scripts/CI that repeatedly create and destroy the same named
  container, add `--rm` to `docker run` so it cleans up automatically on
  exit, or explicitly `docker rm -f <name>` before creating a new one.
- If two unrelated things collided on the same name, rename one of them, or
  stop hardcoding names and let Docker generate one.

## Related
- CannotConnectToTheDockerDaemon
