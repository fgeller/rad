package main

import "testing"

func TestFindEntryMissingPackage(t *testing.T) {
	docs = map[string][]entry{}
	_, err := findEntityMember("scala", "abc", "def", 10)

	if err.Error() != "Package [scala] not installed." {
		t.Errorf("expected error when accessing non existant package, got [%v]", err)
	}
}

func TestFindEntry(t *testing.T) {
	docs = map[string][]entry{
		"scala": []entry{
			entry{
				Namespace: []string{"scala", "sys"},
				Entity:    "SystemProperties",
				Member:    "",
				Signature: "",
			},
			entry{
				Namespace: []string{"scala", "collection"},
				Entity:    "SetProxy",
				Member:    "",
				Signature: "",
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
	docs = map[string][]entry{
		"scala": []entry{
			entry{
				Namespace: []string{"scala", "sys"},
				Entity:    "SystemProperties",
				Member:    "",
				Signature: "",
			},
			entry{
				Namespace: []string{"scala", "sys"},
				Entity:    "SystemThings",
				Member:    "",
				Signature: "",
			},
			entry{
				Namespace: []string{"scala", "collection"},
				Entity:    "SetProxy",
				Member:    "",
				Signature: "",
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
				Entity:    "SystemProperties",
				Member:    "hans",
				Signature: "",
			},
			entry{
				Namespace: []string{"scala", "sys"},
				Entity:    "SYSTEMThings",
				Member:    "HANS",
				Signature: "",
			},
			entry{
				Namespace: []string{"scala", "sys"},
				Entity:    "systemThings",
				Member:    "hAnS",
				Signature: "",
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
				Entity:    "HashMap",
				Member:    "clearTable",
				Signature: "():Unit",
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
				Entity:    "HashMap",
				Member:    "clearTable",
				Signature: "():Unit",
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
