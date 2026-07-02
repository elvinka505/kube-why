// kube-why looks up a Kubernetes error and prints what it means, why it
// usually happens, and how to fix it. No network calls, no dependencies,
// everything ships baked into the binary.
package main

import (
	"embed"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

//go:embed errors/*.md
var errorFiles embed.FS

// version is set at build time via -ldflags "-X main.version=...".
// GoReleaser sets this to the tag on every release build.
var version = "dev"

// Color codes are variables, not constants, so disableColor can blank them
// out once at startup for non-terminal output, NO_COLOR, or --no-color,
// without every print call needing to check a flag individually.
var (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorDim    = "\033[2m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
)

func disableColor() {
	colorReset, colorBold, colorCyan = "", "", ""
	colorGreen, colorDim, colorYellow, colorRed = "", "", "", ""
}

// colorShouldBeOff follows the same convention most terminal tools use:
// respect NO_COLOR (https://no-color.org) if set to anything, and disable
// automatically when output isn't going to a terminal (piped to a file,
// captured in CI, etc.), since ANSI codes just show up as literal garbage
// there.
func colorShouldBeOff(noColorFlag bool) bool {
	if noColorFlag {
		return true
	}
	if _, set := os.LookupEnv("NO_COLOR"); set {
		return true
	}
	return isPiped(os.Stdout)
}

type entry struct {
	slug     string
	title    string
	aliases  []string
	category string
	body     string
}

func main() {
	entries, err := loadEntries()
	if err != nil {
		fmt.Fprintln(os.Stderr, "kube-why: failed to load entries:", err)
		os.Exit(1)
	}

	args := os.Args[1:]

	noColorFlag := false
	filtered := args[:0]
	for _, a := range args {
		if a == "--no-color" {
			noColorFlag = true
			continue
		}
		filtered = append(filtered, a)
	}
	args = filtered

	if colorShouldBeOff(noColorFlag) {
		disableColor()
	}

	if len(args) == 0 {
		if isPiped(os.Stdin) {
			scanPipedInput(entries, os.Stdin)
			return
		}
		printUsage(entries)
		return
	}

	switch args[0] {
	case "-h", "--help", "help":
		printUsage(entries)
	case "-v", "--version", "version":
		fmt.Println("kube-why", version)
	case "list":
		printList(entries)
	case "random":
		printEntry(entries[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(entries))])
	case "lint":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "kube-why: lint requires a file, e.g. kube-why lint deployment.yaml")
			os.Exit(1)
		}
		runLint(args[1])
	default:
		term := normalize(strings.Join(args, " "))
		if e := find(entries, term); e != nil {
			printEntry(*e)
			return
		}
		fmt.Printf("kube-why: no entry found for %q\n\n", strings.Join(args, " "))
		suggestions := suggest(entries, term)
		if len(suggestions) > 0 {
			fmt.Println("Did you mean:")
			for _, s := range suggestions {
				fmt.Printf("  %s\n", s)
			}
			fmt.Println()
		}
		fmt.Println("Run 'kube-why list' to see everything covered so far.")
		fmt.Println("Don't see yours? Add it: https://github.com/Ayushmore1214/kube-why/blob/main/CONTRIBUTING.md")
		os.Exit(1)
	}
}

func loadEntries() ([]entry, error) {
	files, err := errorFiles.ReadDir("errors")
	if err != nil {
		return nil, err
	}

	var entries []entry
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
			continue
		}
		raw, err := errorFiles.ReadFile(path.Join("errors", f.Name()))
		if err != nil {
			return nil, err
		}
		e, err := parseEntry(strings.TrimSuffix(f.Name(), ".md"), string(raw))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", f.Name(), err)
		}
		entries = append(entries, e)
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].title < entries[j].title })
	return entries, nil
}

// parseEntry reads a minimal YAML-ish frontmatter block:
//
//	---
//	title: CrashLoopBackOff
//	aliases: [crashloop, crash-loop-backoff]
//	category: pod
//	---
//	<body>
func parseEntry(slug, raw string) (entry, error) {
	e := entry{slug: slug}
	lines := strings.Split(raw, "\n")

	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return e, fmt.Errorf("missing frontmatter")
	}

	i := 1
	for ; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "---" {
			i++
			break
		}
		key, val, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		switch key {
		case "title":
			e.title = val
		case "category":
			e.category = val
		case "aliases":
			val = strings.TrimPrefix(val, "[")
			val = strings.TrimSuffix(val, "]")
			for _, a := range strings.Split(val, ",") {
				a = strings.TrimSpace(a)
				if a != "" {
					e.aliases = append(e.aliases, a)
				}
			}
		}
	}

	if e.title == "" {
		return e, fmt.Errorf("missing title in frontmatter")
	}

	e.body = strings.TrimLeft(strings.Join(lines[i:], "\n"), "\n")
	return e, nil
}

func normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, "-", "")
	return s
}

