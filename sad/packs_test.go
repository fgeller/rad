package main

import (
	"testing"
)

func TestReadArchivePackInfo(t *testing.T) {
	target := "testdata/go-html.zip"

	actual, err := readArchivePackInfo(target)
	if err != nil {
		t.Errorf("Errow reading archive's pack info: %v", err)
		return
	}

	if actual.Name != "go-html" || actual.Type != "go" {
		t.Errorf("Expected go-html/go pack info, got: %v", actual)
		return
	}
}
