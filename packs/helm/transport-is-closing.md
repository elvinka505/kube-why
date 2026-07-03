---
title: failed to fetch chart, transport is closing
aliases: [failedtofetchchart, transportisclosing, lookingupchartinformationerror]
category: chart
---

# Failed to fetch chart / transport is closing

## What it means
Helm couldn't retrieve the chart itself from its configured repository,
this fails before anything about your values or templates comes into play,
Helm never even got the chart's contents to work with.

## Common causes
1. The chart repository was never added, or was added but its index was
   never refreshed (`helm repo update`), so Helm's local cache doesn't know
   about this chart or chart version.
2. The specific chart version requested doesn't exist in the repository
   (a typo'd version, or a version that was never actually published).
3. Network connectivity issues reaching the repository, corporate proxy,
   firewall, or a genuinely down repository host.
4. An OCI-based chart repository (increasingly common) with an
   authentication or registry-login requirement that wasn't satisfied
   before the pull.

## Check it
```
helm repo list
helm repo update
helm search repo <chart> --versions
```
Confirm the repository is actually registered locally and its index is
current, then confirm the specific version you're requesting genuinely
exists in that repository's index.

## Fix it
- Add the repository if it's missing (`helm repo add <name> <url>`), then
  `helm repo update` to refresh the local index before retrying.
- Fix the chart version if it was a typo, or check the repository directly
  for what versions actually exist.
- For OCI registries, make sure you're authenticated (`helm registry
  login`) the same way you would need to be for `docker login` against a
  private registry.
- If it's a network/proxy issue, confirm the repository host is actually
  reachable from wherever Helm is running (a CI runner might have different
  network access than your local machine).

## Related
- PullAccessDenied
