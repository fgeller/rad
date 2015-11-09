package main

import (
	"../shared"

	"os"
	"testing"
)

func TestParseD3String(t *testing.T) {
	path := "testdata/d3/API-Reference.html"

	fh, err := os.Open(path)
	if err != nil {
		t.Errorf("failed to open test file: %v", err)
		return
	}

	results := parseD3DocFile(path, fh)
	if len(results) < 1 {
		t.Errorf("expected results when parsing [%v], but got: %v", path, results)
		return
	}

	expectedStart := shared.Namespace{
		Path: "",
		Members: []shared.Member{
			{
				Name:   "arc",
				Target: "testdata/d3/SVG-Shapes.html#_arc",
			},
			{
				Name:   "area",
				Target: "testdata/d3/SVG-Shapes.html#_area_radial",
			},
		},
	}

	actualStart := results[0]
	actualStart.Members = actualStart.Members[:2]

	if !expectedStart.Eq(actualStart) {
		t.Errorf("expected start result to be\n%v\nbut got\n%v", expectedStart, actualStart)
		return
	}

	expectedSecond := shared.Namespace{
		Path: "albers",
		Members: []shared.Member{
			{
				Name:   "parallels",
				Target: "testdata/d3/Geo-Projections.html#albers_parallels",
			},
		},
	}

	actualSecond := results[1]
	actualSecond.Members = actualSecond.Members[:1]

	if !expectedSecond.Eq(actualSecond) {
		t.Errorf("expected second result to be\n%v\nbut got\n%v", expectedSecond, actualSecond)
		return
	}
}
