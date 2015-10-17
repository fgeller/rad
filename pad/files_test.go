package main

import (
	"../shared"

	"io"
	"reflect"
	"testing"
)

func TestMergingNamespaces(t *testing.T) {
	path := "testdata"
	namespaces := []shared.Namespace{
		{
			Path: "abc",
			Members: []shared.Member{
				{Name: "def", Target: ""},
				{Name: "ghi", Target: ""},
			},
		},
		{
			Path: "abc",
			Members: []shared.Member{
				{Name: "def", Target: ""},
				{Name: "ghi", Target: ""},
			},
		},
	}
	expected := []shared.Namespace{
		{
			Path: "abc",
			Members: []shared.Member{
				{Name: "def", Target: ""},
				{Name: "ghi", Target: ""},
			},
		},
	}
	testParser := func(p string, i io.Reader) []shared.Namespace {
		return namespaces
	}

	actual, err := scan(path, testParser)
	if err != nil {
		t.Errorf("Unexpected error while scanning: %v\n", err)
		return
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected merged results, but got: %v\n", actual)
		return
	}
}
