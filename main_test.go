package main

import "testing"

func TestLoadEntriesNoErrors(t *testing.T) {
	entries, err := loadEntries()
	if err != nil {
		t.Fatalf("loadEntries() failed: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("loadEntries() returned zero entries")
	}
}

func TestEveryEntryResolvesBySlug(t *testing.T) {
	entries, err := loadEntries()
	if err != nil {
		t.Fatalf("loadEntries() failed: %v", err)
	}
	for _, e := range entries {
		if find(entries, normalize(e.slug)) == nil {
			t.Errorf("entry %q does not resolve by its own slug", e.slug)
		}
	}
}

func TestNoDuplicateAliasesAcrossEntries(t *testing.T) {
	entries, err := loadEntries()
	if err != nil {
		t.Fatalf("loadEntries() failed: %v", err)
	}

	owner := map[string]string{}
	for _, e := range entries {
		candidates := append([]string{e.slug, e.title}, e.aliases...)
		for _, c := range candidates {
			key := normalize(c)
			if key == "" {
				continue
			}
			id := e.pack + "/" + e.slug
			if prev, exists := owner[key]; exists && prev != id {
				t.Errorf("alias %q is claimed by both %q and %q, they'll silently shadow each other", key, prev, id)
			}
			owner[key] = id
		}
	}
}

func TestEveryEntryHasAPack(t *testing.T) {
	entries, err := loadEntries()
	if err != nil {
		t.Fatalf("loadEntries() failed: %v", err)
	}
	for _, e := range entries {
		if e.pack == "" {
			t.Errorf("entry %q has no pack, it won't be reachable via 'kube-why list <pack>'", e.slug)
		}
	}
}

func TestAtLeastTwoPacksExist(t *testing.T) {
	entries, err := loadEntries()
	if err != nil {
		t.Fatalf("loadEntries() failed: %v", err)
	}
	packs := packNames(entries)
	if len(packs) < 2 {
		t.Fatalf("expected at least 2 packs (kubernetes + one more), got %v", packs)
	}
}

func TestEveryEntryHasRequiredFields(t *testing.T) {
	entries, err := loadEntries()
	if err != nil {
		t.Fatalf("loadEntries() failed: %v", err)
	}
	for _, e := range entries {
		if e.title == "" {
			t.Errorf("entry %q has no title", e.slug)
		}
		if e.category == "" {
			t.Errorf("entry %q has no category", e.slug)
		}
		if len(e.aliases) == 0 {
			t.Errorf("entry %q has no aliases, it can only be found by its exact slug", e.slug)
		}
		if e.body == "" {
			t.Errorf("entry %q has an empty body", e.slug)
		}
	}
}

func TestNormalize(t *testing.T) {
	cases := map[string]string{
		"CrashLoopBackOff":     "crashloopbackoff",
		"crash-loop-backoff":   "crashloopbackoff",
		" image pull backoff ": "imagepullbackoff",
		"OOM_Killed":           "oomkilled",
	}
	for in, want := range cases {
		if got := normalize(in); got != want {
			t.Errorf("normalize(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFindMatchesByAliasCaseAndSpacing(t *testing.T) {
	entries, err := loadEntries()
	if err != nil {
		t.Fatalf("loadEntries() failed: %v", err)
	}

	e := find(entries, normalize("OOMKilled"))
	if e == nil {
		t.Fatal("expected a match for OOMKilled")
	}
	if e.slug != "oomkilled" {
		t.Errorf("normalize/find resolved to slug %q, want %q", e.slug, "oomkilled")
	}

	if find(entries, normalize("not-a-real-error-xyz")) != nil {
		t.Error("expected no match for a nonsense term")
	}
}

func TestParseEntryRejectsMissingFrontmatter(t *testing.T) {
	_, err := parseEntry("broken", "# Just a heading\nno frontmatter here")
	if err == nil {
		t.Fatal("expected an error for a file with no frontmatter block")
	}
}

func TestParseEntryRejectsMissingTitle(t *testing.T) {
	raw := "---\naliases: [foo]\ncategory: pod\n---\n# Body\n"
	_, err := parseEntry("broken", raw)
	if err == nil {
		t.Fatal("expected an error for a frontmatter block with no title")
	}
}

func TestParseEntryParsesAliasesList(t *testing.T) {
	raw := "---\ntitle: Example\naliases: [one, two, three]\ncategory: pod\n---\nbody text\n"
	e, err := parseEntry("example", raw)
	if err != nil {
		t.Fatalf("parseEntry failed: %v", err)
	}
	want := []string{"one", "two", "three"}
	if len(e.aliases) != len(want) {
		t.Fatalf("got %d aliases, want %d", len(e.aliases), len(want))
	}
	for i, a := range want {
		if e.aliases[i] != a {
			t.Errorf("alias[%d] = %q, want %q", i, e.aliases[i], a)
		}
	}
}
