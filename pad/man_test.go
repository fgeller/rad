package main

import (
	"../shared"

	"os"
	"testing"
)

func TestParseManString(t *testing.T) {
	path := "testdata/man-pages/dir_all_alphabetic.html"

	fh, err := os.Open(path)
	if err != nil {
		t.Errorf("failed to open test file: %v", err)
		return
	}

	results := parseManDocFile(path, fh)
	if len(results) < 1 {
		t.Errorf("expected results when parsing [%v], but got: %v", path, results)
		return
	}

	expectedStart := shared.Namespace{
		Path: "man1",
		Members: []shared.Member{
			{
				Name:   "ac",
				Target: "testdata/man-pages/man1/ac.1.html",
			},
		},
	}

	actualStart := results[0]
	actualStart.Members = actualStart.Members[:1]

	if !expectedStart.Eq(actualStart) {
		t.Errorf("expected first result to be\n%v\nbut got\n%v", expectedStart, actualStart)
		return
	}
}
