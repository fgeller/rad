package main

import (
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
	if len(results) < 2 {
		t.Errorf("expected results when parsing [%v], but got: %v", path, results)
		return
	}

	fstExpected := entry{
		Namespace: []string{"javax", "xml", "parsers"},
		Entity:    "SAXParser",
		Member:    "SAXParser",
		Signature: "",
		Target:    "testdata/SAXParser.html#SAXParser--",
		Source:    path,
	}

	sndExpected := entry{
		Namespace: []string{"javax", "xml", "parsers"},
		Entity:    "SAXParser",
		Member:    "getParser",
		Signature: "",
		Target:    "testdata/SAXParser.html#getParser--",
		Source:    path,
	}

	if !fstExpected.eq(results[0]) {
		t.Errorf("expected first result to be\n%v\nbut got\n%v", fstExpected, results[0])
		return
	}

	if !sndExpected.eq(results[1]) {
		t.Errorf("expected second result to be\n%v\nbut got\n%v", sndExpected, results[1])
		return
	}

	var foundClone bool
	for _, e := range results {
		if e.Entity == "Object" && e.Member == "clone" {
			foundClone = true
		}
	}

	if !foundClone {
		t.Errorf("expected to find inherited clone method, but wasn't found\n")
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

	fstExpected := entry{
		Namespace: []string{"java", "awt", "event"},
		Entity:    "ActionEvent",
		Member:    "ACTION_FIRST",
		Signature: "",
		Target:    "testdata/ActionEvent.html#ACTION_FIRST",
		Source:    path,
	}

	if !fstExpected.eq(results[0]) {
		t.Errorf("expected first result to be\n%v\nbut got\n%v", fstExpected, results[0])
		return
	}

	var foundGetSource bool
	for _, e := range results {
		if e.Entity == "EventObject" && e.Member == "getSource" {
			foundGetSource = true
		}
	}

	if !foundGetSource {
		t.Errorf("expected to find inherited getSource method, but wasn't found\n")
		return
	}

	var foundActionEventMask bool
	for _, e := range results {
		if e.Entity == "AWTEvent" && e.Member == "ACTION_EVENT_MASK" {
			foundActionEventMask = true
		}
	}

	if !foundActionEventMask {
		t.Errorf("expected to find inherited ACTION_EVENT_MASK field, but wasn't found\n")
		return
	}
}
