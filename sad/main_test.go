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
	"time"
)

func populatePackDir() (map[string][]shared.Namespace, []shared.Pack) {

	p1 := shared.Pack{
		Name:      "p1",
		Type:      "java",
		Created:   time.Now(),
		NameCount: 1,
	}
	p1Data := []shared.Namespace{
		{Path: "A", Members: []shared.Member{{Name: "M1", Target: "T1"}}},
	}

	p2 := shared.Pack{
		Name:      "p2",
		Type:      "go",
		Created:   time.Now(),
		NameCount: 1,
	}
	p2Data := []shared.Namespace{
		{Path: "B", Members: []shared.Member{{Name: "M2", Target: "T2"}}},
	}

	ps := map[string]shared.Pack{p1.Name: p1, p2.Name: p2}
	data := map[string][]shared.Namespace{
		p1.Name: p1Data,
		p2.Name: p2Data,
	}

	for _, p := range ps {
		// make dir
		err := os.MkdirAll(filepath.Join(config.packDir, p.Name), 0755)
		if err != nil {
			log.Fatalf("Failed to create pack %v dir: %v", p.Name, err)
		}

		// create conf file
		cd, err := json.MarshalIndent(p, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal conf for pack %v: %v", p.Name, err)
		}
		cp := filepath.Join(config.packDir, p.Name, "pack.json")
		err = ioutil.WriteFile(cp, cd, 0600)
		if err != nil {
			log.Fatalf("Failed to write conf file for pack %v: %v", p.Name, err)
		}

		// create data file
		dd, err := json.MarshalIndent(data[p.Name], "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal data for pack %v: %v", p.Name, err)
		}
		dp := filepath.Join(config.packDir, p.Name, "data.json")
		err = ioutil.WriteFile(dp, dd, 0600)
		if err != nil {
			log.Fatalf("Failed to write data file for pack %v: %v", p.Name, err)
		}
	}

	return data, []shared.Pack{p2, p1}
}

func TestFindHomeDir(t *testing.T) {

	hd, err := findHomeDir()
	if err != nil {
		t.Errorf("Error finding home dir: %v", err)
		return
	}

	if len(hd) <= 0 {
		t.Errorf("Error finding home dir, empty dir.")
		return
	}
}

func TestLoadInstalledPack(t *testing.T) {
	defer os.RemoveAll(setup())
	expectedDocs, expectedPacks := populatePackDir()

	err := loadInstalled()
	if err != nil {
		t.Errorf("Expected successful loading of installed pack %v, got err: %v", config.packDir, err)
		return
	}

	pcks := installedPacks()
	docs := installedDocs()

	if !reflect.DeepEqual(expectedDocs, docs) {
		t.Errorf("Expected docs:\n%v\nBut got:\n%v\n", expectedDocs, docs)
		return
	}

comparing:
	for _, p := range expectedPacks {
		for _, ap := range pcks {
			if reflect.DeepEqual(ap, p) {
				continue comparing
			}
		}
		t.Errorf("Expected %v to be installed, got %v.", p, pcks)
		return
	}
}
