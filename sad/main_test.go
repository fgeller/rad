package main

import (
	"../shared"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func setup() string {
	docs = map[string][]shared.Entry{}
	tmp, err := ioutil.TempDir("", "sad-main-test-pack-dir")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
	}
	packDir = tmp
	return tmp
}

func TestLoadingLocalPack(t *testing.T) {
	defer os.RemoveAll(setup())

	pp := "testdata/jdk.zip"

	err := load(pp)
	if err != nil {
		t.Errorf("Expected successful loading of local pack %v, got err: %v", pp, err)
		return
	}

	entries, ok := docs["jdk"]
	if !ok {
		t.Errorf("Could not access entries in docs map %v", docs)
		return
	}

	if len(entries) < 1 {
		t.Errorf("Found no entries in docs map %v", docs)
	}

}
