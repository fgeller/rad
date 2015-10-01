package main

import (
	"../shared"
	"reflect"
	"testing"
)

func TestNewSearchResult(t *testing.T) {
	e := shared.Entry{
		Name: "entity",
	}

	expected := searchResult{
		Entity: "entity",
	}

	actual := NewSearchResult(e, 0)

	if !reflect.DeepEqual(expected, actual) {

		t.Errorf(
			"Expected graceful handling of missing members. Expected\n%v\ngot\n%v\n",
			expected,
			actual,
		)
	}
}

func TestFindPackageByPrefix(t *testing.T) {
	docs = map[string][]shared.Entry{
		"aa": []shared.Entry{{Name: "entity1", Members: []shared.Member{{Name: "member1"}}}},
		"ab": []shared.Entry{{Name: "entity1", Members: []shared.Member{{Name: "member1"}}}},
		"cd": []shared.Entry{{Name: "entity1", Members: []shared.Member{{Name: "member1"}}}},
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
	docs = map[string][]shared.Entry{
		"scala": []shared.Entry{
			shared.Entry{
				Namespace: []string{"scala", "sys"},
				Name:      "SystemProperties",
				Members:   []shared.Member{{Name: "", Signature: ""}},
			},
			shared.Entry{
				Namespace: []string{"scala", "collection"},
				Name:      "SetProxy",
				Members:   []shared.Member{{Name: "", Signature: ""}},
			},
		},
	}
	es, err := findEntityMember("scala", "SetProxy", "", 10)

	if err != nil {
		t.Errorf("unexpected error [%v]", err)
		return
	}

	if len(es) != 1 {
		t.Errorf("expected to find one entry but got [%v]", es)
		return
	}

	if es[0].Namespace[0] != "scala" ||
		es[0].Namespace[1] != "collection" ||
		es[0].Entity != "SetProxy" {
		t.Errorf("expected to find SetProxy entry but got [%v]", es[0])
		return
	}

}

func TestFindEntityByPrefix(t *testing.T) {
	docs = map[string][]shared.Entry{
		"scala": []shared.Entry{
			shared.Entry{
				Namespace: []string{"scala", "sys"},
				Name:      "SystemProperties",
				Members:   []shared.Member{{Name: "", Signature: ""}},
			},
			shared.Entry{
				Namespace: []string{"scala", "sys"},
				Name:      "SystemThings",
				Members:   []shared.Member{{Name: "", Signature: ""}},
			},
			shared.Entry{
				Namespace: []string{"scala", "collection"},
				Name:      "SetProxy",
				Members:   []shared.Member{{Name: "", Signature: ""}},
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
	docs = map[string][]shared.Entry{
		"scala": []shared.Entry{
			shared.Entry{
				Namespace: []string{"scala", "sys"},
				Name:      "SystemProperties",
				Members:   []shared.Member{{Name: "hans", Signature: ""}},
			},
			shared.Entry{
				Namespace: []string{"scala", "sys"},
				Name:      "SYSTEMThings",
				Members:   []shared.Member{{Name: "HANS", Signature: ""}},
			},
			shared.Entry{
				Namespace: []string{"scala", "sys"},
				Name:      "systemThings",
				Members:   []shared.Member{{Name: "hAnS", Signature: ""}},
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
	docs = map[string][]shared.Entry{
		"scala": []shared.Entry{
			shared.Entry{
				Namespace: []string{"scala", "collection", "mutable"},
				Name:      "HashMap",
				Members:   []shared.Member{{Name: "clearTable", Signature: "():Unit"}},
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
	docs = map[string][]shared.Entry{
		"scala": []shared.Entry{
			shared.Entry{
				Namespace: []string{"scala", "collection", "mutable"},
				Name:      "HashMap",
				Members:   []shared.Member{{Name: "clearTable", Signature: "():Unit"}},
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
