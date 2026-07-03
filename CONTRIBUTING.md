# Contributing

The whole point of this project is that adding an entry should be easy. One
file, one PR.

## Adding a new error to an existing pack

1. Copy the template below into `packs/<pack>/<slug>.md`, e.g.
   `packs/docker/container-name-conflict.md`. The slug is the filename,
   lowercase, words separated by hyphens.
2. Fill it in based on something you've actually debugged, not a summary of
   someone else's blog post. If you had to figure it out the hard way, that's
   exactly the knowledge this project wants.
3. Run it locally to make sure it parses and reads well in a terminal:
   ```
   go run . <your-slug-or-alias>
   ```
4. Open a PR with just that one file (plus README.md's covered-count if you
   want to keep it in sync, not required).

## Adding a new pack

See [ROADMAP.md](ROADMAP.md) for which packs are already planned (Terraform,
CI/CD). Create `packs/<name>/` with at least one entry following the format
below, the lookup engine picks it up automatically, no code changes needed.
A pack doesn't need to be complete to merge, a handful of real, well-written
entries is a legitimate starting point others can build on.

### Template

```markdown
---
title: YourErrorName
aliases: [lowercase, alternate, spellings]
category: pod
---

# YourErrorName

## What it means
One or two sentences. What is the tool actually telling you here, in plain
language, not a restatement of the error string.

## Common causes
Ranked roughly by how often you'll actually see them, not an exhaustive list
of every theoretical cause. Three to five is usually right.

## Check it
The exact commands to confirm which cause you're dealing with. Real
commands, not "check the logs." Explain what to look for in the output, not
just what command to run.

## Fix it
The actual fix for each common cause, matched to the causes listed above. If
there's a quick mitigation and a real fix, say which is which.

## Related
Two or three other entries someone debugging this might also need.
```

`category` groups entries within a pack's `kube-why list` output, it's
pack-specific. The `kubernetes` pack currently uses: `pod`, `deployment`,
`scheduling`, `node`, `networking`, `ingress`, `job`, `hpa`, `namespace`,
`webhook`, `storage`, `statefulset`, `quota`, `podsecurity`, `apiserver`,
`container-runtime`, `rbac`. The `docker` pack uses: `daemon`, `container`,
`networking`, `storage`, `runtime`, `registry`, `build`. Add a new category
if none of those fit, it's cheap, just say why in the PR.

Aliases matter more than usual for piped-input detection (`<cmd> |
kube-why`): include at least one alias that's a fixed substring of the
real tool's actual error text. If that text interpolates something
variable in the middle (a name, an ID), pick a fixed chunk before or after
the variable part, not the whole phrase, otherwise it will never match
piped input even though direct lookup by name still works fine.

## What makes a good entry

- Causes are ranked by likelihood, not alphabetized or randomly ordered.
- The "check it" commands are things you'd actually run, and say what the
  output tells you, not just the command itself.
- No filler. If a sentence doesn't help someone fix the problem faster,
  cut it.
- Written like you're telling a coworker, not writing documentation.

## Fixing an existing entry

Same process, no new file needed. If something's outdated, wrong, or you have
a better fix, edit the file directly and explain what changed and why in the
PR description.

## Reporting an error you don't have time to write up

Open an issue with the error string and, if you know it, what actually caused
it for you. Someone else can turn it into a full entry.
