package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func load(pack pack, localPath string) error {
	log.Printf("Loading pack [%v] for [%v].\n", pack.name, localPath)

	err := unzip(localPath, mkPath(packDir, pack.name))
	if err != nil {
		log.Fatalf("Failed to unzip archive [%v], err: %v", localPath, err)
		return err
	}

	if docs[pack.name], err = pack.indexer(); err != nil {
		return err
	}

	data, err := json.Marshal(docs[pack.name])
	if err != nil {
		return err
	}

	dataPath := mkPath(packDir, pack.name, "rad-data.json")
	if err = ioutil.WriteFile(dataPath, data, 0644); err != nil {
		return err
	}

	log.Printf("Installed [%v] entries for pack [%v].", len(docs[pack.name]), pack.name)
	return nil
}

func fetchRemote(pack pack) (string, error) {
	log.Printf("Loading remote pack [%v].\n", pack.name)

	local, err := download(http.Get, pack.location)
	if err != nil {
		log.Fatalf("Failed to download [%v] err: %v.\n", pack.location, err)
		return local, err
	}

	return local, nil
}

func install(pack pack) error {
	log.Printf("Installing pack [%v].\n", pack.name)

	dataPath := mkPath(packDir, pack.name, "rad-data.json")

	if fileExists(dataPath) {
		log.Printf("Already installed pack [%v], deserializing entries.", pack.name)
		return unmarshalPack(pack, dataPath)
	}

	local := pack.location
	if u, err := url.Parse(pack.location); err == nil && u.Scheme != "" {
		local, err = fetchRemote(pack)
		if err != nil {
			return err
		}
	}

	return load(pack, local)
}
