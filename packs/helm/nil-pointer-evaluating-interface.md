---
title: nil pointer evaluating interface, template error
aliases: [nilpointerevaluatinginterface, executingattemplatenilpointer]
category: chart
---

# Nil pointer evaluating interface (template error)

## What it means
A template referenced a value (`.Values.something.nested`) where an
intermediate part of that path doesn't exist, Go's templating engine can't
evaluate a field on something that's `nil`, and fails immediately rather
than silently producing an empty result. This is almost always a values
file that's missing a section the template assumes will be there.

## Common causes
1. A values override omits an entire nested section the template expects,
   accessing `.Values.database.host` when `.Values.database` itself was
   never set anywhere.
2. A typo in the values path within the template itself, referencing a key
   that doesn't exist under the correct parent.
3. An optional section in `values.yaml` that some deployments provide and
   others don't, without the template using safe accessors to handle its
   absence.
4. Upgrading a chart whose templates now reference new values fields your
   existing overrides don't include.

## Check it
```
helm template <chart> -f values.yaml --debug
```
The error message names the exact template file and line where evaluation
failed, and usually the specific field path being accessed, go straight to
that line rather than guessing.

## Fix it
- Add the missing section/field to your values override.
- If you maintain the chart, use safe accessors for genuinely optional
  values, `{{ if .Values.database }}` guards, or Helm's `default` function,
  so an absent section degrades gracefully instead of hard-failing:
  ```
  {{ .Values.database.host | default "localhost" }}
  ```
- After a chart upgrade, diff the new chart's default `values.yaml` against
  your override to catch newly-referenced fields before they fail at
  template time.

## Related
- ValuesDontMeetTheSpecificationsOfTheSchema
