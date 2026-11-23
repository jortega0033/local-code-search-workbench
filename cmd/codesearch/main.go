// Command codesearch builds, serves and verifies a local code search snapshot.
package main

import (
	"fmt"
	"io"
	"os"
)

const version = "0.0.0-scaffold"

const usage = `usage: codesearch <command>

commands:
  generate   build the seeded corpus (not implemented yet: see #2)
  serve      run the local search server (not implemented yet: see #5)
  verify     run the parity and timing checks (not implemented yet: see #7)

flags:
  --help     show this help
  --version  show the version
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}
	switch args[0] {
	case "--help", "-h", "help":
		fmt.Fprint(stdout, usage)
		return 0
	case "--version", "version":
		fmt.Fprintln(stdout, "codesearch "+version)
		return 0
	case "generate":
		return notBuilt(stderr, 2)
	case "serve":
		return notBuilt(stderr, 5)
	case "verify":
		return notBuilt(stderr, 7)
	}
	fmt.Fprintf(stderr, "unknown command %q\n\n%s", args[0], usage)
	return 2
}

func notBuilt(stderr io.Writer, ticket int) int {
	fmt.Fprintf(stderr, "not implemented yet: see #%d\n", ticket)
	return 1
}
