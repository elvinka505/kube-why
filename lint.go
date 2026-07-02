package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// runLint checks that a file contains syntactically valid YAML. It doesn't
// know anything about Kubernetes specifically yet, it just catches the kind
// of mistake (bad indentation, a stray tab, an unclosed quote, a duplicate
// key) that would otherwise show up later as a confusing kubectl error.
func runLint(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "kube-why: couldn't read %s: %v\n", path, err)
		os.Exit(1)
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	docIndex := 0
	hadError := false

	for {
		var doc interface{}
		err := decoder.Decode(&doc)
		if err == io.EOF {
			break
		}
		docIndex++

		if err != nil {
			hadError = true
			fmt.Printf("%sdocument %d: %s%s\n", colorRed, docIndex, err.Error(), colorReset)
			// A syntax error leaves the decoder unable to find the next
			// document boundary, re-calling Decode() on it just returns the
			// same error forever instead of reaching EOF. Stop here rather
			// than loop indefinitely.
			break
		}
		if doc == nil {
			// an empty document, e.g. a trailing "---" with nothing after it
			docIndex--
			continue
		}
		fmt.Printf("%sdocument %d: valid YAML%s\n", colorGreen, docIndex, colorReset)
	}

	if docIndex == 0 {
		fmt.Printf("kube-why: %s contains no YAML documents\n", path)
		os.Exit(1)
	}

	fmt.Println()
	if hadError {
		fmt.Printf("%s failed syntax check.\n", path)
		os.Exit(1)
	}
	fmt.Printf("%s: all %d document(s) are syntactically valid.\n", path, docIndex)
}
