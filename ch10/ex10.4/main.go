package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "requires package name")
		os.Exit(1)
	}

	packageName := os.Args[1]
	out, err := exec.Command("go", "list", "-f", "{{ .ImportPath }}", packageName).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get import path of %s: %v\n", packageName, err)
		os.Exit(1)
	}

	target := strings.TrimSpace(string(out))
	//fmt.Printf("target=%q\n", target)

	depMap := make(map[string]struct{})

	out, err = exec.Command("go", "list", "-f", "{{ .ImportPath }} {{join .Deps \" \"}}", "...").Output()
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		packageDeps := strings.Split(line, " ")
		p := packageDeps[0]
		deps := packageDeps[1:]
		for _, dep := range deps {
			if dep == target {
				depMap[p] = struct{}{}
				break
			}
		}
	}

	deps := make([]string, 0, len(depMap))
	fmt.Println(len(deps))
	for dep := range depMap {
		deps = append(deps, dep)
	}
	sort.Strings(deps)

	for _, dep := range deps {
		fmt.Println(dep)
	}
}
