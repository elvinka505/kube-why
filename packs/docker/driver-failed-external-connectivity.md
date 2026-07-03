---
title: driver failed programming external connectivity
aliases: [driverfailedprogrammingexternalconnectivity, failedprogrammingexternalconnectivity]
category: networking
---

# driver failed programming external connectivity

## What it means
Docker successfully decided to bind a port for your container, but failed
while actually setting up the underlying `iptables`/netfilter rules that
make that binding work. This is a layer deeper than "port already
allocated," the port itself might be free, the failure is in the host's
network rule programming.

## Common causes
1. `iptables` rules were manually modified or flushed outside of Docker
   (by a firewall tool, a security agent, or a manual `iptables -F`),
   leaving Docker's own chains in an inconsistent state.
2. A conflicting firewall tool (firewalld, ufw) is fighting with Docker
   over the same `iptables` chains.
3. The Docker daemon was restarted while containers were still running,
   and its in-memory view of network state didn't reconcile properly with
   what's actually on the host.
4. Running Docker inside a restricted environment (a VM, a CI runner)
   without the necessary kernel modules or privileges for netfilter.

## Check it
```
sudo iptables -t nat -L DOCKER
docker network inspect bridge
```
Confirm Docker's own `DOCKER` chain in the `nat` table actually exists and
looks sane. If it's empty or missing entirely, something outside Docker
altered it.

## Fix it
- Restart the Docker daemon, this often lets it reprogram its `iptables`
  rules from scratch and resolves transient inconsistency:
  ```
  sudo systemctl restart docker
  ```
- If a competing firewall tool is the cause, configure it to not manage the
  same chains Docker uses, rather than disabling one or the other blindly.
- In CI/restricted environments, confirm the runner actually has the
  privileges and kernel modules Docker's networking needs, this is a common
  gap in locked-down containerized CI runners.

## Related
- PortIsAlreadyAllocated
- NetworkNotFound
