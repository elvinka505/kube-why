---
title: Cannot connect to the Docker daemon
aliases: [cannotconnecttothedockerdaemon, dockerdaemonnotrunning, isthedockerdaemonrunning]
category: daemon
---

# Cannot connect to the Docker daemon

## What it means
The `docker` CLI is a client, it talks to a background daemon (`dockerd`) over
a socket to actually do anything. This error means the client can't reach
that daemon at all, before any question of permissions or the command itself
comes up.

## Common causes
1. Docker (or Docker Desktop, on Mac/Windows) simply isn't running.
2. The daemon is running, but your user doesn't have permission to access its
   socket, `/var/run/docker.sock` on Linux is usually owned by a `docker`
   group your user isn't in.
3. `DOCKER_HOST` is set to point at a remote or non-default socket that isn't
   actually reachable.
4. On Linux, the docker service crashed or was never started
   (`systemctl status docker` would show this).

## Check it
```
docker info
systemctl status docker
echo $DOCKER_HOST
```
`docker info` gives a clearer error than most commands, if it fails the same
way, this is confirmed to be a daemon connectivity issue, not something
specific to whatever command you were originally trying to run.

## Fix it
- Not running: start it (`sudo systemctl start docker` on Linux, or open
  Docker Desktop on Mac/Windows).
- Permission denied specifically: add your user to the `docker` group
  (`sudo usermod -aG docker $USER`) and start a new shell session for it to
  take effect, or prefix commands with `sudo` as a quicker but less
  convenient workaround.
- Stray `DOCKER_HOST`: unset it unless you specifically meant to point at a
  remote daemon.

## Related
- PermissionDeniedDockerSocket
- DockerCantReachRegistry