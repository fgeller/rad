package main

import (
	"../shared"

	"os"
	"testing"
)

func TestParsePHP(t *testing.T) {
	path := "testdata/php/indexes.functions.php.html"

	fh, err := os.Open(path)
	if err != nil {
		t.Errorf("failed to open test file: %v", err)
		return
	}

	results := parsePHPDocFile(path, fh)
	if len(results) != 3 {
		t.Errorf("expected 3 results when parsing [%v], but got: %v", path, results)
		return
	}

	expectedStart := shared.Namespace{
		Path: "",
		Members: []shared.Member{
			{
				Name:   "abs",
				Target: "testdata/php/function.abs.php.html",
			},
			{
				Name:   "json_decode",
				Target: "testdata/php/function.json-decode.php.html",
			},
		},
	}
	expectedSecond := shared.Namespace{
		Path: "Judy",
		Members: []shared.Member{
			{
				Name:   "byCount",
				Target: "testdata/php/judy.bycount.php.html",
			},
		},
	}
	expectedThird := shared.Namespace{
		Path: "ZMQSocket",
		Members: []shared.Member{
			{
				Name:   "__construct",
				Target: "testdata/php/zmqsocket.construct.php.html",
			},
		},
	}

	actualStart := results[0]
	if !expectedStart.Eq(actualStart) {
		t.Errorf("expected first result to be\n%v\nbut got\n%v", expectedStart, actualStart)
		return
	}

	actualSecond := results[1]
	if !expectedSecond.Eq(actualSecond) {
		t.Errorf("expected second result to be\n%v\nbut got\n%v", expectedSecond, actualSecond)
		return
	}

	actualThird := results[2]
	if !expectedThird.Eq(actualThird) {
		t.Errorf("expected third result to be\n%v\nbut got\n%v", expectedThird, actualThird)
		return
	}
}
