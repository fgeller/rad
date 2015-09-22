package main

import "testing"

func TestFindPackageByPrefix(t *testing.T) {
	docs = map[string][]entry{
		"aa": []entry{{Name: "entity1", Members: []member{{Name: "member1"}}}},
		"ab": []entry{{Name: "entity1", Members: []member{{Name: "member1"}}}},
		"cd": []entry{{Name: "entity1", Members: []member{{Name: "member1"}}}},
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
	docs = map[string][]entry{
		"scala": []entry{
			entry{
				Namespace: []string{"scala", "sys"},
				Name:      "SystemProperties",
				Members:   []member{{Name: "", Signature: ""}},
			},
			entry{
				Namespace: []string{"scala", "collection"},
				Name:      "SetProxy",
				Members:   []member{{Name: "", Signature: ""}},
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
		es[0].Name != "SetProxy" {
		t.Errorf("expected to find SetProxy entry but got [%v]", es[0])
		return
	}

}

func TestFindEntityByPrefix(t *testing.T) {
	docs = map[string][]entry{
		"scala": []entry{
			entry{
				Namespace: []string{"scala", "sys"},
				Name:      "SystemProperties",
				Members:   []member{{Name: "", Signature: ""}},
			},
			entry{
				Namespace: []string{"scala", "sys"},
				Name:      "SystemThings",
				Members:   []member{{Name: "", Signature: ""}},
			},
			entry{
				Namespace: []string{"scala", "collection"},
				Name:      "SetProxy",
				Members:   []member{{Name: "", Signature: ""}},
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

	if !es[0].eq(docs["scala"][0]) ||
		!es[1].eq(docs["scala"][1]) {
		t.Errorf("expected to find System entries but got [%v]", es)
		return
	}

}

func TestFindIsCaseInsentitive(t *testing.T) {
	docs = map[string][]entry{
		"scala": []entry{
			entry{
				Namespace: []string{"scala", "sys"},
				Name:      "SystemProperties",
				Members:   []member{{Name: "hans", Signature: ""}},
			},
			entry{
				Namespace: []string{"scala", "sys"},
				Name:      "SYSTEMThings",
				Members:   []member{{Name: "HANS", Signature: ""}},
			},
			entry{
				Namespace: []string{"scala", "sys"},
				Name:      "systemThings",
				Members:   []member{{Name: "hAnS", Signature: ""}},
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

	if !es[0].eq(docs["scala"][0]) ||
		!es[1].eq(docs["scala"][1]) ||
		!es[2].eq(docs["scala"][2]) {
		t.Errorf("expected to find System entries but got [%v]", es)
		return
	}

}

func TestFindMember(t *testing.T) {
	docs = map[string][]entry{
		"scala": []entry{
			entry{
				Namespace: []string{"scala", "collection", "mutable"},
				Name:      "HashMap",
				Members:   []member{{Name: "clearTable", Signature: "():Unit"}},
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
	expected := docs["scala"][0]
	if !actual.eq(expected) {
		t.Errorf("expected to find\n%v\nbut got\n%v", expected, actual)
	}

}

func TestFindMemberByPrefix(t *testing.T) {
	docs = map[string][]entry{
		"scala": []entry{
			entry{
				Namespace: []string{"scala", "collection", "mutable"},
				Name:      "HashMap",
				Members:   []member{{Name: "clearTable", Signature: "():Unit"}},
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
	expected := docs["scala"][0]
	if !actual.eq(expected) {
		t.Errorf("expected to find\n%v\nbut got\n%v", expected, actual)
	}

}
