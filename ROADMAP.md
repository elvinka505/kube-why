# Roadmap

## What this project actually is

`kube-why` started as a Kubernetes error reference. It's becoming something
narrower in purpose but broader in scope: a terminal-native reference for
the cryptic errors across the whole cloud-native toolchain, not just
Kubernetes, browsable and pipeable the same way regardless of which tool
produced the error.

The mechanism for that is packs.

## How packs work

Content lives under `packs/<pack-name>/`, one directory per ecosystem. The
lookup engine (`main.go`) doesn't know anything about Kubernetes, Docker, or
Helm specifically, it just walks `packs/*/`, loads every `.md` file it finds,
and treats the parent directory name as the pack. Adding a new ecosystem is
adding a new directory, not changing code.

Every entry, regardless of pack, follows the same format documented in
[CONTRIBUTING.md](CONTRIBUTING.md): a frontmatter block (`title`, `aliases`,
`category`) and a body (what it means, common causes, how to check it, how
to fix it, related entries).

```
kube-why list docker     # see everything in one pack
kube-why list             # see everything, grouped by pack, then category
```

## Current packs

| Pack | Entries | Status |
|---|---|---|
| `kubernetes` | 48 | Established, the original scope |
| `docker` | 25 | Established |
| `helm` | 15 | Established |

## Planned packs, open for contribution

These don't exist yet. Each one is a real, self-contained unit of work,
create `packs/<name>/`, add a handful of entries following the existing
format, open a PR. No changes to the core engine required.

- **Terraform** — `Error: Cycle`, state lock errors, provider auth failures,
  drift-related apply failures. See [issue #1](https://github.com/Ayushmore1214/kube-why/issues/1).
- **CI/CD** (GitHub Actions specifically to start) — cryptic workflow syntax
  errors, permission errors on `GITHUB_TOKEN`, matrix build failures.

Open an issue before starting a new pack from scratch if you want a second
opinion on scope, or just open the PR directly if you're confident in it.
Either works.

## Not packs, general engine improvements

Separate from content, things that would improve the lookup engine itself
regardless of which pack someone's using:

- Smarter matching for packs whose real-world error text interpolates
  variable data mid-string (a container name, a resource ID), which can
  break simple substring detection on piped input. The Kubernetes pack
  mostly avoids this because `Reason` fields are fixed enum-style tokens,
  it's a real, live consideration for Docker and will be for Terraform too.
- `kube-why scan` — point it at a real cluster/host and get a live diagnosis
  of what's actually unhealthy right now, instead of pasting text in by
  hand.
