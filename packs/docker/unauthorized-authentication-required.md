---
title: unauthorized, authentication required
aliases: [unauthorizedauthenticationrequired, dockerloginrequired]
category: registry
---

# unauthorized, authentication required

## What it means
The registry is telling Docker it needs credentials to even evaluate the
request, this is subtly different from `pull access denied`, which usually
means you're authenticated but not authorized for that specific image.
`unauthorized` here typically means no valid credentials were presented at
all.

## Common causes
1. Never having run `docker login` for this registry in the current
   environment (a fresh CI runner, a new machine).
2. A previously valid login session expired, registries commonly issue
   short-lived tokens that need periodic re-authentication.
3. Pushing (not just pulling) to a registry without having authenticated for
   write access specifically, some registries separate read and write auth.
4. A credential helper (`docker-credential-*`) misconfigured or not
   installed, so Docker can't actually retrieve stored credentials even
   though they exist somewhere.

## Check it
```
docker login <registry>
cat ~/.docker/config.json
```
Check whether `~/.docker/config.json` actually has an entry for the
registry you're targeting, an empty or missing entry confirms you were
never actually authenticated for it, regardless of what you might expect
from a previous session.

## Fix it
- Run `docker login <registry>` explicitly with valid credentials before
  pulling or pushing.
- In CI, ensure the login step actually runs before the push/pull step, and
  that credentials are injected fresh each run rather than assumed to
  persist between jobs.
- If using a credential helper, confirm it's actually installed and
  functioning (`docker-credential-<helper> list`), a broken helper silently
  produces this same error.

## Related
- PullAccessDenied
- PullRateLimitExceeded
