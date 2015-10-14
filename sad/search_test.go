package main

import (
	"../shared"

	"log"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func TestNewSearchResult(t *testing.T) {
	n := shared.Namespace{
		Path: "entity",
	}

	expected := searchResult{
		Namespace: "entity",
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
	resetGlobals()

	pck := shared.Pack{Name: "go"}
	nss := []shared.Namespace{
		{
			Path:    "io.ioutil",
			Members: []shared.Member{{Name: "ReadAll"}, {Name: "ReadDir"}},
		},
	}
	installPack(pck, nss)

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
					Namespace: "io.ioutil",
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
					Namespace: "io.ioutil",
					Member:    "ReadAll",
					Target:    "/pack/",
				},
				{
					Namespace: "io.ioutil",
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
					Namespace: "io.ioutil",
					Member:    "ReadAll",
					Target:    "/pack/",
				},
				{
					Namespace: "io.ioutil",
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
					Namespace: "io.ioutil",
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
					Namespace: "io.ioutil",
					Member:    "ReadAll",
					Target:    "/pack/",
				},
				{
					Namespace: "io.ioutil",
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
		control := make(chan struct{})
		go find(results, control, params)

		var actual []searchResult
	readresults:
		for {
			r, ok := <-results
			if !ok {
				break readresults
			}
			actual = append(actual, r)
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
				Path:    "io.ioutil",
				Members: []shared.Member{{Name: "ReadAll" + string(i)}},
			},
		)
	}

	pck := shared.Pack{Name: "go"}
	installPack(pck, lots)

	params := searchParams{
		pack:   regexp.MustCompile("."),
		path:   regexp.MustCompile("."),
		member: regexp.MustCompile("."),
	}

	var writtenResults []searchResult
	results := make(chan searchResult)
	control := make(chan struct{}, 1)
	control <- struct{}{}
	go find(results, control, params)
	go func() {
		for {
			res, ok := <-results
			if !ok {
				return
			}
			writtenResults = append(writtenResults, res)
		}
	}()

	time.Sleep(1 * time.Millisecond)
	if len(writtenResults) >= 1 {
		t.Errorf("Expected find to stop searching but found %v element(s) on channel.\n", len(writtenResults))
		return
	}
}
