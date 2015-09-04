package main

import "testing"

func TestFindEntryMissingPackage(t *testing.T) {
	docs = map[string][]entry{}
	_, err := findEntityFunction("scala", "abc", "def")

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
				Function:  "",
				Signature: "",
			},
			entry{
				Namespace: []string{"scala", "collection"},
				Entity:    "SetProxy",
				Function:  "",
				Signature: "",
			},
		},
	}
	es, err := findEntityFunction("scala", "SetProxy", "")

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

func TestFindFunction(t *testing.T) {
	docs = map[string][]entry{
		"scala": []entry{
			entry{
				Namespace: []string{"scala", "collection", "mutable"},
				Entity:    "HashMap",
				Function:  "clearTable",
				Signature: "():Unit",
			},
		},
	}

	es, err := findEntityFunction("scala", "HashMap", "clearTable")

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

func TestFindFunctionByPrefix(t *testing.T) {
	docs = map[string][]entry{
		"scala": []entry{
			entry{
				Namespace: []string{"scala", "collection", "mutable"},
				Entity:    "HashMap",
				Function:  "clearTable",
				Signature: "():Unit",
			},
		},
	}

	es, err := findEntityFunction("scala", "HashMap", "clear")

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
