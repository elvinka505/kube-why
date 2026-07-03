---
title: dockerfile parse error
aliases: [dockerfileparseerror, failedtosolvewithfrontenddockerfile]
category: build
---

# Dockerfile parse error

## What it means
BuildKit (the modern Docker build engine) couldn't even parse the
Dockerfile into a valid build plan, this fails before any instruction
actually runs, it's a syntax problem in the Dockerfile itself, not a
problem with anything it's trying to do.

## Common causes
1. A malformed instruction, a `RUN`/`COPY`/`FROM` line with unbalanced
   quotes, an unescaped backslash at the end of a line that wasn't meant to
   continue, or invalid JSON-array syntax for exec-form instructions.
2. Mixing shell-form and exec-form syntax incorrectly in the same
   instruction.
3. A `FROM` referencing a build stage by a name that was never actually
   defined earlier in the file (typo in a multi-stage build's stage name).
4. An instruction that doesn't exist or is misspelled entirely.

## Check it
```
docker build --no-cache .
```
The parse error message names the exact line number and usually the
specific token it choked on, that's almost always enough to spot the typo
or malformed syntax directly without further digging.

## Fix it
- Fix the specific line the error names, check for unbalanced quotes,
  unintended line continuations (`\` at the end of a line you didn't mean
  to continue), and correct exec-form JSON array syntax
  (`["executable", "param1"]`, not a bare string).
- For multi-stage builds, double check every `FROM ... AS <name>` and every
  `COPY --from=<name>` reference the exact same stage name, a typo here
  produces a confusing "stage not found" variant of this same class of
  error.

## Related
- CopyFailedFileNotFoundInBuildContext
- BuildkitFailedToSolve
