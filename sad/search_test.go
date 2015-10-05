package main

import (
	"../shared"

	"reflect"
	"regexp"
	"testing"
	"time"
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

	testData := []struct {
		name     string
		packPat  string
		pathPat  string
		memPat   string
		expected []searchResult
	}{
		{
			name:    "exact matching on full path",
			packPat: "go",
			pathPat: "io.ioutil",
			memPat:  "ReadAll",
			expected: []searchResult{
				{
					Namespace: []string{"io", "ioutil"},
					Member:    "ReadAll",
					Target:    "/pack/",
				},
			},
		},

		{
			name:    "regexp matches 1",
			packPat: "o",
			pathPat: "ou",
			memPat:  "ea",
			expected: []searchResult{
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
			},
		},

		{
			name:    "regexp matches 2",
			packPat: "g.",
			pathPat: "i.\\.i.u.i.",
			memPat:  "^Rea.+$",
			expected: []searchResult{
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
			},
		},

		{
			name:    "case insensitive when all lower case",
			packPat: "go",
			pathPat: "io.ioutil",
			memPat:  "readall",
			expected: []searchResult{
				{
					Namespace: []string{"io", "ioutil"},
					Member:    "ReadAll",
					Target:    "/pack/",
				},
			},
		},

		{
			name:    "empty string matches anything",
			packPat: "go",
			pathPat: "io.ioutil",
			memPat:  "",
			expected: []searchResult{
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
			},
		},
	}

	for _, data := range testData {
		params, err := compileParams(data.packPat, data.pathPat, data.memPat)
		if err != nil {
			t.Errorf("Unexpected error compiling params for test [%v]: %v", data.name, err)
			return
		}

		results := make(chan searchResult)
		control := make(chan bool)
		go find(results, control, params)

		var actual []searchResult
	readresults:
		for {
			select {
			case <-control:
				break readresults
			case r := <-results:
				actual = append(actual, r)
			}
		}

		if err != nil {
			t.Errorf("Unexpected error for test %v: %v", data.name, err)
			return
		}
		if !reflect.DeepEqual(actual, data.expected) {
			t.Errorf("Test [%v] expected\n%v\nbut got:\n%v\n", data.name, data.expected, actual)
			return
		}
	}
}

func TestFindObeysControl(t *testing.T) {
	lots := []shared.Namespace{}
	for i := 0; i < 1000; i++ {
		lots = append(
			lots,
			shared.Namespace{
				Path:    []string{"io", "ioutil"},
				Members: []shared.Member{{Name: "ReadAll" + string(i)}},
			},
		)
	}
	docs = map[string][]shared.Namespace{"go": lots}
	params := searchParams{
		pack:   regexp.MustCompile("."),
		path:   regexp.MustCompile("."),
		member: regexp.MustCompile("."),
	}

	var writtenResults []searchResult
	results := make(chan searchResult)
	control := make(chan bool)
	go func() {
		control <- true
	}()
	go find(results, control, params)
	go func() {
		time.Sleep(time.Millisecond)
		for {
			writtenResults = append(writtenResults, <-results)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	if len(writtenResults) >= len(lots) {
		t.Errorf("Expected find to stop searching but found element on channel.\n")
		return
	}
}
