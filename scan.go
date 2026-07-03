package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// maxInspectPerScan caps how many individual `kubectl describe` or `docker
// logs`/`inspect` calls a single scan makes, so a badly broken cluster with
// thousands of crashing pods doesn't turn `kube-why scan` into thousands of
// subprocess launches.
const maxInspectPerScan = 20

// lookPath and runCommand are package vars rather than exec.LookPath and
// exec.Command called directly, so tests can substitute fakes without a
// real kubectl or docker installed.
var lookPath = func(name string) (string, error) { return exec.LookPath(name) }

var runCommand = func(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).CombinedOutput()
	return string(out), err
}

func commandAvailable(name string) bool {
	_, err := lookPath(name)
	return err == nil
}

// scanSource is one ecosystem's contribution to a scan. scanned is true
// only if the listing command (kubectl get pods / docker ps) actually
// succeeded, not merely that the binary was on PATH, a kubectl that exists
// but can't reach any cluster, or a docker CLI with no daemon behind it,
// must not be reported the same as "ran clean and found nothing", that
// would turn a broken scan into a false "all clear". Only kubectl
// get+describe and docker ps/logs/inspect are ever run, all read-only, scan
// never applies, deletes, or restarts anything.
type scanSource struct {
	name      string
	scanned   bool
	unhealthy int
	text      string
	skipped   string
}

type podRef struct{ namespace, name string }

// parseUnhealthyPods reads `kubectl get pods --all-namespaces --no-headers`
// output and returns every pod whose STATUS column isn't Running or
// Completed. Columns are NAMESPACE NAME READY STATUS RESTARTS AGE; only the
// first two and the fourth are ever read, so a RESTARTS column that itself
// contains a space (recent kubectl prints "3 (5m ago)") doesn't shift
// anything this function cares about.
func parseUnhealthyPods(output string) []podRef {
	var pods []podRef
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}
		namespace, name, status := fields[0], fields[1], fields[3]
		if status == "Running" || status == "Completed" {
			continue
		}
		pods = append(pods, podRef{namespace: namespace, name: name})
	}
	return pods
}

func gatherKubernetes() scanSource {
	s := scanSource{name: "kubernetes"}
	if !commandAvailable("kubectl") {
		s.skipped = "kubectl not found on PATH"
		return s
	}

	var buf strings.Builder

	podsOut, err := runCommand("kubectl", "get", "pods", "--all-namespaces", "--no-headers")
	if err != nil {
		s.skipped = "kubectl get pods failed: " + firstLine(podsOut+" "+err.Error())
		return s
	}
	s.scanned = true

	unhealthyPods := parseUnhealthyPods(podsOut)
	s.unhealthy += len(unhealthyPods)
	for i, p := range unhealthyPods {
		if i >= maxInspectPerScan {
			break
		}
		desc, _ := runCommand("kubectl", "describe", "pod", p.name, "-n", p.namespace)
		buf.WriteString(desc)
		buf.WriteString("\n")
	}

	if eventsOut, err := runCommand("kubectl", "get", "events", "--all-namespaces", "--field-selector", "type=Warning", "--no-headers"); err == nil {
		buf.WriteString(eventsOut)
		buf.WriteString("\n")
	}

	if nodesOut, err := runCommand("kubectl", "get", "nodes", "--no-headers"); err == nil {
		for _, line := range strings.Split(nodesOut, "\n") {
			fields := strings.Fields(line)
			if len(fields) >= 2 && fields[1] != "Ready" {
				s.unhealthy++
				buf.WriteString(line)
				buf.WriteString("\n")
			}
		}
	}

	s.text = buf.String()
	return s
}

func gatherDocker() scanSource {
	s := scanSource{name: "docker"}
	if !commandAvailable("docker") {
		s.skipped = "docker not found on PATH"
		return s
	}

	psOut, err := runCommand("docker", "ps", "-a", "--format", "{{.Names}}\t{{.Status}}")
	if err != nil {
		s.skipped = "docker ps failed: " + firstLine(psOut+" "+err.Error())
		return s
	}
	s.scanned = true

	var unhealthy []string
	scanner := bufio.NewScanner(strings.NewReader(psOut))
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), "\t", 2)
		if len(parts) != 2 {
			continue
		}
		name, status := parts[0], parts[1]
		if strings.HasPrefix(status, "Up") {
			continue
		}
		unhealthy = append(unhealthy, name)
	}
	s.unhealthy = len(unhealthy)

	var buf strings.Builder
	for i, name := range unhealthy {
		if i >= maxInspectPerScan {
			break
		}
		logs, _ := runCommand("docker", "logs", "--tail", "30", name)
		buf.WriteString(logs)
		buf.WriteString("\n")
		inspect, _ := runCommand("docker", "inspect", name,
			"--format", "{{.State.Status}} {{.State.Error}} exitcode={{.State.ExitCode}} oomkilled={{.State.OOMKilled}}")
		buf.WriteString(inspect)
		buf.WriteString("\n")
	}

	s.text = buf.String()
	return s
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

