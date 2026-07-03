---
title: found in Chart.yaml, but missing in charts/ directory
aliases: [foundinchartyamlbutmissinginchartsdirectory, chartdependencymissing]
category: chart
---

# Found in Chart.yaml, but missing in charts/ directory

## What it means
The chart declares a dependency in `Chart.yaml`'s `dependencies` section,
but the actual dependency chart archive isn't present in the `charts/`
subdirectory. Helm needs the dependency physically downloaded and vendored
locally before it can install or template anything that uses it.

## Common causes
1. `helm dependency update` (or `helm dependency build`) was never run after
   adding or changing a dependency in `Chart.yaml`.
2. The `charts/` directory (or its contents) was excluded from version
   control (a `.gitignore` entry) and never regenerated in this environment
   or CI pipeline.
3. A dependency's version constraint in `Chart.yaml` no longer matches
   what's actually vendored in `charts/`, after bumping the constraint
   without re-running the dependency update.
4. CI caches `charts/` from a previous run, and it's now stale relative to
   the current `Chart.yaml`.

## Check it
```
cat Chart.yaml
ls charts/
```
Compare what's declared in `Chart.yaml` against what's actually present in
`charts/`, the mismatch is almost always this simple to spot once you look
at both side by side.

## Fix it
- Run `helm dependency update` (fetches based on `Chart.lock` if present,
  or resolves fresh) to populate `charts/` correctly.
- Decide deliberately whether `charts/` should be committed to version
  control or regenerated in CI, and be consistent, half-committing it leads
  to exactly this kind of drift.
- In CI, add a `helm dependency update` step before every install/template
  step rather than assuming a cached or committed `charts/` is current.

## Related
- HelmUpgradeFailed
