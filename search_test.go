package rad

import "testing"

func TestFindEntryMissingPackage(t *testing.T) {
	docs = map[string][]entry{}
	_, err := findEntries("scala", "abc")

	if err.Error() != "Package [scala] not installed." {
		t.Errorf("expected error when accessing non existant package, got [%v]", err)
	}
}

func TestFindEntry(t *testing.T) {
	docs = map[string][]entry{
		"scala": []entry{
			entry{
				namespace: []string{"scala", "sys"},
				entity:    "SystemProperties",
				function:  "",
				signature: "",
			},
			entry{
				namespace: []string{"scala", "collection"},
				entity:    "SetProxy",
				function:  "",
				signature: "",
			},
		},
	}
	es, err := findEntries("scala", "SetProxy")

	if err != nil {
		t.Errorf("unexpected error [%v]", err)
		return
	}

	if len(es) != 1 {
		t.Errorf("expected to find one entry but got [%v]", es)
		return
	}

	if es[0].namespace[0] != "scala" ||
		es[0].namespace[1] != "collection" ||
		es[0].entity != "SetProxy" {
		t.Errorf("expected to find SetProxy entry but got [%v]", es[0])
		return
	}

}
