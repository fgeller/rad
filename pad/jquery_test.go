package main

import (
	"../shared"

	"os"
	"testing"
)

func TestReactParseJQueryAjax(t *testing.T) {
	path := "testdata/jquery/jQuery.ajax/index.html"

	fh, err := os.Open(path)
	if err != nil {
		t.Errorf("failed to open test file: %v", err)
		return
	}

	results := parseJQueryDocFile(path, fh)
	if len(results) < 1 {
		t.Errorf("expected results when parsing [%v], but got: %v", path, results)
		return
	}

	var foundAjaxSettings bool

	for _, ns := range results {
		for _, m := range ns.Members {
			if m.Name == "ajax(settings)" {
				foundAjaxSettings = true
			}
		}
	}

	if !foundAjaxSettings {
		t.Errorf("Expected to find ajax(settings)")
		return
	}

	expectedStart := shared.Namespace{
		Path: "jQuery",
		Members: []shared.Member{
			{
				Name:   "ajax(url,settings)",
				Target: "testdata/jquery/jQuery.ajax/index.html#jQuery-ajax-url-settings",
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

func TestReactParseJQueryGetJSON(t *testing.T) {
	path := "testdata/jquery/jQuery.getJSON/index.html"

	fh, err := os.Open(path)
	if err != nil {
		t.Errorf("failed to open test file: %v", err)
		return
	}

	results := parseJQueryDocFile(path, fh)
	if len(results) < 1 {
		t.Errorf("expected results when parsing [%v], but got: %v", path, results)
		return
	}

	expectedStart := shared.Namespace{
		Path: "jQuery",
		Members: []shared.Member{
			{
				Name:   "getJSON(url,data,success)",
				Target: "testdata/jquery/jQuery.getJSON/index.html#jQuery-getJSON-url-data-success",
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
