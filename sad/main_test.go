package main

import (
	"../shared"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func setup() string {
	docs = map[string][]shared.Namespace{}
	tmp, err := ioutil.TempDir("", "sad-main-test-pack-dir")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
	}
	packDir = tmp
	return tmp
}

func TestInstallingLocalPack(t *testing.T) {
	defer os.RemoveAll(setup())

	pp := "testdata/jdk.zip"

	err := install(pp)
	if err != nil {
		t.Errorf("Expected successful installing of local pack %v, got err: %v", pp, err)
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

func populatePackDir() (map[string][]shared.Namespace, map[string]shared.Pack) {
	p1 := shared.Pack{Name: "p1", Type: "java"}
	p1Data := []shared.Namespace{
		{Path: []string{"A"}, Members: []shared.Member{{Name: "M1", Target: "T1"}}},
	}

	p2 := shared.Pack{Name: "p2", Type: "go"}
	p2Data := []shared.Namespace{
		{Path: []string{"B"}, Members: []shared.Member{{Name: "M2", Target: "T2"}}},
	}

	packs := map[string]shared.Pack{p1.Name: p1, p2.Name: p2}
	data := map[string][]shared.Namespace{
		p1.Name: p1Data,
		p2.Name: p2Data,
	}

	for _, p := range packs {

		// make dir
		err := os.MkdirAll(filepath.Join(packDir, p.Name), 0755)
		if err != nil {
			log.Fatalf("Failed to create pack %v dir: %v", p.Name, err)
		}

		// create conf file
		cd, err := json.MarshalIndent(p, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal conf for pack %v: %v", p.Name, err)
		}
		cp := filepath.Join(packDir, p.Name, "pack.json")
		err = ioutil.WriteFile(cp, cd, 0600)
		if err != nil {
			log.Fatalf("Failed to write conf file for pack %v: %v", p.Name, err)
		}

		// create data file
		dd, err := json.MarshalIndent(data[p.Name], "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal data for pack %v: %v", p.Name, err)
		}
		dp := filepath.Join(packDir, p.Name, "data.json")
		err = ioutil.WriteFile(dp, dd, 0600)
		if err != nil {
			log.Fatalf("Failed to write data file for pack %v: %v", p.Name, err)
		}
	}

	return data, packs
}

func TestLoadInstalledPack(t *testing.T) {
	defer os.RemoveAll(setup())
	expectedDocs, expectedPacks := populatePackDir()

	err := loadInstalled()
	if err != nil {
		t.Errorf("Expected successful loading of installed pack %v, got err: %v", packDir, err)
		return
	}

	if !reflect.DeepEqual(expectedDocs, docs) {
		t.Errorf("Expected docs:\n%v\nBut got:\n%v\n", expectedDocs, docs)
		return
	}

	if !reflect.DeepEqual(expectedPacks, packs) {
		t.Errorf("Expected packs:\n%v\nBut got:\n%v\n", expectedPacks, packs)
		return
	}
}
