package main

import (
	"../shared"
	"reflect"
	"testing"
)

func TestNewSearchResult(t *testing.T) {
	n := shared.Namespace{
		Path: []string{"entity"},
	}

	expected := searchResult{
		Namespace: []string{"entity"},
	}

	actual := NewSearchResult(n, 0)

	if !reflect.DeepEqual(expected, actual) {

		t.Errorf(
			"Expected graceful handling of missing members. Expected\n%v\ngot\n%v\n",
			expected,
			actual,
		)
	}
}

func TestFind(t *testing.T) {

	docs = map[string][]shared.Namespace{
		"go": []shared.Namespace{
			{
				Path:    []string{"io", "ioutil"},
				Members: []shared.Member{{Name: "ReadAll"}, {Name: "ReadDir"}},
			},
		},
	}

	// exact matching on full path
	limit := 10
	packPat := "go"
	pathPat := "io.ioutil"
	memPat := "ReadAll"
	expected := []searchResult{
		{
			Namespace: []string{"io", "ioutil"},
			Member:    "ReadAll",
			Target:    "/pack/",
		},
	}

	actual, err := find(packPat, pathPat, memPat, limit)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected\n%v\nbut got:\n%v\n", expected, actual)
		return
	}

	// regexp matches
	limit = 10
	packPat = "o"
	pathPat = "ou"
	memPat = "ea"
	expected = []searchResult{
		{
			Namespace: []string{"io", "ioutil"},
			Member:    "ReadAll",
			Target:    "/pack/",
		},
		{
			Namespace: []string{"io", "ioutil"},
			Member:    "ReadDir",
			Target:    "/pack/",
		},
	}

	actual, err = find(packPat, pathPat, memPat, limit)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected\n%v\nbut got:\n%v\n", expected, actual)
		return
	}

	limit = 10
	packPat = "g."
	pathPat = "i.\\.i.u.i."
	memPat = "^Rea.+$"
	expected = []searchResult{
		{
			Namespace: []string{"io", "ioutil"},
			Member:    "ReadAll",
			Target:    "/pack/",
		},
		{
			Namespace: []string{"io", "ioutil"},
			Member:    "ReadDir",
			Target:    "/pack/",
		},
	}

	actual, err = find(packPat, pathPat, memPat, limit)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected\n%v\nbut got:\n%v\n", expected, actual)
		return
	}

	// limit results
	limit = 1
	packPat = "o"
	pathPat = "ou"
	memPat = "ea"
	expected = []searchResult{
		{
			Namespace: []string{"io", "ioutil"},
			Member:    "ReadAll",
			Target:    "/pack/",
		},
	}

	actual, err = find(packPat, pathPat, memPat, limit)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected\n%v\nbut got:\n%v\n", expected, actual)
		return
	}

	// case insensitive when all lower case
	limit = 10
	packPat = "go"
	pathPat = "io.ioutil"
	memPat = "readall"
	expected = []searchResult{
		{
			Namespace: []string{"io", "ioutil"},
			Member:    "ReadAll",
			Target:    "/pack/",
		},
	}

	actual, err = find(packPat, pathPat, memPat, limit)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected\n%v\nbut got:\n%v\n", expected, actual)
		return
	}

	// empty string matches anything
	limit = 10
	packPat = "go"
	pathPat = "io.ioutil"
	memPat = ""
	expected = []searchResult{
		{
			Namespace: []string{"io", "ioutil"},
			Member:    "ReadAll",
			Target:    "/pack/",
		},
		{
			Namespace: []string{"io", "ioutil"},
			Member:    "ReadDir",
			Target:    "/pack/",
		},
	}

	actual, err = find(packPat, pathPat, memPat, limit)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected\n%v\nbut got:\n%v\n", expected, actual)
		return
	}
}
