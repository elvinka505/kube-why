---
title: invalid reference format
aliases: [invalidreferenceformat, invalidreferenceformatrepositoryname]
category: registry
---

# invalid reference format

## What it means
The image name/tag string you gave Docker isn't syntactically valid at all,
this fails before Docker ever tries to contact a registry. It's a parsing
error, not a network or auth error, similar in spirit to Kubernetes'
`InvalidImageName` but at the Docker CLI/daemon level directly.

## Common causes
1. Uppercase letters in the repository name, Docker image names must be
   lowercase (tags can have uppercase, the repository/name part cannot).
2. An unresolved shell variable or template placeholder left literally in
   the command (`docker run $IMAGE_NAME` where the variable was never set).
3. A stray space or quote character introduced when the image name was
   built up from a script or CI variable.
4. Too many colons, or a malformed digest reference
   (`name@sha256:...` with a broken hash format).

## Check it
```
echo "$IMAGE_NAME"
```
Print the exact string being passed to `docker run`/`docker pull` before it
gets used, this immediately reveals unresolved variables, stray whitespace,
or unintended casing that isn't obvious just reading the script.

## Fix it
- Lowercase the repository name, this is a hard requirement, not a style
  preference.
- Fix the variable substitution in the script/CI pipeline so the value is
  actually resolved before being passed to Docker.
- Validate image references before use in automation (a simple regex check
  or `docker inspect` dry run) rather than discovering it at deploy time.

## Related
- PullAccessDenied
