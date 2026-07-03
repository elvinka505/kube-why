// kube-why looks up a cloud-native error and prints what it means, why it
// usually happens, and how to fix it. No network calls, no dependencies,
// everything ships baked into the binary.
//
// Content is organized into "packs", one directory per ecosystem
// (packs/kubernetes, packs/docker, packs/helm, ...). The lookup engine
// itself doesn't know anything about Kubernetes specifically, a pack is
// just a directory of entries, so adding a new ecosystem is a new
// directory, not a code change.
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

//go:embed packs
var packFiles embed.FS

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
	pack     string
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
	jsonFlag := false
	filtered := args[:0]
	for _, a := range args {
		switch a {
		case "--no-color":
			noColorFlag = true
		case "--json":
			jsonFlag = true
		default:
			filtered = append(filtered, a)
		}
	}
	args = filtered

	if colorShouldBeOff(noColorFlag) || jsonFlag {
		disableColor()
	}

	if len(args) == 0 {
		if isPiped(os.Stdin) {
			scanPipedInput(entries, os.Stdin, jsonFlag)
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
		packFilter := ""
		if len(args) > 1 {
			packFilter = normalize(args[1])
		}
		if jsonFlag {
			printListJSON(entries, packFilter)
		} else {
			printList(entries, packFilter)
		}
	case "random":
		e := entries[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(entries))]
		if jsonFlag {
			printEntryJSON(e)
		} else {
			printEntry(e)
		}
	case "lint":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "kube-why: lint requires a file, e.g. kube-why lint deployment.yaml")
			os.Exit(1)
		}
		runLint(args[1], jsonFlag)
	case "scan":
		packFilter := ""
		if len(args) > 1 {
			packFilter = normalize(args[1])
		}
		runScan(entries, packFilter, jsonFlag)
	case "completion":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "kube-why: completion requires a shell, e.g. kube-why completion bash")
			os.Exit(1)
		}
		printCompletion(args[1])
	case "__candidates":
		printCandidates(entries)
	default:
		term := normalize(strings.Join(args, " "))
		if e := find(entries, term); e != nil {
			if jsonFlag {
				printEntryJSON(*e)
			} else {
				printEntry(*e)
			}
			return
		}
		suggestions := suggest(entries, term)
		if jsonFlag {
			printJSONError(fmt.Sprintf("no entry found for %q", strings.Join(args, " ")), suggestions)
			os.Exit(1)
		}
		fmt.Printf("kube-why: no entry found for %q\n\n", strings.Join(args, " "))
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
	packs, err := packFiles.ReadDir("packs")
	if err != nil {
		return nil, err
	}

	var entries []entry
	for _, p := range packs {
		if !p.IsDir() {
			continue
		}
		pack := p.Name()

		files, err := packFiles.ReadDir(path.Join("packs", pack))
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
				continue
			}
			raw, err := packFiles.ReadFile(path.Join("packs", pack, f.Name()))
			if err != nil {
				return nil, err
			}
			e, err := parseEntry(strings.TrimSuffix(f.Name(), ".md"), string(raw))
			if err != nil {
				return nil, fmt.Errorf("%s/%s: %w", pack, f.Name(), err)
			}
			e.pack = pack
			entries = append(entries, e)
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].pack != entries[j].pack {
			return entries[i].pack < entries[j].pack
		}
		return entries[i].title < entries[j].title
	})
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

// matchEntries scans arbitrary text (piped kubectl/docker output, or
// anything `kube-why scan` gathers by running commands itself) for any
// known error signature and returns every entry whose slug, title, or an
// alias shows up in it. Reason strings like CrashLoopBackOff or OOMKilled
// appear verbatim, with no spaces or hyphens, in real kubectl output, so a
// plain substring check on a de-hyphenated, lowercased copy of the input is
// enough, no parsing of the source command's output format required.
// Aliases under 5 characters are skipped here (but still work for direct
// lookups) since short tokens like "oom" risk matching unrelated words in
// arbitrary text.
func matchEntries(entries []entry, text string) []entry {
	haystack := normalize(text)

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
	return matches
}

