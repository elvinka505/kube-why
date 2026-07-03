---
title: No space left on device
aliases: [nospaceleftondevice, dockerdiskfull, outofdiskspacedocker]
category: storage
---

# No space left on device

## What it means
Docker's storage (image layers, container writable layers, volumes, build
cache) has filled up the disk it lives on. This can happen even when `df -h`
on your actual project directory looks fine, Docker often keeps its data in
a separate location (a VM disk image, on Docker Desktop) that fills up
independently of your main filesystem.

## Common causes
1. Old, unused images accumulating over time, especially from frequent
   rebuilds during development where each layer change creates new
   intermediate images.
2. Stopped containers and their writable layers never getting cleaned up.
3. Build cache growing unbounded from repeated `docker build` runs.
4. On Docker Desktop specifically, the VM's disk image has a fixed size
   allocation that fills up independently of your host machine's actual free
   space.

## Check it
```
docker system df
docker system df -v
```
This shows exactly where the space is going, images, containers, local
volumes, and build cache, broken out separately, rather than guessing.

## Fix it
- Quick reclaim: `docker system prune` removes stopped containers, unused
  networks, dangling images, and build cache. Add `-a` to also remove
  unused (not just dangling) images, and `--volumes` to include unused
  volumes, understand `--volumes` can delete data you actually wanted to
  keep, check first.
- For Docker Desktop specifically, increase the VM's disk image size limit
  in settings if pruning doesn't free enough, or if this recurs frequently.
- In CI, add a prune step on a schedule rather than letting it accumulate
  indefinitely across every build.

## Related
- ContainerNameConflict
