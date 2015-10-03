package main

import (
	"../shared"
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

func loadInstalled() error {
	log.Printf("Loading installed packs from %v\n", packDir)

	dirs, err := ioutil.ReadDir(packDir)
	if err != nil {
		log.Printf("Failed to read contents of packDir %v: %v\n", packDir, err)
		return err
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			log.Printf("Skipping: Expected only directories, found %v.", dir.Name())
		}

		pack := dir.Name()
		log.Printf("Loading pack %v\n", pack)
		df := filepath.Join(packDir, pack, "data.json")
		dc, err := ioutil.ReadFile(df)
		if err != nil {
			log.Printf("Skipping: Could not load data for %v (err: %v).", pack, err)
		}
		var data []shared.Namespace
		err = json.Unmarshal(dc, &data)
		docs[pack] = data
		log.Printf("Found %v entries for %v.", len(data), pack)
	}

	return nil
}

func install(path string) error {

	log.Printf("Installing %v\n", path)

	// unzip it
	err := shared.Unzip(path, packDir)
	if err != nil {
		return err
	}

	fn := filepath.Base(path)                  // jdk.zip
	packName := fn[:strings.Index(fn, ".zip")] // jdk

	// load data
	dp := filepath.Join(packDir, packName, "data.json")
	log.Printf("Reading data from %v\n", dp)
	db, err := ioutil.ReadFile(dp)
	if err != nil {
		return err
	}
	log.Printf("Unmarshalling data.\n")
	var namespaces []shared.Namespace
	err = json.Unmarshal(db, &namespaces)
	if err != nil {
		return err
	}

	// add to global var
	docs[packName] = namespaces
	log.Printf("Found %v entries for pack %v\n", len(namespaces), packName)

	return nil
}
