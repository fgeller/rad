package main

import (
	"../shared"

	"os"
	"testing"
)

func TestParseMDNJavascriptString(t *testing.T) {
	path := "testdata/mdn/String.html"

	fh, err := os.Open(path)
	if err != nil {
		t.Errorf("failed to open test file: %v", err)
		return
	}

	results := parseMDNDocFile(path, fh)
	if len(results) < 1 {
		t.Errorf("expected results when parsing [%v], but got: %v", path, results)
		return
	}

	expectedStart := shared.Namespace{
		Path: "String",
		Members: []shared.Member{
			{
				Name:   "fromCharCode()",
				Target: "testdata/mdn/String/fromCharCode.html",
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
