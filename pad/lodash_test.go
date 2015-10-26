package main

import (
	"../shared"

	"os"
	"testing"
)

func TestLodashParseFileMethods(t *testing.T) {
	path := "testdata/lodash/docs.html"

	fh, err := os.Open(path)
	if err != nil {
		t.Errorf("failed to open test file: %v", err)
		return
	}

	results := parseLodashDocFile(path, fh)
	if len(results) < 1 {
		t.Errorf("expected results when parsing [%v], but got: %v", path, results)
		return
	}

	expectedStart := shared.Namespace{
		Path: "Array",
		Members: []shared.Member{
			{
				Name:   "chunk",
				Target: path + "#chunk",
			},
		},
	}

	actualStart := results[0]
	actualStart.Members = actualStart.Members[:1]

	if !expectedStart.Eq(actualStart) {
		t.Errorf("expected first result to be\n%v\nbut got\n%v", expectedStart, actualStart)
		return
	}

	foundEscape := false
	for _, ns := range results {
		if ns.Path == "Properties.templateSettings" {
			for _, m := range ns.Members {
				if m.Name == "escape" {
					foundEscape = true
				}
			}
		}
	}

	if !foundEscape {
		t.Errorf("Expected to find Properties.templateSettings.escape")
		return
	}
}
