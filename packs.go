package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var packDir = "packs"

func load(pack pack, dataPath string) error {
	local, err := download(http.Get, pack.url)
	if err != nil {
		log.Fatalf("Failed to download [%v] err: %v.\n", pack.url, err)
		return err
	}
	defer os.Remove(local)

	err = unzip(local, mkPath(packDir, pack.name))
	if err != nil {
		log.Fatalf("Failed to unzip archive [%v], err: %v", local, err)
		return err
	}

	if docs[pack.name], err = pack.indexer(); err != nil {
		return err
	}

	data, err := json.Marshal(docs[pack.name])
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(dataPath, data, 0644); err != nil {
		return err
	}

	log.Printf("Installed [%v] entries for pack [%v].", len(docs[pack.name]), pack.name)

	return nil
}

func install(pack pack) error {
	log.Printf("Installing pack [%v].\n", pack.name)

	dataPath := mkPath(packDir, pack.name, "rad-data.json")

	if fileExists(dataPath) {
		log.Printf("Already installed pack [%v], deserializing entries.", pack.name)
		return unmarshalPack(pack, dataPath)
	}

	return load(pack, dataPath)
}
