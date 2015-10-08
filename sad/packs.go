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
	log.Printf("Loading installed packs from %v\n", config.packDir)

	dirs, err := ioutil.ReadDir(config.packDir)
	if err != nil {
		log.Printf("Failed to read contents of packDir %v: %v\n", config.packDir, err)
		return err
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			log.Printf("Skipping: Expected only directories, found %v.", dir.Name())
		}

		pack := dir.Name()
		log.Printf("Loading pack %v\n", pack)

		pf := filepath.Join(config.packDir, pack, "pack.json")
		pc, err := ioutil.ReadFile(pf)
		if err != nil {
			log.Printf("Skipping: Could not load pack info for %v (err: %v).", pack, err)
		}
		var packInfo shared.Pack
		err = json.Unmarshal(pc, &packInfo)
		global.packs[pack] = packInfo
		log.Printf("Found info %v for %v.", packInfo, pack)

		df := filepath.Join(config.packDir, pack, "data.json")
		dc, err := ioutil.ReadFile(df)
		if err != nil {
			log.Printf("Skipping: Could not load data for %v (err: %v).", pack, err)
		}
		var data []shared.Namespace
		err = json.Unmarshal(dc, &data)
		global.docs[pack] = data
		log.Printf("Found %v entries for %v.", len(data), pack)
	}

	return nil
}

func install(path string) error {

	log.Printf("Installing %v\n", path)

	// unzip it
	err := shared.Unzip(path, config.packDir)
	if err != nil {
		return err
	}

	fn := filepath.Base(path)                  // jdk.zip
	packName := fn[:strings.Index(fn, ".zip")] // jdk

	// load data
	dp := filepath.Join(config.packDir, packName, "data.json")
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
	global.docs[packName] = namespaces
	log.Printf("Found %v entries for pack %v\n", len(namespaces), packName)

	return nil
}
