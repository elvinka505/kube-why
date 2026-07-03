---
title: Docker Desktop is unable to start / the VM is not running
aliases: [dockerdesktopunabletostart, dockerdesktopvmnotrunning]
category: daemon
---

# Docker Desktop is unable to start / the VM is not running

## What it means
On Mac and Windows, Docker Desktop runs containers inside a lightweight
virtual machine, the `docker` CLI you use talks to a daemon running inside
that VM, not directly to your host OS. When the VM itself fails to start,
every Docker command fails, often with a misleadingly generic connection
error, because the whole backend it depends on isn't up.

## Common causes
1. A previous crash or forced shutdown left the VM's internal state
   corrupted.
2. Insufficient resources allocated to Docker Desktop (memory/disk) for the
   VM to start, especially after other applications have claimed more of
   the host's resources.
3. A conflicting hypervisor (another virtualization tool using the same
   virtualization framework) preventing Docker Desktop's VM from
   initializing.
4. A pending OS update or a Docker Desktop update that didn't fully
   complete.

## Check it
```
docker info
```
On Docker Desktop, check the Desktop application's own diagnostic/log
viewer (accessible from its UI), the CLI alone often just reports
generic connection failures without surfacing the VM-level cause.

## Fix it
- Restart Docker Desktop entirely, not just retry the failing command, this
  resolves the majority of transient VM-start failures.
- If restarting doesn't help, Docker Desktop's own "reset to factory
  defaults" option resolves corrupted VM state, understand this removes all
  local images/containers/volumes, back up anything you need first.
- Check allocated resources in Docker Desktop settings if the host was
  under heavy memory pressure when the failure started.
- For hypervisor conflicts, check whether another virtualization tool was
  running at the same time and whether both can coexist on this OS version.

## Related
- CannotConnectToTheDockerDaemon
