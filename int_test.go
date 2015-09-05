package main

import "io/ioutil"
import "os"
import "testing"
import "time"
import "reflect"
import "net/http"
import "encoding/json"

func TestInstallPack(t *testing.T) {

	sampleEntry := entry{[]string{"main"}, "Entity", "Function", "Signature", "Target", "source"}
	indexer := func() ([]entry, error) { return []entry{sampleEntry}, nil }
	serveZip := func(addr string) {
		http.HandleFunc("/test.zip", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./testdata/test.zip")
		})
		http.ListenAndServe(addr, nil)
	}

	// possible race condition here?
	go serveZip(":8881")

	conf := pack{
		name:    "blubb",
		url:     "http://localhost:8881/test.zip",
		indexer: indexer,
	}

	install(conf)

	actual, err := findEntityFunction(conf.name, "Entity", "")
	if err != nil || len(actual) < 1 || !reflect.DeepEqual(sampleEntry, actual[0]) {
		t.Errorf("Expected to find sample entry, got %v [err: %v].\n", actual, err)
	}

	os.RemoveAll("packs/" + conf.name)
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
	}
	docs = map[string][]entry{}
	docs[pack] = samples

	addr := "0.0.0.0:3025"
	go serve(addr)
	time.Sleep(200)

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

	expected, err := json.Marshal(samples[1:])
	if err != nil {
		t.Errorf("unexpected error while marshaling to json: %v", err)
		return
	}

	if string(byts) != string(expected) {
		t.Errorf("unexpected response: %v", string(byts))
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

	expectedEntries := []entry{
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
}
