package main

import (
	"fmt"
	"os"
	"sort"
)

// subcommands are the fixed, non-error-lookup words the CLI reserves. They
// show up in completion candidates alongside every entry's slug.
var subcommands = []string{"list", "random", "lint", "scan", "completion", "help", "version"}

// candidateNames returns every reserved subcommand plus every entry's slug,
// sorted. Shell completion scripts fetch this at completion time
// (`kube-why __candidates`) rather than a static list baked in at install
// time, so completions stay correct as packs grow without the completion
// script itself needing to change.
func candidateNames(entries []entry) []string {
	names := append([]string{}, subcommands...)
	for _, e := range entries {
		names = append(names, e.slug)
	}
	sort.Strings(names)
	return names
}

func printCandidates(entries []entry) {
	for _, n := range candidateNames(entries) {
		fmt.Println(n)
	}
}

func printCompletion(shell string) {
	switch shell {
	case "bash":
		fmt.Print(bashCompletion)
	case "zsh":
		fmt.Print(zshCompletion)
	case "fish":
		fmt.Print(fishCompletion)
	default:
		fmt.Fprintf(os.Stderr, "kube-why: unsupported shell %q, expected bash, zsh, or fish\n", shell)
		os.Exit(1)
	}
}

const bashCompletion = `# kube-why bash completion
# Install: kube-why completion bash > /etc/bash_completion.d/kube-why
# Or, for the current shell only: source <(kube-why completion bash)
_kube_why_complete() {
    local cur candidates
    cur="${COMP_WORDS[COMP_CWORD]}"
    candidates="$(kube-why __candidates 2>/dev/null)"
    COMPREPLY=($(compgen -W "${candidates}" -- "${cur}"))
}
complete -F _kube_why_complete kube-why
`

const zshCompletion = `#compdef kube-why
# kube-why zsh completion
# Install: kube-why completion zsh > "${fpath[1]}/_kube-why"
_kube_why() {
    local -a candidates
    candidates=(${(f)"$(kube-why __candidates 2>/dev/null)"})
    _describe 'kube-why' candidates
}
_kube_why
`

const fishCompletion = `# kube-why fish completion
# Install: kube-why completion fish > ~/.config/fish/completions/kube-why.fish
function __kube_why_candidates
    kube-why __candidates 2>/dev/null
end
complete -c kube-why -f -a '(__kube_why_candidates)'
`
