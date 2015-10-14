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

type reqType int

const (
	Install reqType = iota
	Remove
	Read
	Reset
)

type packReq struct {
	tpe reqType
	pck shared.Pack
	nss []shared.Namespace
	res chan packResp
}

type packResp struct {
	pck shared.Pack
	nss []shared.Namespace
}

func packMaster() {
	log.Printf("Started packMaster\n")

	docs := map[string][]shared.Namespace{}
	packs := map[string]shared.Pack{}

	sendPacks := func(req packReq) {
		for n, p := range packs {
			req.res <- packResp{
				pck: p,
				nss: docs[n],
			}
		}
	}

	installPack := func(req packReq) {
		packs[req.pck.Name] = req.pck
		docs[req.pck.Name] = req.nss
	}

	removePack := func(req packReq) {
		if _, ok := packs[req.pck.Name]; !ok {
			return // pack doesn't exist (anymore?)
		}
		delete(packs, req.pck.Name)
		delete(docs, req.pck.Name)
		os.RemoveAll(filepath.Join(config.packDir, req.pck.Name)) // TODO log
	}

	reset := func(req packReq) {
		docs = map[string][]shared.Namespace{}
		packs = map[string]shared.Pack{}
	}

	for {
		select {
		case req := <-global.packs:
			switch {
			case req.tpe == Install:
				installPack(req)
				close(req.res)
			case req.tpe == Remove:
				removePack(req)
				close(req.res)
			case req.tpe == Read:
				sendPacks(req)
				close(req.res)
			case req.tpe == Reset:
				reset(req)
				close(req.res)
			default:
				log.Printf("Unsupported req.tpe %+v in packMaster\n")
			}
		}
	}
}

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

		if err := loadPack(pack); err != nil {
			log.Printf("Failed to install pack %v: %v", pack, err)
			return err
		}
	}

	return nil
}

func remove(pn string) {
	res := make(chan packResp)
	pck := shared.Pack{Name: pn}
	req := packReq{tpe: Remove, pck: pck, res: res}
	global.packs <- req
	_, ok := <-res
	log.Printf("Removed pack %+v (ok: %v)\n", pck, ok)
}

func loadPack(name string) error {

	pf := filepath.Join(config.packDir, name, "pack.json")
	pc, err := ioutil.ReadFile(pf)
	if err != nil {
		log.Printf("Skipping: Could not load pack info for %v (err: %v).", name, err)
	}
	var pck shared.Pack
	err = json.Unmarshal(pc, &pck)
	log.Printf("Found info %v for %v.", pck, name)

	df := filepath.Join(config.packDir, name, "data.json")
	dc, err := ioutil.ReadFile(df)
	if err != nil {
		log.Printf("Skipping: Could not load data for %v (err: %v).", name, err)
	}
	var nss []shared.Namespace
	err = json.Unmarshal(dc, &nss)
	log.Printf("Found %v entries for %v.", len(nss), name)

	res := make(chan packResp)
	req := packReq{Install, pck, nss, res}
	global.packs <- req
	_, ok := <-res
	log.Printf("Successfully installed pack %+v (ok: %v)\n", pck, ok)

	return nil
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

	pn := fs[0].Name()

	log.Printf("Copying contents for [%v] into pack dir.\n", pn)
	_, err = shared.CopyDir(
		filepath.Join(tmp, pn),
		filepath.Join(config.packDir),
	)

	if err != nil {
		log.Printf("Failed to copy directory into packdir: %v\n", err)
		return err
	}

	return loadPack(pn)
}

func installedPacks() []shared.Pack {
	res := make(chan packResp)
	req := packReq{tpe: Read, res: res}
	global.packs <- req
	pcks := []shared.Pack{}
	for {
		select {
		case resp, ok := <-res:
			if !ok {
				return pcks
			}
			pcks = append(pcks, resp.pck)
		}
	}
}

func installedDocs() map[string][]shared.Namespace {
	res := make(chan packResp)
	req := packReq{tpe: Read, res: res}
	global.packs <- req
	docs := map[string][]shared.Namespace{}
	for {
		select {
		case resp, ok := <-res:
			if !ok {
				return docs
			}
			docs[resp.pck.Name] = resp.nss
		}
	}
}

func installPack(pck shared.Pack, nss []shared.Namespace) {
	res := make(chan packResp)
	req := packReq{Install, pck, nss, res}
	log.Printf("Installing pack: %v\n", req.pck.Name)
	global.packs <- req
	<-res
}
