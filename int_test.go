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
	e := entry{
		Namespace: []string{"main"},
		Entity:    "Entity",
		Function:  "Function",
		Signature: "Signature",
		Target:    "Target",
		Source:    "source",
	}
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

	actual, err := findEntityFunction(conf.name, "Entity", "", 10)
	if err != nil {
		t.Errorf("unexpected error while finding entries: %v", err)
		return
	}

	if len(actual) != 1 {
		t.Errorf("expected to find 1 entry, got: %v", len(actual))
		return
	}

	if !e.eq(actual[0]) {
		t.Errorf("Expected to find sample entry, got \n%v\nbut expected\n%v", actual, es)
	}
}

func TestInstallExistingSerializedPack(t *testing.T) {

	e := entry{
		Namespace: []string{"main"},
		Entity:    "Entity",
		Function:  "Function",
		Signature: "Signature",
		Target:    "Target",
		Source:    "source",
	}
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

	actual, err := findEntityFunction(conf.name, "Entity", "", 10)
	if err != nil {
		t.Errorf("unexpected error while finding entries: %v", err)
		return
	}

	if len(actual) != 1 {
		t.Errorf("expected to find 1 entry, got: %v", len(actual))
		return
	}

	if !e.eq(actual[0]) {
		t.Errorf("Expected to find sample entry, got \n%v\nbut expected\n%v", actual, es)
	}
}

func TestFindEntityFunctions(t *testing.T) {
	pack := "test"
	entity := "Entity"
	samples := []entry{
		{[]string{"main"}, "AnotherEntity", "abc", "Signature", "Target", "source"},
		{[]string{"main"}, entity, "a", "Signature", "Target", "source"},
		{[]string{"main"}, entity, "ab", "Signature", "Target", "source"},
		{[]string{"main"}, entity, "abc", "Signature", "Target", "source"},
		{[]string{"main"}, entity, "d", "Signature", "Target", "source"},
		{[]string{"main"}, entity + "suffix", "x", "Signature", "Target", "source"},
		{[]string{"main"}, entity + "suffix", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_zzz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_zyz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_zzy", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yzz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_zyy", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yyz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yyy", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yzy", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yxz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_xyz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yyx", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_xxz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yxz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_xxx", "a", "Signature", "Target", "source"},
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

	expectedEntries := []entry{
		{[]string{"main"}, entity, "a", "Signature", "Target", "source"},
		{[]string{"main"}, entity, "ab", "Signature", "Target", "source"},
		{[]string{"main"}, entity, "abc", "Signature", "Target", "source"},
		{[]string{"main"}, entity, "d", "Signature", "Target", "source"},
		{[]string{"main"}, entity + "suffix", "x", "Signature", "Target", "source"},
		{[]string{"main"}, entity + "suffix", "a", "Signature", "Target", "source"},
	}
	expected, err := json.Marshal(expectedEntries)
	if err != nil {
		t.Errorf("unexpected error while marshaling to json: %v", err)
		return
	}

	if string(byts) != string(expected) {
		t.Errorf("unexpected response, got\n%v\nbut expected\n%v\n", string(byts), string(expected))
	}

	// find all with given function prefix

	res, err = http.Get("http://" + addr + "/s?p=" + pack + "&e=" + entity + "&f=a")
	if err != nil {
		t.Errorf("unexpected error while finding entries: %v", err)
		return
	}

	byts, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unexpected error while reading response body: %v", err)
		return
	}

	expectedEntries = []entry{
		{[]string{"main"}, entity, "a", "Signature", "Target", "source"},
		{[]string{"main"}, entity, "ab", "Signature", "Target", "source"},
		{[]string{"main"}, entity, "abc", "Signature", "Target", "source"},
		{[]string{"main"}, entity + "suffix", "a", "Signature", "Target", "source"},
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

	expectedEntries = []entry{
		{[]string{"main"}, "_zzz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_zyz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_zzy", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yzz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_zyy", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yyz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yyy", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yzy", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yxz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_xyz", "a", "Signature", "Target", "source"},
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

	expectedEntries = []entry{
		{[]string{"main"}, "_zzz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_zyz", "a", "Signature", "Target", "source"},
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
