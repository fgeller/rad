package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

type zipServe struct{}

func (z *zipServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/ping":
		w.Write([]byte("pong"))
	case r.URL.Path == "/test.zip":
		http.ServeFile(w, r, "./testdata/test.zip")
	}

}

func TestInstallPack(t *testing.T) {
	os.RemoveAll("test.zip") // TODO: clean this up, after added conf for packsDir

	e := entry{
		Namespace: []string{"main"},
		Name:      "Entity",
		Members:   []member{{Name: "Member", Signature: "Signature", Target: "Target"}},
	}
	sr := NewSearchResult(e, 0)
	es := []entry{e}
	indexer := func() ([]entry, error) { return es, nil }
	serveZip := func(addr string) { http.ListenAndServe(addr, &zipServe{}) }
	zipAddr := "0.0.0.0:8881"
	conf := pack{
		name:     "blubb",
		location: "http://" + zipAddr + "/test.zip",
		indexer:  indexer,
	}

	defer os.RemoveAll("packs/" + conf.name)
	defer os.RemoveAll("test.zip")

	go serveZip(zipAddr)
	awaitPing(zipAddr)

	install(conf)

	actual, err := findEntityMember(conf.name, "Entity", "", 10)
	if err != nil {
		t.Errorf("unexpected error while finding entries: %v", err)
		return
	}

	if len(actual) != 1 {
		t.Errorf("expected to find 1 entry, got: %v", len(actual))
		return
	}

	if !sr.eq(actual[0]) {
		t.Errorf("Expected to find sample entry, got \n%v\nbut expected\n%v", actual, es)
	}
}

func TestInstallLocalPack(t *testing.T) {
	e := entry{
		Namespace: []string{"com", "example"},
		Name:      "Entity",
		Members:   []member{{Name: "Member", Signature: "Signature", Target: "Target"}},
	}
	sr := NewSearchResult(e, 0)
	es := []entry{e}
	indexer := func() ([]entry, error) { return es, nil }
	p := pack{
		name:     "testpack",
		indexer:  indexer,
		location: "testdata/test.zip",
	}
	docs = map[string][]entry{}
	os.RemoveAll("packs/" + p.name)
	defer os.RemoveAll("packs/" + p.name)

	install(p)

	found, err := findEntityMember(p.name, e.Name, e.Members[0].Name, 10)
	if err != nil {
		t.Errorf("unexpected error when trying to find entries: %v\n", err)
		return
	}

	if len(found) != 1 {
		t.Errorf("expected to find single test entry [%v, %v, %v] got\n%v\ndocs %v\n", p.name, e.Name, e.Members[0].Name, found, docs)
		return
	}

	if !found[0].eq(sr) {
		t.Errorf("expected to find test entry\n%v\ngot\n%v\n", e, found[0])
		return
	}

}

func TestInstallExistingSerializedPack(t *testing.T) {

	e := entry{
		Namespace: []string{"main"},
		Name:      "Entity",
		Members:   []member{{Name: "Member", Signature: "Signature", Target: "Target"}},
	}
	sr := NewSearchResult(e, 0)
	es := []entry{e}
	indexer := func() ([]entry, error) { return []entry{}, nil }
	conf := pack{
		name:     "blubb",
		location: "http://localhost:8881/test.zip",
		indexer:  indexer,
	}
	os.RemoveAll("packs/" + conf.name)

	err := os.MkdirAll("packs/"+conf.name, 0755)
	if err != nil {
		t.Errorf("unexpected error when creating dir: %v", err)
		return
	}
	data, err := json.Marshal(es)
	if err != nil {
		t.Errorf("unexpected error when serializing data: %v", err)
		return
	}
	dataPath := "packs/" + conf.name + "/rad-data.json"
	err = ioutil.WriteFile(dataPath, data, 0644)
	if err != nil {
		t.Errorf("unexpected error when writing serialized data: %v", err)
		return
	}

	install(conf)

	actual, err := findEntityMember(conf.name, "Entity", "", 10)
	if err != nil {
		t.Errorf("unexpected error while finding entries: %v", err)
		return
	}

	if len(actual) != 1 {
		t.Errorf("expected to find 1 entry, got: %v", len(actual))
		return
	}

	if !sr.eq(actual[0]) {
		t.Errorf("Expected to find sample entry, got \n%v\nbut expected\n%v", actual, es)
	}
}