type scanSourceJSON struct {
	Name      string `json:"name"`
	Scanned   bool   `json:"scanned"`
	Unhealthy int    `json:"unhealthy"`
	Skipped   string `json:"skipped,omitempty"`
}

func toScanSourceJSON(sources []scanSource) []scanSourceJSON {
	out := make([]scanSourceJSON, len(sources))
	for i, s := range sources {
		out[i] = scanSourceJSON{Name: s.name, Scanned: s.scanned, Unhealthy: s.unhealthy, Skipped: s.skipped}
	}
	return out
}

// runScan is `kube-why scan`: gather live signal from whatever of
// kubectl/docker is on PATH and run it through the exact matcher piped
// input uses. It's "auto-detect from piped kubectl output" with kube-why
// running kubectl itself instead of waiting for you to paste it in.
//
// Exit codes are deliberately distinct so this is usable as a CI/alerting
// check: 0 nothing unhealthy, 1 found unhealthy resources, 2 couldn't scan
// at all (no kubectl or docker on PATH, the pack filter isn't scannable, or
// every requested source failed to actually run, wrong kubeconfig, daemon
// not reachable, etc.). 2 is deliberately distinct from 0, a scan that
// never ran is not the same thing as a clean bill of health.
func runScan(entries []entry, packFilter string, jsonMode bool) {
	var sources []scanSource
	switch packFilter {
	case "":
		sources = []scanSource{gatherKubernetes(), gatherDocker()}
	case "kubernetes":
		sources = []scanSource{gatherKubernetes()}
	case "docker":
		sources = []scanSource{gatherDocker()}
	default:
		msg := fmt.Sprintf("kube-why scan doesn't support pack %q, only kubernetes and docker have a live target to scan", packFilter)
		if jsonMode {
			printJSONError(msg, nil)
		} else {
			fmt.Fprintln(os.Stderr, "kube-why:", msg)
		}
		os.Exit(2)
	}

	anyScanned := false
	totalUnhealthy := 0
	var haystack strings.Builder
	for _, s := range sources {
		if s.scanned {
			anyScanned = true
			totalUnhealthy += s.unhealthy
			haystack.WriteString(s.text)
			haystack.WriteString("\n")
		}
	}

	if !anyScanned {
		msg := "nothing could be scanned, neither kubectl nor docker is usable right now"
		if len(sources) == 1 {
			msg = sources[0].skipped
		}
		if jsonMode {
			printJSONError(msg, nil)
		} else {
			fmt.Fprintln(os.Stderr, "kube-why:", msg)
		}
		os.Exit(2)
	}

	matches := matchEntries(entries, haystack.String())

	if jsonMode {
		printJSON(struct {
			Sources   []scanSourceJSON `json:"sources"`
			Unhealthy int              `json:"unhealthy"`
			Matches   []jsonEntry      `json:"matches"`
		}{
			Sources:   toScanSourceJSON(sources),
			Unhealthy: totalUnhealthy,
			Matches:   toJSONEntries(matches),
		})
		if totalUnhealthy > 0 {
			os.Exit(1)
		}
		return
	}

	for _, s := range sources {
		switch {
		case !s.scanned:
			fmt.Printf("%s: skipped, %s\n", s.name, s.skipped)
		case s.unhealthy == 0:
			fmt.Printf("%s: nothing unhealthy found\n", s.name)
		default:
			fmt.Printf("%s: %d unhealthy resource(s) found\n", s.name, s.unhealthy)
		}
	}

	if totalUnhealthy == 0 {
		fmt.Println("\nAll clear.")
		return
	}

	if len(matches) == 0 {
		fmt.Println("\nFound unhealthy resources, but didn't recognize a known error pattern in their output.")
		fmt.Println("Try describing/inspecting them directly, or run 'kube-why list' to see what's covered.")
		os.Exit(1)
	}

	fmt.Println()
	for i, e := range matches {
		if i > 0 {
			fmt.Println(strings.Repeat("-", 60))
		}
		printEntry(e)
	}
	os.Exit(1)
}
