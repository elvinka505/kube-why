---
title: service starts before its dependency is actually ready
aliases: [servicestartsbeforedependencyisready, dependsonnotwaitingforready, composedependsonnotready]
category: build
---

# Service starts before its dependency is actually ready

## What it means
`depends_on` in Compose controls start *order*, it does not wait for the
dependency to be *ready* to accept connections, only for its container
process to have started (or, with a healthcheck condition, for the
healthcheck to pass, if one is actually configured). Without a healthcheck
condition, a database container can report "started" long before it's
actually accepting connections, and the dependent service crashes trying
to connect too early.

## Common causes
1. `depends_on` is used without a `condition`, so Compose only guarantees
   the dependency's container started, not that it's ready to serve
   traffic.
2. The dependency image itself has a slow startup/initialization phase
   (a database running migrations, a cache warming up) that outlasts
   however long the dependent service waits before its first connection
   attempt.
3. No healthcheck is defined on the dependency at all, so Compose has no
   way to know "ready" even if you wanted it to wait for that.

## Check it
```
docker compose logs <dependency-service>
docker inspect <dependency-container> --format '{{json .State.Health}}'
```
Check whether the dependency even has a healthcheck configured at all, if
`.State.Health` is empty, that confirms `depends_on` can only ever wait for
"started," not "ready," for this service as currently configured.

## Fix it
- Add a `healthcheck` to the dependency and use `depends_on` with
  `condition: service_healthy` so Compose actually waits for readiness, not
  just process start:
  ```yaml
  depends_on:
    db:
      condition: service_healthy
  ```
- If you can't add a healthcheck, add retry/backoff logic in the dependent
  application itself so it doesn't crash on the first failed connection
  attempt, this is more robust than relying on startup ordering alone
  regardless.

## Related
- ComposeVersionUnsupported
