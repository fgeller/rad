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

func TestFindPackageByPrefix(t *testing.T) {
	docs = map[string][]shared.Namespace{
		"aa": []shared.Namespace{{Path: []string{"entity1"}, Members: []shared.Member{{Name: "member1"}}}},
		"ab": []shared.Namespace{{Path: []string{"entity1"}, Members: []shared.Member{{Name: "member1"}}}},
		"cd": []shared.Namespace{{Path: []string{"entity1"}, Members: []shared.Member{{Name: "member1"}}}},
	}
	res, err := findEntityMember("a", "entity", "member", 10)

	if err != nil {
		t.Errorf("unexpected error when accessing packages, got [%v]", err)
	}

	if len(res) != 2 {
		t.Errorf("expected two results, got [%v]", res)
	}
}

func TestFindEntry(t *testing.T) {
	docs = map[string][]shared.Namespace{
		"scala": []shared.Namespace{
			shared.Namespace{
				Path:    []string{"scala", "sys", "SystemProperties"},
				Members: []shared.Member{{Name: ""}},
			},
			shared.Namespace{
				Path:    []string{"scala", "collection", "SetProxy"},
				Members: []shared.Member{{Name: ""}},
			},
		},
	}
	ns, err := findEntityMember("scala", "SetProxy", "", 10)

	if err != nil {
		t.Errorf("unexpected error [%v]", err)
		return
	}

	if len(ns) != 1 {
		t.Errorf("expected to find one result but got [%v]", ns)
		return
	}

	// TODO: just compare the slices?
	if ns[0].Namespace[0] != "scala" ||
		ns[0].Namespace[1] != "collection" ||
		ns[0].Namespace[len(ns[0].Namespace)-1] != "SetProxy" {
		t.Errorf("expected to find SetProxy entry but got [%v]", ns[0])
		return
	}

}

func TestFindEntityByPrefix(t *testing.T) {
	docs = map[string][]shared.Namespace{
		"scala": []shared.Namespace{
			shared.Namespace{
				Path:    []string{"scala", "sys", "SystemProperties"},
				Members: []shared.Member{{Name: ""}},
			},
			shared.Namespace{
				Path:    []string{"scala", "sys", "SystemThings"},
				Members: []shared.Member{{Name: ""}},
			},
			shared.Namespace{
				Path:    []string{"scala", "collection", "SetProxy"},
				Members: []shared.Member{{Name: ""}},
			},
		},
	}
	es, err := findEntityMember("scala", "Syst", "", 10)

	if err != nil {
		t.Errorf("unexpected error [%v]", err)
		return
	}

	if len(es) != 2 {
		t.Errorf("expected to find two entries but got [%v]", es)
		return
	}

	if !es[0].eq(NewSearchResult(docs["scala"][0], 0)) ||
		!es[1].eq(NewSearchResult(docs["scala"][1], 0)) {
		t.Errorf("expected to find System entries but got [%v]", es)
		return
	}

}

func TestFindIsCaseInsentitive(t *testing.T) {
	docs = map[string][]shared.Namespace{
		"scala": []shared.Namespace{
			shared.Namespace{
				Path:    []string{"scala", "sys", "SystemProperties"},
				Members: []shared.Member{{Name: "hans"}},
			},
			shared.Namespace{
				Path:    []string{"scala", "sys", "SYSTEMThings"},
				Members: []shared.Member{{Name: "HANS"}},
			},
			shared.Namespace{
				Path:    []string{"scala", "sys", "systemThings"},
				Members: []shared.Member{{Name: "hAnS"}},
			},
		},
	}
	es, err := findEntityMember("scala", "SyS", "haNS", 10)

	if err != nil {
		t.Errorf("unexpected error [%v]", err)
		return
	}

	if len(es) != 3 {
		t.Errorf("expected to find three entries but got [%v]", es)
		return
	}

	if !es[0].eq(NewSearchResult(docs["scala"][0], 0)) ||
		!es[1].eq(NewSearchResult(docs["scala"][1], 0)) ||
		!es[2].eq(NewSearchResult(docs["scala"][2], 0)) {
		t.Errorf("expected to find System entries but got [%v]", es)
		return
	}

}

func TestFindMember(t *testing.T) {
	docs = map[string][]shared.Namespace{
		"scala": []shared.Namespace{
			shared.Namespace{
				Path:    []string{"scala", "collection", "mutable", "HashMap"},
				Members: []shared.Member{{Name: "clearTable"}},
			},
		},
	}

	es, err := findEntityMember("scala", "HashMap", "clearTable", 10)

	if err != nil {
		t.Errorf("unexpected error [%v]", err)
		return
	}

	if len(es) != 1 {
		t.Errorf("expected to find one entry but got [%v]", es)
		return
	}

	actual := es[0]
	expected := NewSearchResult(docs["scala"][0], 0)
	if !actual.eq(expected) {
		t.Errorf("expected to find\n%v\nbut got\n%v", expected, actual)
	}

}

func TestFindMemberByPrefix(t *testing.T) {
	docs = map[string][]shared.Namespace{
		"scala": []shared.Namespace{
			shared.Namespace{
				Path:    []string{"scala", "collection", "mutable", "HashMap"},
				Members: []shared.Member{{Name: "clearTable"}},
			},
		},
	}

	es, err := findEntityMember("scala", "HashMap", "clear", 10)

	if err != nil {
		t.Errorf("unexpected error [%v]", err)
		return
	}

	if len(es) != 1 {
		t.Errorf("expected to find one entry but got [%v]", es)
		return
	}

	actual := es[0]
	expected := NewSearchResult(docs["scala"][0], 0)
	if !actual.eq(expected) {
		t.Errorf("expected to find\n%v\nbut got\n%v", expected, actual)
	}

}
