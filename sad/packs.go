package main

import (
	"../shared"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func loadInstalled() error {
	log.Printf("Loading installed packs from %v\n", config.packDir)

	err := os.MkdirAll(config.packDir, 0755)
	if err != nil {
		return err
	}

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

func remove(pack string) error {
	return os.RemoveAll(filepath.Join(config.packDir, pack))
}

func install(path string) error {

	log.Printf("Installing %v\n", path)

	tmp, err := ioutil.TempDir("", "unzipped")
	if err != nil {
		return err
	}

	if err = shared.Unzip(path, tmp); err != nil {
		log.Printf("Failed to unzip %v: %v\n", path, err)
		return err
	}

	fs, err := ioutil.ReadDir(tmp)
	if err != nil {
		log.Printf("Failed to read directory contents: %v\n", err)
		return err
	}

	if len(fs) != 1 {
		return fmt.Errorf("Expected one file in pack directory, got: %v", len(fs))
	}

	packName := fs[0].Name()

	log.Printf("Copying contents for [%v] into pack dir.\n", packName)
	_, err = shared.CopyDir(
		filepath.Join(tmp, packName),
		filepath.Join(config.packDir),
	)

	if err != nil {
		log.Printf("Failed to copy directory into packdir: %v\n", err)
		return err
	}

	return loadInstalled() // TODO overkill?
}