// scanPipedInput reads raw kubectl/docker output from stdin and reports any
// known error signature found in it.
func scanPipedInput(entries []entry, r io.Reader, jsonMode bool) {
	data, err := io.ReadAll(r)
	if err != nil {
		if jsonMode {
			printJSONError(fmt.Sprintf("failed to read stdin: %v", err), nil)
		} else {
			fmt.Fprintln(os.Stderr, "kube-why: failed to read stdin:", err)
		}
		os.Exit(1)
	}

	matches := matchEntries(entries, string(data))

	if len(matches) == 0 {
		if jsonMode {
			printJSONError("didn't recognize an error pattern in that input", nil)
		} else {
			fmt.Println("kube-why: didn't recognize an error pattern in that input.")
			fmt.Println("Run 'kube-why list' to see everything covered, or pass the error name directly.")
		}
		os.Exit(1)
	}

	if jsonMode {
		printJSON(toJSONEntries(matches))
		return
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

func packNames(entries []entry) []string {
	seen := map[string]bool{}
	var packs []string
	for _, e := range entries {
		if !seen[e.pack] {
			seen[e.pack] = true
			packs = append(packs, e.pack)
		}
	}
	sort.Strings(packs)
	return packs
}

func printUsage(entries []entry) {
	fmt.Printf("%skube-why%s — look up what a cloud-native error means and how to fix it\n\n", colorBold, colorReset)
	fmt.Println("Usage:")
	fmt.Println("  kube-why <error>          print what it means, why it happens, how to fix it")
	fmt.Println("  kube-why list [pack]      list everything covered, optionally for one pack")
	fmt.Println("  kube-why random           print a random one")
	fmt.Println("  kube-why lint <file>      check a YAML file's syntax before you apply it")
	fmt.Println("  kube-why scan [pack]      run kubectl/docker yourself and diagnose what's unhealthy now")
	fmt.Println("  kube-why completion <shell>  print a completion script (bash, zsh, fish)")
	fmt.Println("  <cmd> | kube-why          auto-detect the error from piped command output")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  kube-why crashloopbackoff")
	fmt.Println("  kube-why oomkilled")
	fmt.Println("  kube-why \"image pull backoff\"")
	fmt.Println("  kube-why lint deployment.yaml")
	fmt.Println("  kube-why list docker")
	fmt.Println("  kube-why scan")
	fmt.Println("  kube-why oomkilled --json")
	fmt.Println("  source <(kube-why completion zsh)")
	fmt.Println()
	fmt.Println("Add --no-color to disable colored output, or set NO_COLOR.")
	fmt.Println("Add --json to any lookup, list, random, lint, or scan for machine-readable output.")
	fmt.Println("scan only ever runs read-only kubectl get/describe and docker ps/logs/inspect,")
	fmt.Println("it never applies, deletes, or restarts anything. Exit codes: 0 all clear,")
	fmt.Println("1 found something, 2 couldn't scan (no kubectl/docker on PATH).")
	fmt.Printf("\n%d entries across %d packs: %s. Run 'kube-why list' to see them all.\n",
		len(entries), len(packNames(entries)), strings.Join(packNames(entries), ", "))
}

func printList(entries []entry, packFilter string) {
	byPack := map[string][]entry{}
	for _, e := range entries {
		if packFilter != "" && normalize(e.pack) != packFilter {
			continue
		}
		byPack[e.pack] = append(byPack[e.pack], e)
	}

	if len(byPack) == 0 {
		fmt.Printf("kube-why: no pack named %q. Packs: %s\n", packFilter, strings.Join(packNames(entries), ", "))
		os.Exit(1)
	}

	var packs []string
	for p := range byPack {
		packs = append(packs, p)
	}
	sort.Strings(packs)

	for _, pack := range packs {
		fmt.Printf("%s%s%s\n", colorBold, strings.ToUpper(pack), colorReset)

		byCategory := map[string][]entry{}
		for _, e := range byPack[pack] {
			byCategory[e.category] = append(byCategory[e.category], e)
		}
		var categories []string
		for c := range byCategory {
			categories = append(categories, c)
		}
		sort.Strings(categories)

		for _, c := range categories {
			fmt.Printf("  %s%s%s%s\n", colorCyan, strings.ToUpper(c), colorReset, "")
			width := 0
			for _, e := range byCategory[c] {
				if len(e.title) > width {
					width = len(e.title)
				}
			}
			for _, e := range byCategory[c] {
				fmt.Printf("    %-*s  %s\n", width, e.title, colorDim+e.slug+colorReset)
			}
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

// colorizeInline highlights lines that look like real commands to run,
// recognized by the tool name a pack's examples actually use (kubectl for
// the kubernetes pack, docker for the docker pack, and so on).
func colorizeInline(line string) string {
	trimmed := strings.TrimSpace(line)
	commandPrefixes := []string{"kubectl", "docker", "helm", "$"}
	for _, p := range commandPrefixes {
		if strings.HasPrefix(trimmed, p) {
			return colorGreen + line + colorReset
		}
	}
	return line
}
