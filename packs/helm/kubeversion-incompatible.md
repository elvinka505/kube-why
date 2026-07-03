---
title: chart requires kubeVersion which is incompatible with the cluster
aliases: [chartrequireskubeversionwhichisincompatible, kubeversionincompatible]
category: chart
---

# Chart requires kubeVersion which is incompatible

## What it means
A chart's `Chart.yaml` can declare a `kubeVersion` constraint, the range of
Kubernetes versions it's actually designed to work with. If the cluster
you're targeting falls outside that range, Helm refuses to install at all,
before rendering or validating anything else.

## Common causes
1. The cluster was upgraded (or downgraded) to a version outside the
   chart's declared supported range.
2. Using an older chart version against a much newer cluster, where the
   chart's constraint was set conservatively and never revisited.
3. Using a bleeding-edge chart version against an older, not-yet-upgraded
   cluster.
4. A typo or overly narrow constraint in the chart's own `Chart.yaml` that
   doesn't actually reflect real compatibility, a chart-authoring mistake
   rather than a genuine incompatibility.

## Check it
```
kubectl version
helm show chart <chart> | grep kubeVersion
```
Compare the cluster's actual server version against the constraint the
chart declares, this tells you immediately whether the cluster or the
chart choice needs to change.

## Fix it
- Use a chart version whose declared `kubeVersion` actually matches your
  cluster.
- If you maintain the chart yourself and believe the constraint is overly
  strict or simply wrong, widen it in `Chart.yaml` after actually verifying
  compatibility, don't just remove it to make the error go away.
- If the cluster genuinely needs to be upgraded to use a chart's newer
  features, that's a real infrastructure decision, not something to route
  around by ignoring the constraint.

## Related
- UnableToBuildKubernetesObjectsFromReleaseManifest
