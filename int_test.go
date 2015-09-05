package main

import "io/ioutil"
import "os"
import "testing"
import "time"
import "reflect"
import "net/http"

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

func TestFindEntries(t *testing.T) {
	pack := "test"
	entity := "Entity"
	fun := "Function"
	sample := entry{[]string{"main"}, entity, fun, "Signature", "Target", "source"}
	docs = map[string][]entry{}
	docs[pack] = []entry{sample}

	addr := "0.0.0.0:3025"
	go serve(addr)
	time.Sleep(200)

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

	expected := `[{"Namespace":["main"],"Entity":"Entity","Function":"Function","Signature":"Signature","Target":"Target"}]`
	bdy := string(byts)
	if bdy != expected {
		t.Errorf("unexpected response: %v", bdy)
	}
}
