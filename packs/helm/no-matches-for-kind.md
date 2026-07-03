---
title: unable to recognize, no matches for kind
aliases: [unabletorecognizenomatchesforkind, nomatchesforkindinversion]
category: chart
---

# Unable to recognize "", no matches for kind

## What it means
The chart's manifest references a `kind` (a Custom Resource, usually) that
the target cluster's API server doesn't know about. This is the same root
cause as the Kubernetes pack's `The server doesn't have a resource type`,
but hit specifically through a Helm install/upgrade, which makes it easy to
mistake for a chart bug rather than a missing CRD.

## Common causes
1. The chart installs a Custom Resource, but the CRD that defines it wasn't
   installed first, either a separate CRD chart/step that was skipped, or
   ordering that assumed something already existed.
2. Installing a chart's CRDs and its main resources in the same Helm
   release, Helm applies CRDs from a chart's `crds/` directory before
   templates, but a CRD referenced from a *sub-chart* or *dependency*
   doesn't always get the same guarantee depending on chart structure.
3. Targeting the wrong cluster entirely, one where the operator providing
   this CRD was never installed at all.
4. An operator was uninstalled (removing its CRDs) while resources of that
   kind still exist or are still referenced by charts expected to run
   there.

## Check it
```
kubectl get crd | grep -i <kind>
kubectl api-resources | grep -i <kind>
```
Confirm whether the CRD is installed in this cluster at all, if it's
missing entirely, that's confirmed as the cause, not a chart templating
bug.

## Fix it
- Install the operator/CRD-providing chart first, before the chart that
  creates instances of its Custom Resource, check the chart's own
  documentation for the expected install order.
- If CRDs live in a separate chart specifically to control this ordering,
  don't skip that step, it exists for exactly this reason.
- Confirm you're targeting the intended cluster if the CRD was expected to
  already be present and genuinely isn't.

## Related
- NoResourceTypeFound
