package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// jsonSection is one "## Heading" block from an entry's body. Sections are
// kept in document order (What it means, Common causes, Check it, Fix it,
// Related) rather than a map, since consumers scripting against this almost
// always want the same order a human reading the entry would see.
type jsonSection struct {
	Heading string `json:"heading"`
	Body    string `json:"body"`
}

// jsonEntry is the full, single-entry shape returned by `kube-why <error>
// --json` and `kube-why random --json`. It mirrors entry, but the body is
// parsed into sections instead of shipped as raw markdown, since the whole
// point of --json is that something other than a human is going to read it.
type jsonEntry struct {
	Pack     string        `json:"pack"`
	Slug     string        `json:"slug"`
	Title    string        `json:"title"`
	Category string        `json:"category"`
	Aliases  []string      `json:"aliases"`
	Sections []jsonSection `json:"sections"`
}

// jsonListEntry is the lighter shape used by `kube-why list --json`, an
// index rather than full content, one per entry.
type jsonListEntry struct {
	Pack     string   `json:"pack"`
	Slug     string   `json:"slug"`
	Title    string   `json:"title"`
	Category string   `json:"category"`
	Aliases  []string `json:"aliases"`
}

// parseSections splits an entry's body into its "## Heading" blocks. The
// leading "# Title" line is dropped, it's already available as Title.
func parseSections(body string) []jsonSection {
	var sections []jsonSection
	var heading string
	var buf []string
	started := false

	flush := func() {
		if started {
			sections = append(sections, jsonSection{
				Heading: heading,
				Body:    strings.TrimSpace(strings.Join(buf, "\n")),
			})
		}
	}

	for _, line := range strings.Split(body, "\n") {
		switch {
		case strings.HasPrefix(line, "## "):
			flush()
			heading = strings.TrimPrefix(line, "## ")
			buf = nil
			started = true
		case strings.HasPrefix(line, "# "):
			continue
		default:
			buf = append(buf, line)
		}
	}
	flush()
	return sections
}

func toJSONEntry(e entry) jsonEntry {
	return jsonEntry{
		Pack:     e.pack,
		Slug:     e.slug,
		Title:    e.title,
		Category: e.category,
		Aliases:  e.aliases,
		Sections: parseSections(e.body),
	}
}

func printEntryJSON(e entry) {
	printJSON(toJSONEntry(e))
}

func toJSONEntries(entries []entry) []jsonEntry {
	out := make([]jsonEntry, len(entries))
	for i, e := range entries {
		out[i] = toJSONEntry(e)
	}
	return out
}

func printListJSON(entries []entry, packFilter string) {
	var out []jsonListEntry
	for _, e := range entries {
		if packFilter != "" && normalize(e.pack) != packFilter {
			continue
		}
		out = append(out, jsonListEntry{
			Pack:     e.pack,
			Slug:     e.slug,
			Title:    e.title,
			Category: e.category,
			Aliases:  e.aliases,
		})
	}
	if packFilter != "" && len(out) == 0 {
		printJSONError(fmt.Sprintf("no pack named %q", packFilter), packNames(entries))
		os.Exit(1)
	}
	printJSON(out)
}

// printJSONError is the --json equivalent of the plain "no entry found"
// text path, same information, machine-readable shape instead of prose.
func printJSONError(message string, suggestions []string) {
	printJSON(struct {
		Error       string   `json:"error"`
		Suggestions []string `json:"suggestions,omitempty"`
	}{Error: message, Suggestions: suggestions})
}

// printJSON encodes v as indented JSON to stdout. It disables Go's default
// HTML-escaping of <, >, and & (encoding/json's json.Marshal does this by
// default for safe embedding in <script> tags), since entry bodies are full
// of literal placeholders like <pod> and <deployment> in command examples,
// escaping them to unicode sequences would make every "check it" command
// unreadable to whatever's consuming --json output.
func printJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Fprintln(os.Stderr, "kube-why: failed to encode JSON:", err)
		os.Exit(1)
	}
}
