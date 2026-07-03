---
title: values don't meet the specifications of the schema
aliases: [valuesdontmeetthespecificationsoftheschema, valuesschemajsonvalidationfailed]
category: chart
---

# Values don't meet the specifications of the schema

## What it means
The chart ships a `values.schema.json` that defines what shape and type its
`values.yaml` (and any overrides you pass) must have. Helm validates your
merged values against that schema before templating anything, and rejects
the install/upgrade immediately if they don't match, this happens before
any Kubernetes object is even rendered.

## Common causes
1. A required field defined in the schema was never set in your values
   override.
2. A value's type doesn't match what the schema expects (a string where the
   schema requires a number, or vice versa, common when overriding via
   `--set` which doesn't always infer types the way you expect).
3. The chart was upgraded to a newer version with new required fields, and
   an existing values file that worked on the previous chart version
   doesn't satisfy the new schema.
4. A typo in a field name, so it doesn't match any property the schema
   actually defines, depending on the schema's strictness this can fail
   validation instead of just being silently ignored.

## Check it
```
helm lint <chart> -f values.yaml
helm template <chart> -f values.yaml --debug
```
The validation error itself usually names the exact field and constraint it
violated, `helm lint` surfaces this without needing to attempt a real
install first.

## Fix it
- Add the missing required field, or fix the type mismatch the error names
  directly.
- After upgrading to a newer chart version, diff the new chart's
  `values.schema.json` (or its default `values.yaml`) against your existing
  overrides to catch newly required fields before they fail at deploy time.
- Be deliberate with `--set` for typed fields, use `--set-string` or a
  proper values file when the schema expects a specific type `--set`'s
  inference might get wrong.

## Related
- UnableToBuildKubernetesObjectsFromReleaseManifest
