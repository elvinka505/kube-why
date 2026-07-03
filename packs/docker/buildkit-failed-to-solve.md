---
title: failed to solve (BuildKit)
aliases: [failedtosolve, buildkitfailedtosolve]
category: build
---

# failed to solve (BuildKit)

## What it means
"failed to solve" is BuildKit's generic wrapper error, it means some step
in the build graph failed, and the real cause is in the more specific
message that follows it, not in these three words themselves. Treat this as
a header, not the actual diagnosis.

## Common causes
1. A `RUN` instruction's command actually failed (non-zero exit), the
   underlying reason is whatever that command printed right before this
   wrapper message.
2. A `COPY`/`ADD` source path doesn't exist in the build context (see
   `COPY failed, file not found`).
3. A base image in `FROM` couldn't be pulled (see `pull access denied` or
   `manifest unknown`).
4. A build stage referenced by name in `COPY --from=` doesn't exist (a
   Dockerfile syntax/reference problem, see `dockerfile parse error`).

## Check it
```
docker build --progress=plain --no-cache .
```
`--progress=plain` shows full, unbuffered output for every build step
instead of BuildKit's default collapsed TTY view, this is usually the
fastest way to find the actual failing command and its real output, rather
than trying to interpret the summarized "failed to solve" wrapper alone.

## Fix it
- Read past "failed to solve" to the actual nested error, it always names a
  more specific problem, fix that specific problem, this message itself
  isn't independently actionable.
- If a `RUN` command failed, reproduce it manually inside a container based
  on the same intermediate state to debug interactively rather than
  guessing from build logs alone.
- Cross-reference the more specific underlying error against the other
  Docker pack entries here, it likely matches one of them exactly.

## Related
- DockerfileParseError
- CopyFailedFileNotFoundInBuildContext
- PullAccessDenied