func TestFindEntityMembers(t *testing.T) {
	pack := "test"
	entity := "Entity"
	samples := []entry{
		{Namespace: []string{"main"}, Name: "AnotherEntity", Members: []member{{Name: "abc", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: entity, Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: entity, Members: []member{{Name: "ab", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: entity, Members: []member{{Name: "abc", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: entity, Members: []member{{Name: "d", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: entity + "suffix", Members: []member{{Name: "x", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: entity + "suffix", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: "_zzz", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: "_zyz", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: "_zzy", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: "_yzz", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: "_zyy", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: "_yyz", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: "_yyy", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: "_yzy", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: "_yxz", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: "_xyz", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: "_yyx", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: "_xxz", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: "_yxz", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}},
		{Namespace: []string{"main"}, Name: "_xxx", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}},
	}
	docs = map[string][]entry{}
	docs[pack] = samples

	addr := "0.0.0.0:3025"
	go serve(addr)
	awaitPing(addr)

	// find all for entity
	res, err := http.Get("http://" + addr + "/s?p=" + pack + "&e=" + entity)
	if err != nil {
		t.Errorf("unexpected error while finding entries: %v", err)
		return
	}

	byts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unexpected error while reading response body: %v", err)
		return
	}

	expectedEntries := []searchResult{
		NewSearchResult(entry{Namespace: []string{"main"}, Name: entity, Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}}, 0),
		NewSearchResult(entry{Namespace: []string{"main"}, Name: entity, Members: []member{{Name: "ab", Signature: "Signature", Target: "Target"}}}, 0),
		NewSearchResult(entry{Namespace: []string{"main"}, Name: entity, Members: []member{{Name: "abc", Signature: "Signature", Target: "Target"}}}, 0),
		NewSearchResult(entry{Namespace: []string{"main"}, Name: entity, Members: []member{{Name: "d", Signature: "Signature", Target: "Target"}}}, 0),
		NewSearchResult(entry{Namespace: []string{"main"}, Name: entity + "suffix", Members: []member{{Name: "x", Signature: "Signature", Target: "Target"}}}, 0),
		NewSearchResult(entry{Namespace: []string{"main"}, Name: entity + "suffix", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}}, 0),
	}
	expected, err := json.Marshal(expectedEntries)
	if err != nil {
		t.Errorf("unexpected error while marshaling to json: %v", err)
		return
	}

	if string(byts) != string(expected) {
		t.Errorf("unexpected response, got\n%v\nbut expected\n%v\n", string(byts), string(expected))
	}

	// find all with given member prefix

	res, err = http.Get("http://" + addr + "/s?p=" + pack + "&e=" + entity + "&m=a")
	if err != nil {
		t.Errorf("unexpected error while finding entries: %v", err)
		return
	}

	byts, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unexpected error while reading response body: %v", err)
		return
	}

	expectedEntries = []searchResult{
		NewSearchResult(entry{Namespace: []string{"main"}, Name: entity, Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}}, 0),
		NewSearchResult(entry{Namespace: []string{"main"}, Name: entity, Members: []member{{Name: "ab", Signature: "Signature", Target: "Target"}}}, 0),
		NewSearchResult(entry{Namespace: []string{"main"}, Name: entity, Members: []member{{Name: "abc", Signature: "Signature", Target: "Target"}}}, 0),
		NewSearchResult(entry{Namespace: []string{"main"}, Name: entity + "suffix", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}}, 0),
	}
	expected, err = json.Marshal(expectedEntries)
	if err != nil {
		t.Errorf("unexpected error while marshaling to json: %v", err)
		return
	}

	if string(byts) != string(expected) {
		t.Errorf("unexpected response, got \n%v\nbut expected\n%v\n", string(byts), string(expected))
	}

	// limit the results by default to 10
	res, err = http.Get("http://" + addr + "/s?p=" + pack + "&e=_")
	if err != nil {
		t.Errorf("unexpected error while finding entries: %v", err)
		return
	}

	byts, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unexpected error while reading response body: %v", err)
		return
	}

	expectedEntries = []searchResult{
		NewSearchResult(entry{Namespace: []string{"main"}, Name: "_zzz", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}}, 0),
		NewSearchResult(entry{Namespace: []string{"main"}, Name: "_zyz", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}}, 0),
		NewSearchResult(entry{Namespace: []string{"main"}, Name: "_zzy", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}}, 0),
		NewSearchResult(entry{Namespace: []string{"main"}, Name: "_yzz", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}}, 0),
		NewSearchResult(entry{Namespace: []string{"main"}, Name: "_zyy", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}}, 0),
		NewSearchResult(entry{Namespace: []string{"main"}, Name: "_yyz", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}}, 0),
		NewSearchResult(entry{Namespace: []string{"main"}, Name: "_yyy", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}}, 0),
		NewSearchResult(entry{Namespace: []string{"main"}, Name: "_yzy", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}}, 0),
		NewSearchResult(entry{Namespace: []string{"main"}, Name: "_yxz", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}}, 0),
		NewSearchResult(entry{Namespace: []string{"main"}, Name: "_xyz", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}}, 0),
	}
	expected, err = json.Marshal(expectedEntries)
	if err != nil {
		t.Errorf("unexpected error while marshaling to json: %v", err)
		return
	}

	if string(byts) != string(expected) {
		t.Errorf("unexpected response, got \n%v\nbut expected\n%v\n", string(byts), string(expected))
	}

	// allow custom limit
	limit := 2
	res, err = http.Get(fmt.Sprintf("http://%v/s?p=%v&e=%v&limit=%v", addr, pack, "_", limit))
	if err != nil {
		t.Errorf("unexpected error while finding entries: %v", err)
		return
	}

	byts, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unexpected error while reading response body: %v", err)
		return
	}

	expectedEntries = []searchResult{
		NewSearchResult(entry{Namespace: []string{"main"}, Name: "_zzz", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}}, 0),
		NewSearchResult(entry{Namespace: []string{"main"}, Name: "_zyz", Members: []member{{Name: "a", Signature: "Signature", Target: "Target"}}}, 0),
	}
	expected, err = json.Marshal(expectedEntries)
	if err != nil {
		t.Errorf("unexpected error while marshaling to json: %v", err)
		return
	}

	if string(byts) != string(expected) {
		t.Errorf("unexpected response, got \n%v\nbut expected\n%v\n", string(byts), string(expected))
	}
}

func awaitPing(addr string) {
	for i := 0; i < 10; i++ {
		_, err := http.Get("http://" + addr + "/ping")
		if err == nil {
			return
		}
		time.Sleep(100)
	}
}
