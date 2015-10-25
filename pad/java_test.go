package main

import (
	"../shared"
	"os"
	"testing"
)

func TestJavaParseFileMethods(t *testing.T) {

	path := "testdata/SAXParser.html"

	fh, err := os.Open(path)
	if err != nil {
		t.Errorf("failed to open test file: %v", err)
		return
	}

	results := parseJavaDocFile(path, fh)
	if len(results) < 1 {
		t.Errorf("expected results when parsing [%v], but got: %v", path, results)
		return
	}

	expectedStart := shared.Namespace{
		Path: "java.lang.Object.SAXParser",
		Members: []shared.Member{
			{Name: "Object.html", Target: "testdata/SAXParser.html#methods.inherited.from.class.java.lang.Object"}, // TODO: Object.html??
			{Name: "clone", Target: "testdata/SAXParser.html#methods.inherited.from.class.java.lang.Object"},
		},
	}

	actualStart := results[0]
	actualStart.Members = actualStart.Members[:2]

	if !expectedStart.Eq(actualStart) {
		t.Errorf("expected first results to be\n%v\nbut got\n%v", expectedStart, actualStart)
		return
	}
}

func TestJavaParseFileFields(t *testing.T) {

	path := "testdata/ActionEvent.html"
	fh, err := os.Open(path)
	if err != nil {
		t.Errorf("failed to open test file: %v", err)
		return
	}

	results := parseJavaDocFile(path, fh)
	if len(results) < 2 {
		t.Errorf("expected results when parsing [%v], but got: %v", path, results)
		return
	}

	fstExpected := shared.Namespace{
		Path:    "java.awt.AWTEvent.ActionEvent",
		Members: []shared.Member{{Name: "ACTION_EVENT_MASK", Target: "testdata/ActionEvent.html#fields.inherited.from.class.java.awt.AWTEvent"}},
	}

	fstActual := results[0]
	fstActual.Members = fstActual.Members[:1]

	if !fstExpected.Eq(fstActual) {
		t.Errorf("expected first result to be\n%v\nbut got\n%v", fstExpected, fstActual)
		return
	}

	var foundGetSource bool
	for _, n := range results {
		for _, m := range n.Members {
			if n.Last() == "ActionEvent" && m.Name == "getSource" {
				foundGetSource = true
			}
		}
	}

	if !foundGetSource {
		t.Errorf("expected to find inherited getSource method, but wasn't found\n")
		return
	}

	var foundActionEventMask bool
	for _, n := range results {
		for _, m := range n.Members {
			if n.Last() == "ActionEvent" && m.Name == "ACTION_EVENT_MASK" {
				foundActionEventMask = true
			}
		}
	}

	if !foundActionEventMask {
		t.Errorf("expected to find inherited ACTION_EVENT_MASK field, but wasn't found\n")
		return
	}
}

func TestParseHref(t *testing.T) {
	// 2015/09/14 21:31:46 href ../../../../com/sun/source/util/DocTreeScanner.html#DocTreeScanner--
	// 2015/09/14 21:31:46 fstart 20 fend 14 -- for last DocTreeScanner.html#DocTreeScanner--

	path := "/x/y/z"

	href := "../../../../com/sun/source/util/DocTreeScanner.html#DocTreeScanner--"
	expected := shared.Namespace{
		Path:    "com.sun.source.util.DocTreeScanner",
		Members: []shared.Member{{Name: "DocTreeScanner", Target: "/x/y/DocTreeScanner.html#DocTreeScanner--"}},
	}
	actual := parseJavaHref(href, path)
	if !expected.Eq(actual) {
		t.Errorf("expected to parse\n%v\nto entry\n%v\nbut got\n%v\n", href, expected, actual)
		return
	}
}
