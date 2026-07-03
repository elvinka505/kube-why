---
title: container can't reach the host machine
aliases: [containercantreachhostmachine, hostdockerinternalnotresolving]
category: networking
---

# Container can't reach the host machine

## What it means
A process inside a container needs to reach something running directly on
the host machine (a local dev server, a database running outside Docker),
and the usual way to do that, `host.docker.internal`, either isn't
resolving or isn't behaving the way it does on other platforms. This varies
meaningfully between Docker Desktop (Mac/Windows) and native Linux Docker,
which is the root of most confusion here.

## Common causes
1. On native Linux (not Docker Desktop), `host.docker.internal` doesn't
   resolve by default, it's a Docker Desktop convenience feature, Linux
   requires explicit configuration to get equivalent behavior.
2. The host's local server is bound to `localhost`/`127.0.0.1` only, which
   is not reachable from inside a container regardless of the networking
   setup, it needs to bind to `0.0.0.0` to accept connections from
   anywhere, including containers.
3. A host firewall blocking the container's bridge network from reaching
   host-bound ports.

## Check it
```
docker exec <container> getent hosts host.docker.internal
```
Confirm whether the hostname resolves at all first, on native Linux this
commonly fails unless explicitly configured, which immediately narrows down
the cause.

## Fix it
- On native Linux, add the host gateway explicitly:
  ```
  docker run --add-host=host.docker.internal:host-gateway ...
  ```
  (Docker Desktop provides this automatically, native Linux doesn't.)
- Make sure whatever's running on the host binds to `0.0.0.0`, not just
  `localhost`, or containers can never reach it no matter how networking is
  configured.
- Check host firewall rules if the hostname resolves correctly but
  connections still time out.

## Related
- NetworkNotFound
- DriverFailedProgrammingExternalConnectivity
