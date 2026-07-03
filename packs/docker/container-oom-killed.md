---
title: Container killed with exit code 137 (OOM)
aliases: [exitcode137, containeroomkilled, dockeroomkilled]
category: runtime
---

# Container killed with exit code 137 (OOM)

## What it means
Exit code 137 means the container was killed by `SIGKILL` (128 + signal 9).
On its own that just means something forcibly killed it, but combined with
no application error in the logs and a memory limit set on the container,
it's almost always the kernel's OOM killer stepping in because the
container hit its `--memory` limit.

## Common causes
1. The container's `--memory` limit is set lower than what the application
   actually needs under real load.
2. No memory limit was set explicitly, but the host itself ran out of
   memory across everything running on it, and this container was the OOM
   killer's target.
3. A memory leak in the application, usage climbs steadily until it's
   killed, rather than being consistently high from the start.
4. A sudden spike in load or a batch operation pushes memory usage past a
   limit that's normally fine.

## Check it
```
docker inspect <container> --format '{{.State.OOMKilled}}'
docker stats <container>
```
`inspect` tells you definitively whether this specific exit was an OOM kill
versus some other cause of exit code 137. `docker stats` while it's running
shows real-time memory usage against its limit, useful for confirming
whether it's climbing steadily (leak) or spiking (load).

## Fix it
- Quick fix: raise the `--memory` limit to a realistic number based on
  actually observed usage, not a guess.
- Better fix: profile the application under production-like load to find
  its real working set before setting a limit.
- If it's a leak, raising the limit buys time but doesn't fix the
  underlying bug, treat it as a mitigation.
- If the host itself is out of memory across multiple containers, that's a
  capacity problem, not a single container's config, look at total memory
  allocated versus what the host actually has.

## Related
- NoSpaceLeftOnDevice
