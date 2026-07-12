---
title: Docker can't reach registry-1.docker.io
aliases: [dockercantreachregistry, registry1dockerionosuchhost, lookupregistry1dockerio, nosuchhost]
category: registry
---

# Docker can't reach registry-1.docker.io

## What it means
Docker reached the daemon successfully, but the daemon itself couldn't
resolve or reach Docker Hub over the network. This is different from
`cannot-connect-daemon`, where the CLI can't reach the daemon at all, and
from `pull-access-denied`, where the registry responds but rejects the
request.

## Common causes
1. DNS can't resolve `registry-1.docker.io`, often because the DNS server
   Docker or the host OS is using doesn't have a route to public DNS at all.
2. A corporate proxy, VPN, or firewall intercepts or blocks the request
   before it reaches Docker Hub, common on office networks.
3. A `dns` entry in `/etc/docker/daemon.json` overrides the host's working
   DNS with one that can't actually resolve public hostnames.
4. The machine is genuinely offline, or Docker Hub itself is down.

## Check it
```
docker pull hello-world
nslookup registry-1.docker.io
curl -v https://registry-1.docker.io/v2/
cat /etc/docker/daemon.json
```
If `nslookup` also fails, this is a DNS problem outside Docker, the host
itself can't resolve the name, not just the daemon. If `nslookup` succeeds
but `curl` hangs or is refused, that points to a proxy, VPN, or firewall
blocking the connection rather than DNS. Check `daemon.json` for a `dns` key,
if present, that's overriding whatever DNS is actually working for the rest
of the host.

## Fix it
- DNS not resolving at the host level: fix the host's DNS configuration
  (`/etc/resolv.conf` or your OS network settings) before touching Docker at
  all, Docker inherits whatever DNS the host provides unless overridden.
- Corporate proxy/VPN/firewall: configure Docker's proxy settings
  (`~/.docker/config.json` or `/etc/systemd/system/docker.service.d/` on
  Linux) to route through the corporate proxy, or ask IT to whitelist
  `registry-1.docker.io` and `auth.docker.io`.
- Bad `dns` override in `daemon.json`: remove or correct the `dns` entry,
  then `sudo systemctl restart docker` for it to take effect.
- Genuinely offline or Docker Hub down: wait it out, or point at a mirror/
  pull-through cache if this needs to keep working without direct internet
  access.

## Related
- Cannot connect to the Docker daemon
- Pull access denied / manifest unknown