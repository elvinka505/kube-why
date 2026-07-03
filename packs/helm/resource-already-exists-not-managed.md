---
title: rendered manifests contain a resource that already exists
aliases: [renderedmanifestscontainaresourcethatalreadyexists, invalidownershipmetadata]
category: install
---

# Rendered manifests contain a resource that already exists

## What it means
Helm found a resource in the cluster with the same name/kind/namespace as
one it's trying to create, but that existing resource isn't tagged as being
owned by this Helm release. Helm refuses to silently adopt or overwrite
resources it doesn't already manage, this is a safety check, not a bug.

## Common causes
1. The resource was created manually with `kubectl apply` before the Helm
   release existed, and now the chart tries to create the "same" resource.
2. Migrating an existing, manually-managed set of resources to be managed
   by Helm for the first time, without properly adopting them first.
3. Two different Helm releases (or a Helm release and a separate
   `kubectl`-applied manifest) both trying to own the same resource name.
4. A previous release using the same resource name was uninstalled with
   some cleanup step that didn't fully remove everything Helm expected to
   own.

## Check it
```
kubectl get <kind> <name> -n <namespace> -o jsonpath='{.metadata.labels}{"\n"}{.metadata.annotations}'
```
Check for the `app.kubernetes.io/managed-by: Helm` label and
`meta.helm.sh/release-name` / `meta.helm.sh/release-namespace` annotations,
their absence (or a different release name in them) confirms this resource
genuinely isn't tracked as belonging to the release trying to create it.

## Fix it
- To properly adopt an existing resource into Helm management, manually add
  the required labels/annotations Helm expects
  (`app.kubernetes.io/managed-by: Helm`, and the `meta.helm.sh/release-name`
  /`release-namespace` annotations matching the release that should own
  it), then retry.
- If the resource was created by a different release or process entirely
  and shouldn't be touched by this one, rename this chart's resource, or
  remove/migrate the conflicting one deliberately, don't just force an
  overwrite without understanding which one should actually own it.

## Related
- ReleaseNameAlreadyExists