func find(entries []entry, term string) *entry {
	for i, e := range entries {
		if normalize(e.slug) == term || normalize(e.title) == term {
			return &entries[i]
		}
		for _, a := range e.aliases {
			if normalize(a) == term {
				return &entries[i]
			}
		}
	}
	// fall back to substring match if nothing matched exactly
	for i, e := range entries {
		if strings.Contains(normalize(e.title), term) {
			return &entries[i]
		}
	}
	return nil
}

// isPiped reports whether stdin is connected to a pipe rather than a
// terminal, so `kubectl describe pod X | kube-why` is auto-detected without
// needing an explicit flag.
func isPiped(f *os.File) bool {
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

// scanPipedInput reads raw kubectl output (describe pod, get events, etc.)
// and matches any known error signature found in it. Reason strings like
// CrashLoopBackOff or OOMKilled appear verbatim, with no spaces or hyphens,
// in real kubectl output, so a plain substring check on a de-hyphenated,
// lowercased copy of the input is enough, no parsing of kubectl's output
// format required. Aliases under 5 characters are skipped here (but still
// work for direct lookups) since short tokens like "oom" risk matching
// unrelated words in arbitrary pasted text.
func scanPipedInput(entries []entry, r io.Reader) {
	data, err := io.ReadAll(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, "kube-why: failed to read stdin:", err)
		os.Exit(1)
	}
	haystack := strings.ReplaceAll(strings.ToLower(string(data)), "-", "")

	var matches []entry
	for _, e := range entries {
		candidates := append([]string{e.slug, e.title}, e.aliases...)
		for _, c := range candidates {
			needle := normalize(c)
			if len(needle) < 5 {
				continue
			}
			if strings.Contains(haystack, needle) {
				matches = append(matches, e)
				break
			}
		}
	}

	if len(matches) == 0 {
		fmt.Println("kube-why: didn't recognize an error pattern in that input.")
		fmt.Println("Run 'kube-why list' to see everything covered, or pass the error name directly.")
		os.Exit(1)
	}

	for i, e := range matches {
		if i > 0 {
			fmt.Println(strings.Repeat("-", 60))
		}
		printEntry(e)
	}
}

func suggest(entries []entry, term string) []string {
	var out []string
	for _, e := range entries {
		if strings.Contains(normalize(e.title), term[:min(3, len(term))]) {
			out = append(out, e.title)
		}
	}
	if len(out) > 5 {
		out = out[:5]
	}
	return out
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func printUsage(entries []entry) {
	fmt.Printf("%skube-why%s — look up what a Kubernetes error means and how to fix it\n\n", colorBold, colorReset)
	fmt.Println("Usage:")
	fmt.Println("  kube-why <error>          print what it means, why it happens, how to fix it")
	fmt.Println("  kube-why list             list every error currently covered")
	fmt.Println("  kube-why random           print a random one")
	fmt.Println("  kube-why lint <file>      check a YAML file's syntax before you apply it")
	fmt.Println("  <kubectl cmd> | kube-why  auto-detect the error from piped kubectl output")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  kube-why crashloopbackoff")
	fmt.Println("  kube-why oomkilled")
	fmt.Println("  kube-why \"image pull backoff\"")
	fmt.Println("  kube-why lint deployment.yaml")
	fmt.Println()
	fmt.Println("Add --no-color to disable colored output, or set NO_COLOR.")
	fmt.Printf("\n%d errors covered. Run 'kube-why list' to see them all.\n", len(entries))
}

func printList(entries []entry) {
	byCategory := map[string][]entry{}
	for _, e := range entries {
		byCategory[e.category] = append(byCategory[e.category], e)
	}
	var categories []string
	for c := range byCategory {
		categories = append(categories, c)
	}
	sort.Strings(categories)

	for _, c := range categories {
		fmt.Printf("%s%s%s%s\n", colorBold, colorCyan, strings.ToUpper(c), colorReset)
		for _, e := range byCategory[c] {
			fmt.Printf("  %-28s %s\n", e.title, colorDim+e.slug+colorReset)
		}
		fmt.Println()
	}
}

func printEntry(e entry) {
	for _, line := range strings.Split(e.body, "\n") {
		switch {
		case strings.HasPrefix(line, "# "):
			fmt.Printf("%s%s%s\n", colorBold, strings.TrimPrefix(line, "# "), colorReset)
		case strings.HasPrefix(line, "## "):
			fmt.Printf("%s%s%s%s\n", colorBold, colorCyan, strings.TrimPrefix(line, "## "), colorReset)
		case strings.HasPrefix(strings.TrimSpace(line), "```"):
			// skip fence markers (indented or not), color the code lines instead
			continue
		case strings.HasPrefix(line, "- "), strings.HasPrefix(line, "  - "):
			fmt.Printf("%s%s%s\n", colorYellow, line, colorReset)
		default:
			fmt.Println(colorizeInline(line))
		}
	}
}

// colorizeInline dims lines that look like shell commands (rough heuristic:
// starts with kubectl, or is indented as part of a fenced block).
func colorizeInline(line string) string {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "kubectl") || strings.HasPrefix(trimmed, "$") {
		return colorGreen + line + colorReset
	}
	return line
}
