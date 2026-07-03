---
title: COPY failed, file not found in build context
aliases: [copyfailednosuchfileordirectory, filenotfoundinbuildcontext, copyfailedforbiddenpath, copyfailedfilenotfoundinbuildcontext]
category: build
---

# COPY failed, file not found in build context

## What it means
A `COPY` or `ADD` instruction in a Dockerfile references a file or path that
Docker can't find inside the build context, the set of files actually sent
to the Docker daemon when the build starts, which isn't necessarily every
file in your project directory.

## Common causes
1. The file is excluded by `.dockerignore`, so it's never part of the build
   context even though it exists in your project.
2. The `COPY` source path is relative to the build context root, not to the
   Dockerfile's own location, easy to get wrong when the Dockerfile lives in
   a subdirectory.
3. The build is run with the wrong context path (`docker build -f
   path/to/Dockerfile .` still uses `.` as the context, not the Dockerfile's
   directory).
4. The file genuinely wasn't created yet at build time, expected to exist
   from a step that runs outside the Docker build (a separate compile step
   that was supposed to run first).

## Check it
```
cat .dockerignore
docker build --no-cache -f <dockerfile> <context-path> 2>&1 | head -30
```
Check `.dockerignore` first, it's the most common silent cause. Confirm
which directory is actually being used as the build context, it's the last
argument to `docker build`, not necessarily where the Dockerfile lives.

## Fix it
- If `.dockerignore` is excluding a file you actually need in the image,
  remove or narrow that ignore rule.
- Fix the `COPY` path to be relative to the actual build context root, not
  the Dockerfile's directory.
- If a prior build step is supposed to produce the file, confirm that step
  actually ran and produced it before the `docker build` step in your
  pipeline, rather than assuming ordering.

## Related
- ExecFormatError
