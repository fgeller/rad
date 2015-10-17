package main

import (
	"../shared"

	"os"
	"testing"
)

func TestPy27ParseFileMethods(t *testing.T) {
	path := "testdata/py27/datetime.html"

	fh, err := os.Open(path)
	if err != nil {
		t.Errorf("failed to open test file: %v", err)
		return
	}

	results := parsePy27DocFile(path, fh)
	if len(results) < 1 {
		t.Errorf("expected results when parsing [%v], but got: %v", path, results)
		return
	}

	expectedStart := shared.Namespace{
		Path: "",
		Members: []shared.Member{
			{
				Name:   "datetime",
				Target: "testdata/py27/datetime.html#module-datetime",
			},
		},
	}
	expectedSecond := shared.Namespace{
		Path: "datetime",
		Members: []shared.Member{
			{
				Name:   "MINYEAR",
				Target: "testdata/py27/datetime.html#datetime.MINYEAR",
			},
		},
	}

	actualStart := results[0]
	actualStart.Members = actualStart.Members[:1]

	if !expectedStart.Eq(actualStart) {
		t.Errorf("expected first result to be\n%v\nbut got\n%v", expectedStart, actualStart)
		return
	}

	actualSecond := results[1]
	actualSecond.Members = actualSecond.Members[:1]

	if !expectedSecond.Eq(actualSecond) {
		t.Errorf("expected second result to be\n%v\nbut got\n%v", expectedSecond, actualSecond)
		return
	}
}
