package main

import (
	"../shared"
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

func load(path string) error {

	log.Printf("Loading pack from path %s\n", path)

	// unzip it
	err := shared.Unzip(path, packDir)
	if err != nil {
		return err
	}

	fn := filepath.Base(path)                  // jdk.zip
	packName := fn[:strings.Index(fn, ".zip")] // jdk

	// load data
	dp := filepath.Join(packDir, packName, "data.json")
	db, err := ioutil.ReadFile(dp)
	if err != nil {
		return err
	}
	var entries []shared.Entry
	err = json.Unmarshal(db, &entries)
	if err != nil {
		return err
	}

	// add to global var
	docs[packName] = entries

	log.Printf("Found %v entries for pack %v\n", len(entries), packName)

	return nil
}
