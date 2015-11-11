package main

import (
	"../shared"

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

func TestNameCountForInstalledPacks(t *testing.T) {
	resetGlobals()
	loadPack(
		shared.Pack{Name: "blubb"},
		[]shared.Namespace{{Members: []shared.Member{{Name: "m2"}}}},
	)

	actual := installedPacks()
	if len(actual) != 1 {
		t.Errorf("Expected one installed pack got %v\n", actual)
		return
	}

	if actual[0].NameCount != 1 {
		t.Errorf(
			"Expected one name for installed pack got %v where pack=%v\n",
			actual[0].NameCount,
			actual,
		)
		return
	}
}
