---
title: unable to build kubernetes objects from release manifest
aliases: [unabletobuildkubernetesobjectsfromreleasemanifest]
category: chart
---

# Unable to build kubernetes objects from release manifest

## What it means
Helm successfully rendered the chart's templates into YAML text, but that
YAML doesn't describe valid Kubernetes objects, a required field is
missing, a value has the wrong type for that API field, or the structure
doesn't match what the target `apiVersion`/`kind` expects. This is a step
after template rendering succeeds and before anything is actually sent to
the cluster.

## Common causes
1. A template renders a field with the wrong type (a template producing a
   quoted string where the API expects an integer, a common gotcha with Go
   templating and YAML's implicit typing).
2. A required field for that resource kind was omitted because a value
   controlling it was left empty or misconfigured.
3. The chart targets an `apiVersion` that's been removed in the Kubernetes
   version you're deploying to (see the `kubeVersion` note on chart
   compatibility).
4. Indentation or structural errors introduced by a template's conditional
   logic producing malformed YAML under certain value combinations that
   weren't tested.

## Check it
```
helm template <chart> -f values.yaml > /tmp/rendered.yaml
kube-why lint /tmp/rendered.yaml
kubectl apply --dry-run=server -f /tmp/rendered.yaml
```
Render the manifest to a file first and inspect it directly, the error from
Helm alone is often vague, seeing the actual generated YAML almost always
makes the malformed field or structure obvious immediately.

## Fix it
- Fix the specific field or template logic producing invalid output, once
  you can see the rendered YAML directly this is usually a quick, obvious
  fix.
- Watch for values that control types implicitly, an empty string versus an
  actual `null` versus an integer `0` can render very differently depending
  on how the template handles them.
- Add `kubectl apply --dry-run=server` as a pipeline step before real
  deploys, it catches this category of error without actually touching the
  cluster.

## Related
- ValuesDontMeetTheSpecificationsOfTheSchema
- KubeVersionIncompatible
