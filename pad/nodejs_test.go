package main

import (
	"../shared"

	"os"
	"testing"
)

func TestParseNodeJsConsole(t *testing.T) {
	path := "testdata/nodejs/console.html"

	fh, err := os.Open(path)
	if err != nil {
		t.Errorf("failed to open test file: %v", err)
		return
	}

	results := parseNodeJsDocFile(path, fh)
	if len(results) < 1 {
		t.Errorf("expected results when parsing [%v], but got: %v", path, results)
		return
	}

	expectedStart := shared.Namespace{
		Path: "",
		Members: []shared.Member{
			{
				Name:   "Console",
				Target: "testdata/nodejs/console.html#console_console",
			},
		},
	}

	actualStart := results[0]
	actualStart.Members = actualStart.Members[:1]

	if !expectedStart.Eq(actualStart) {
		t.Errorf("expected first result to be\n%v\nbut got\n%v", expectedStart, actualStart)
		return
	}

	var foundNewConsole bool

	for _, ns := range results {
		for _, m := range ns.Members {
			if m.Name == "new Console(stdout, stderr)" {
				foundNewConsole = true
			}
		}
	}

	if !foundNewConsole {
		t.Errorf("Expected to find new Console(stdout, stderr)")
		return
	}

	var foundAssert bool

	for _, ns := range results {
		for _, m := range ns.Members {
			if m.Name == "assert(value, message, ...)" {
				foundAssert = true
			}
		}
	}

	if !foundAssert {
		t.Errorf("Expected to find assert(value, message, ...)")
		return
	}

}
