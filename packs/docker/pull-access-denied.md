---
title: Pull access denied / manifest unknown
aliases: [pullaccessdenied, manifestunknown, repositorydoesnotexistormaynotpermission]
category: registry
---

# Pull access denied / manifest unknown

## What it means
Docker couldn't pull the image, either because it genuinely doesn't exist at
that name/tag, or because it exists but you're not authenticated to access
it. Docker deliberately gives the same vague error for both "doesn't exist"
and "you don't have access," so it doesn't leak whether a private image
exists to someone who isn't authorized to see it.

## Common causes
1. A typo in the image name or tag.
2. The image is private, and you're not logged in to the registry at all
   (`docker login` was never run, or the session expired).
3. You're logged into the wrong registry or the wrong account for a private
   image under a different organization/namespace.
4. The image genuinely doesn't exist, was deleted, or the tag was never
   pushed.

## Check it
```
docker pull <image>:<tag>
docker login <registry>
```
Try logging in explicitly and re-pulling. If it succeeds after login, it was
an auth issue, if it still fails identically, the image/tag likely doesn't
exist as named.

## Fix it
- Typo: fix the image name or tag.
- Private image, not logged in: `docker login` to the correct registry with
  an account that actually has access.
- Wrong account/org: log out and back in with the correct credentials,
  double check which organization's private images you're supposed to have
  access to.
- Genuinely missing: confirm with whoever publishes the image that the tag
  was actually pushed, tags can be deleted or never published for a given
  build.

## Related
- ExecFormatError
- DockerCantReachRegistry