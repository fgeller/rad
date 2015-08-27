package main

import "os"
import "testing"
import "reflect"
import "net/http"

func TestInstallPack(t *testing.T) {

	sampleEntry := entry{[]string{"main"}, "Entity", "Function", "Signature", "Target", "source"}
	indexer := func() []entry { return []entry{sampleEntry} }
	serveZip := func(addr string) {
		http.HandleFunc("/test.zip", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./testdata/test.zip")
		})
		http.ListenAndServe(addr, nil)
	}

	// possible race condition here?
	go serveZip(":8881")

	conf := packConfig{
		name:    "blubb",
		url:     "http://localhost:8881/test.zip",
		indexer: indexer,
	}

	install(conf)

	actual, err := findEntries(conf.name, "Entity")
	if err != nil || len(actual) < 1 || !reflect.DeepEqual(sampleEntry, actual[0]) {
		t.Errorf("Expected to find sample entry, got %v [err: %v].\n", actual, err)
	}

	os.RemoveAll("packs/" + conf.name)
}
