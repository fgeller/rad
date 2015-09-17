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
		Name:      "SAXParser",
		Members:   []member{{Name: "SAXParser", Signature: "", Target: "testdata/SAXParser.html#SAXParser--", Source: path}},
		Source:    path,
	}

	sndExpected := entry{
		Namespace: []string{"javax", "xml", "parsers"},
		Name:      "SAXParser",
		Members:   []member{{Name: "getParser", Signature: "", Target: "testdata/SAXParser.html#getParser--", Source: path}},
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
		if e.Name == "SAXParser" {
			for _, m := range e.Members {
				if m.Name == "clone" &&
					m.Target == "testdata/SAXParser.html#methods.inherited.from.class.java.lang.Object" {
					foundClone = true
				}
			}
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
		Name:      "ActionEvent",
		Members:   []member{{Name: "ACTION_FIRST", Signature: "", Target: "testdata/ActionEvent.html#ACTION_FIRST", Source: path}},
		Source:    path,
	}

	if !fstExpected.eq(results[0]) {
		t.Errorf("expected first result to be\n%v\nbut got\n%v", fstExpected, results[0])
		return
	}

	var foundGetSource bool
	for _, e := range results {
		if e.Name == "ActionEvent" && e.Members[0].Name == "getSource" {
			foundGetSource = true
		}
	}

	if !foundGetSource {
		t.Errorf("expected to find inherited getSource method, but wasn't found\n")
		return
	}

	var foundActionEventMask bool
	for _, e := range results {
		if e.Name == "ActionEvent" && e.Members[0].Name == "ACTION_EVENT_MASK" {
			foundActionEventMask = true
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
	expected := entry{
		Namespace: []string{"com", "sun", "source", "util"},
		Name:      "DocTreeScanner",
		Members:   []member{{Name: "DocTreeScanner", Signature: "", Source: path, Target: "/x/y/DocTreeScanner.html#DocTreeScanner--"}},
		Source:    path,
	}
	actual := parseHref(href, path)
	if !expected.eq(actual) {
		t.Errorf("expected to parse\n%v\nto entry\n%v\nbut got\n%v\n", href, expected, actual)
		return
	}
}
