package main

import (
	"../shared"

	"log"
	"os"
	"testing"
)

func TestReactParseFileMethods(t *testing.T) {
	path := "testdata/react/component-specs.html"

	fh, err := os.Open(path)
	if err != nil {
		t.Errorf("failed to open test file: %v", err)
		return
	}

	results := parseReactDocFile(path, fh)
	if len(results) < 1 {
		t.Errorf("expected results when parsing [%v], but got: %v", path, results)
		return
	}

	var foundComponentWillMount bool
	var foundGetDefaultProps bool

	for _, ns := range results {
		for _, m := range ns.Members {
			if m.Name == "componentWillMount" {
				foundComponentWillMount = true
			}
			if m.Name == "getDefaultProps" {
				foundGetDefaultProps = true
			}
		}
	}

	if !foundGetDefaultProps {
		t.Errorf("Expected to find getDefaultProps")
		return
	}

	if !foundComponentWillMount {
		t.Errorf("Expected to find componentWillMount")
		return
	}

	expectedStart := shared.Namespace{
		Path: "component-specs",
		Members: []shared.Member{
			{
				Name:   "render",
				Target: "testdata/react/component-specs.html#render",
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

func TestReactParseComponentApi(t *testing.T) {
	path := "testdata/react/component-api.html"

	fh, err := os.Open(path)
	if err != nil {
		t.Errorf("failed to open test file: %v", err)
		return
	}

	results := parseReactDocFile(path, fh)
	if len(results) < 1 {
		t.Errorf("expected results when parsing [%v], but got: %v", path, results)
		return
	}

	var foundGetDOMNode bool
	var foundReplaceProps bool

	for _, ns := range results {
		for _, m := range ns.Members {
			if m.Name == "getDOMNode" {
				foundGetDOMNode = true
			}
			if m.Name == "replaceProps" {
				foundReplaceProps = true
			}
		}
	}

	if !foundGetDOMNode {
		t.Errorf("Expected to find getDOMNode")
		return
	}

	if !foundReplaceProps {
		t.Errorf("Expected to find replaceProps")
		return
	}

	expectedStart := shared.Namespace{
		Path: "component-api",
		Members: []shared.Member{
			{
				Name:   "setState",
				Target: "testdata/react/component-api.html#setstate",
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

func TestReactParseTopLevelApi(t *testing.T) {
	path := "testdata/react/top-level-api.html"

	fh, err := os.Open(path)
	if err != nil {
		t.Errorf("failed to open test file: %v", err)
		return
	}

	results := parseReactDocFile(path, fh)
	if len(results) < 1 {
		t.Errorf("expected results when parsing [%v], but got: %v", path, results)
		return
	}

	var foundCreateClass bool
	var foundRenderToStaticMarkup bool

	for _, ns := range results {
		log.Printf("Found namespace with path: %v\n", ns.Path)
		for _, m := range ns.Members {
			if m.Name == "createClass" {
				foundCreateClass = true
			}
			if m.Name == "renderToStaticMarkup" {
				foundRenderToStaticMarkup = true
			}
			log.Printf("  Member: [%v] Target: [%v]\n", m.Name, m.Target)
		}
	}

	if !foundCreateClass {
		t.Errorf("Expected to find createClass")
		return
	}

	if !foundRenderToStaticMarkup {
		t.Errorf("Expected to find renderToStaticMarkup")
		return
	}

	expectedStart := shared.Namespace{
		Path: "top-level-api.React",
		Members: []shared.Member{
			{
				Name:   "Component",
				Target: "testdata/react/top-level-api.html#react.component",
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
