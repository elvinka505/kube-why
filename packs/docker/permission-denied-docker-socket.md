---
title: Permission denied while connecting to the Docker socket
aliases: [permissiondeniedwhiletryingtoconnecttothedockerdaemonsocket, gotpermissiondeniedvarrundockersock, permissiondenieddockersocket]
category: daemon
---

# Permission denied while connecting to the Docker socket

## What it means
This is a more specific version of "cannot connect to the Docker daemon,"
the daemon is running and reachable, but your current user doesn't have
permission to read/write its Unix socket, so every command fails with a
permission error instead of a connection error.

## Common causes
1. Your user was never added to the `docker` group, so only root (or users
   in that group) can access `/var/run/docker.sock` on Linux.
2. You were added to the `docker` group, but haven't started a new login
   session since, group membership changes don't apply to already-running
   shells.
3. Running Docker commands inside a container or CI job that mounts the
   host's Docker socket, but the process inside runs as a user whose UID
   doesn't match what the socket's permissions expect.
4. A misconfigured Docker installation where the socket's group ownership
   is wrong entirely.

## Check it
```
ls -l /var/run/docker.sock
groups
```
Check the socket's owning group (usually `docker`) against your current
user's groups. If `docker` isn't listed in your `groups` output, that's
confirmed.

## Fix it
- Add your user to the group and start a fresh session for it to apply:
  ```
  sudo usermod -aG docker $USER
  ```
  then log out and back in (or open a new terminal, depending on your OS),
  a change to group membership doesn't retroactively apply to shells opened
  before the change.
- For containers mounting the host socket, ensure the process inside runs
  as a UID that's part of the socket's group, or as root if that's
  acceptable for your use case.
- As an immediate but less ideal workaround, prefix commands with `sudo`.

## Related
- CannotConnectToTheDockerDaemon
