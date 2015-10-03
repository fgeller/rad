package main

import (
	"../shared"
	"os"
	"testing"
)

func TestGoParseFileMethods(t *testing.T) {

	path := "testdata/go/godoc.org/archive/tar.html"

	fh, err := os.Open(path)
	if err != nil {
		t.Errorf("failed to open test file: %v", err)
		return
	}

	results := parseGoDocFile(path, fh)
	if len(results) < 1 {
		t.Errorf("expected results when parsing [%v], but got: %v", path, results)
		return
	}

	expectedStart := shared.Namespace{
		Path: []string{"archive", "tar"},
		Members: []shared.Member{
			{
				Name:   "Header",
				Target: "testdata/go/godoc.org/archive/tar.html#Header",
			},
		},
	}

	actualStart := results[0] // TODO: add more test cases

	actualStart.Members = actualStart.Members[:1] // TODO: fail gracefully

	if !expectedStart.Eq(actualStart) {
		t.Errorf("expected first results to be\n%v\nbut got\n%v", expectedStart, actualStart)
		return
	}

}
