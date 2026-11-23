package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"--help"}, &out, &errOut); code != 0 {
		t.Fatalf("exit code %d, want 0", code)
	}
	if !strings.Contains(out.String(), "usage: codesearch") {
		t.Errorf("help text missing usage line: %q", out.String())
	}
}

func TestVersion(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"--version"}, &out, &errOut); code != 0 {
		t.Fatalf("exit code %d, want 0", code)
	}
	if !strings.HasPrefix(out.String(), "codesearch ") {
		t.Errorf("unexpected version output: %q", out.String())
	}
}

func TestStubExitsNonZero(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"serve"}, &out, &errOut); code == 0 {
		t.Fatal("serve stub exited 0")
	}
	if got := errOut.String(); !strings.Contains(got, "not implemented yet: see #5") {
		t.Errorf("stub message = %q", got)
	}
}
