package main

import (
	"../shared"

	"os"
	"testing"
)

func TestDjangoParseFileMethods(t *testing.T) {
	path := "testdata/django/django.readthedocs.org/en/1.6.x/ref/request-response.html"

	fh, err := os.Open(path)
	if err != nil {
		t.Errorf("failed to open test file: %v", err)
		return
	}

	results := parseDjangoDocFile(path, fh)
	if len(results) < 1 {
		t.Errorf("expected results when parsing [%v], but got: %v", path, results)
		return
	}

	expectedStart := shared.Namespace{
		Path: "django.http",
		Members: []shared.Member{
			{
				Name:   "HttpRequest",
				Target: path + "#django.http.HttpRequest",
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
