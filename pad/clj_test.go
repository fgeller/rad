package main

import (
	"../shared"

	"os"
	"testing"
)

func TestClojureParseFileMethods(t *testing.T) {
	path := "testdata/clj/clojure.core.async/<!!.html"

	fh, err := os.Open(path)
	if err != nil {
		t.Errorf("failed to open test file: %v", err)
		return
	}

	results := parseClojureDocFile(path, fh)
	if len(results) < 1 {
		t.Errorf("expected results when parsing [%v], but got: %v", path, results)
		return
	}

	expectedStart := shared.Namespace{
		Path: []string{"clojure", "core", "async"},
		Members: []shared.Member{
			{
				Name:   "<!!",
				Target: "testdata/clj/clojure.core.async/<!!.html",
			},
		},
	}

	actualStart := results[0]
	actualStart.Members = actualStart.Members[:1]

	if !expectedStart.Eq(actualStart) {
		t.Errorf("expected first results to be\n%v\nbut got\n%v", expectedStart, actualStart)
		return
	}
}
