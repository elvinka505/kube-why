package main

import (
	"errors"
	"strings"
	"testing"
)

// withFakeCommands swaps lookPath and runCommand for the duration of fn, so
// gatherKubernetes/gatherDocker can be tested without a real kubectl or
// docker installed. failing lists command prefixes that should return an
// error, so tests can simulate "the binary is on PATH but the command
// itself fails" (wrong kubeconfig, daemon unreachable, etc.), distinct from
// "not installed at all".
func withFakeCommands(t *testing.T, available map[string]bool, outputs map[string]string, failing []string, fn func()) {
	t.Helper()
	origLookPath, origRunCommand := lookPath, runCommand
	defer func() { lookPath, runCommand = origLookPath, origRunCommand }()

	lookPath = func(name string) (string, error) {
		if available[name] {
			return "/usr/bin/" + name, nil
		}
		return "", errors.New("not found")
	}
	runCommand = func(name string, args ...string) (string, error) {
		key := name + " " + strings.Join(args, " ")
		for _, f := range failing {
			if strings.HasPrefix(key, f) {
				return "", errors.New("command failed")
			}
		}
		for k, v := range outputs {
			if strings.HasPrefix(key, k) {
				return v, nil
			}
		}
		return "", nil
	}
	fn()
}

func TestParseUnhealthyPodsSkipsRunningAndCompleted(t *testing.T) {
	output := strings.Join([]string{
		"default     web-1        1/1   Running     0     3d",
		"default     migrate-job  0/1   Completed   0     1h",
		"default     web-2        0/1   CrashLoopBackOff   5     10m",
		"kube-system dns-1        0/1   Pending     0     2m",
	}, "\n")

	pods := parseUnhealthyPods(output)
	if len(pods) != 2 {
		t.Fatalf("parseUnhealthyPods() returned %d pods, want 2, got %+v", len(pods), pods)
	}
	if pods[0] != (podRef{namespace: "default", name: "web-2"}) {
		t.Errorf("pods[0] = %+v", pods[0])
	}
	if pods[1] != (podRef{namespace: "kube-system", name: "dns-1"}) {
		t.Errorf("pods[1] = %+v", pods[1])
	}
}

func TestGatherKubernetesSkippedWhenKubectlMissing(t *testing.T) {
	withFakeCommands(t, map[string]bool{}, nil, nil, func() {
		s := gatherKubernetes()
		if s.scanned {
			t.Fatal("expected scanned=false when kubectl isn't on PATH")
		}
		if s.skipped == "" {
			t.Fatal("expected a skipped reason")
		}
	})
}

// TestGatherKubernetesSkippedWhenGetPodsFails is a regression test: an
// earlier version reported this as scanned=true with unhealthy=0, which
// runScan then printed as "nothing unhealthy found" / exit 0, a false "all
// clear" for a kubectl that can't actually reach any cluster (wrong
// kubeconfig, VPN down, etc.).
func TestGatherKubernetesSkippedWhenGetPodsFails(t *testing.T) {
	withFakeCommands(t,
		map[string]bool{"kubectl": true},
		nil,
		[]string{"kubectl get pods"},
		func() {
			s := gatherKubernetes()
			if s.scanned {
				t.Fatal("expected scanned=false when `kubectl get pods` itself fails")
			}
			if s.skipped == "" {
				t.Fatal("expected a skipped reason explaining why")
			}
		},
	)
}

func TestGatherKubernetesFindsUnhealthyPods(t *testing.T) {
	withFakeCommands(t,
		map[string]bool{"kubectl": true},
		map[string]string{
			"kubectl get pods":           "default web-1 0/1 CrashLoopBackOff 5 10m",
			"kubectl describe pod web-1": "Back-off restarting failed container\nReason: CrashLoopBackOff",
			"kubectl get events":         "",
			"kubectl get nodes":          "node-1 Ready",
		},
		nil,
		func() {
			s := gatherKubernetes()
			if !s.scanned {
				t.Fatal("expected scanned=true when kubectl get pods succeeds")
			}
			if s.unhealthy != 1 {
				t.Fatalf("unhealthy = %d, want 1", s.unhealthy)
			}
			if !strings.Contains(s.text, "CrashLoopBackOff") {
				t.Fatalf("expected gathered text to contain the describe output, got %q", s.text)
			}
		},
	)
}

func TestGatherDockerSkippedWhenDockerMissing(t *testing.T) {
	withFakeCommands(t, map[string]bool{}, nil, nil, func() {
		s := gatherDocker()
		if s.scanned {
			t.Fatal("expected scanned=false when docker isn't on PATH")
		}
	})
}

// TestGatherDockerSkippedWhenPsFails is the docker-side counterpart to
// TestGatherKubernetesSkippedWhenGetPodsFails: a docker CLI that's
// installed but can't reach a daemon must not be reported the same as "ran
// clean, found nothing".
func TestGatherDockerSkippedWhenPsFails(t *testing.T) {
	withFakeCommands(t,
		map[string]bool{"docker": true},
		nil,
		[]string{"docker ps"},
		func() {
			s := gatherDocker()
			if s.scanned {
				t.Fatal("expected scanned=false when `docker ps` itself fails")
			}
			if s.skipped == "" {
				t.Fatal("expected a skipped reason explaining why")
			}
		},
	)
}

func TestGatherDockerFindsUnhealthyContainers(t *testing.T) {
	withFakeCommands(t,
		map[string]bool{"docker": true},
		map[string]string{
			"docker ps":      "myapp\tExited (137) 2 minutes ago",
			"docker logs":    "standard_init_linux.go:228: exec user process caused: exec format error",
			"docker inspect": "exited  exitcode=1 oomkilled=false",
		},
		nil,
		func() {
			s := gatherDocker()
			if !s.scanned {
				t.Fatal("expected scanned=true when docker ps succeeds")
			}
			if s.unhealthy != 1 {
				t.Fatalf("unhealthy = %d, want 1", s.unhealthy)
			}
			if !strings.Contains(s.text, "exec format error") {
				t.Fatalf("expected gathered text to contain the container's logs, got %q", s.text)
			}
		},
	)
}

func TestGatherDockerIgnoresRunningContainers(t *testing.T) {
	withFakeCommands(t,
		map[string]bool{"docker": true},
		map[string]string{
			"docker ps": "myapp\tUp 3 hours",
		},
		nil,
		func() {
			s := gatherDocker()
			if s.unhealthy != 0 {
				t.Fatalf("unhealthy = %d, want 0 for a container that's Up", s.unhealthy)
			}
		},
	)
}

func TestMatchEntriesFindsKnownSignatureInArbitraryText(t *testing.T) {
	entries, err := loadEntries()
	if err != nil {
		t.Fatalf("loadEntries() failed: %v", err)
	}
	matches := matchEntries(entries, "pod is stuck with reason CrashLoopBackOff after 5 restarts")
	found := false
	for _, m := range matches {
		if m.slug == "crashloopbackoff" {
			found = true
		}
	}
	if !found {
		t.Errorf("matchEntries() didn't find crashloopbackoff in text containing it, got %+v", matches)
	}
}
