---
title: version in docker-compose.yml is unsupported
aliases: [versionindockercomposeisunsupported, composefileversionunsupported]
category: build
---

# Version in docker-compose.yml is unsupported

## What it means
The `version:` field at the top of a `docker-compose.yml` declares which
Compose file format the rest of the file follows. If it names a version
your installed Compose doesn't recognize, or is malformed, Compose refuses
to parse the rest of the file at all.

## Common causes
1. An old Compose file using a `version` string that newer Compose releases
   have dropped support for entirely.
2. A typo in the version string (`"3,8"` instead of `"3.8"`).
3. Mixing syntax from a newer Compose spec (which, in the latest Compose
   versions, doesn't require a `version` key at all) with an explicitly
   declared older version, causing a mismatch between what the file expects
   and what's actually written below it.
4. Using Compose V1 (`docker-compose`, hyphenated) syntax expectations
   against Compose V2 (`docker compose`, a Docker CLI plugin), which handles
   versioning slightly differently.

## Check it
```
docker compose version
docker compose config
```
`docker compose config` renders and validates the file without starting
anything, it surfaces the exact parsing complaint clearly, and confirms
what Compose version you actually have installed to run it.

## Fix it
- Fix the typo or update the `version` field to one your installed Compose
  actually supports.
- For modern Compose (V2), consider dropping the `version` key entirely,
  current Compose infers format from content and this removes an entire
  class of this error.
- If migrating from `docker-compose` (V1) to `docker compose` (V2), check
  the migration notes for syntax differences rather than assuming full
  compatibility.

## Related
- ComposeDependsOnNotReady
