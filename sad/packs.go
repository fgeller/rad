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
	Load
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
	err error
	pck shared.Pack
	nss []shared.Namespace
}

func loadRemotePack(fn string) (shared.Pack, []shared.Namespace, error) {
	// TODO: review for cleanup after failure
	var pck shared.Pack
	var nss []shared.Namespace

	path, err := shared.DownloadToTemp("http://" + config.sapAddr + "/pack/" + fn)
	if err != nil {
		log.Printf("Error downloading pack: %v\n", err)
		return pck, nss, err
	}
	defer os.RemoveAll(path)

	tmp, err := ioutil.TempDir("", "unzipped")
	if err != nil {
		return pck, nss, err
	}

	if err = shared.Unzip(path, tmp); err != nil {
		log.Printf("Failed to unzip %v: %v\n", path, err)
		return pck, nss, err
	}

	fs, err := ioutil.ReadDir(tmp)
	if err != nil {
		log.Printf("Failed to read directory contents: %v\n", err)
		return pck, nss, err
	}

	if len(fs) != 1 || !fs[0].IsDir() {
		return pck, nss, fmt.Errorf("Expected one directory in pack directory, got: %v", fs)
	}

	pn := fs[0].Name()

	log.Printf("Copying contents for [%v] into pack dir.\n", pn)
	if _, err = shared.CopyDir(filepath.Join(tmp, pn), config.packDir); err != nil {
		log.Printf("Failed to copy directory into packdir: %v\n", err)
		return pck, nss, err
	}

	pf := filepath.Join(config.packDir, pn, "pack.json")
	pc, err := ioutil.ReadFile(pf)
	if err != nil {
		log.Printf("Could not load pack info for %v (err: %v).", pn, err)
		return pck, nss, err
	}

	err = json.Unmarshal(pc, &pck)
	log.Printf("Found info %v for %v.", pck, pn)

	df := filepath.Join(config.packDir, pn, "data.json")
	dc, err := ioutil.ReadFile(df)
	if err != nil {
		log.Printf("Skipping: Could not load data for %v (err: %v).", pn, err)
		return pck, nss, err
	}

	err = json.Unmarshal(dc, &nss)
	log.Printf("Found %v entries for %v.", len(nss), pn)

	return pck, nss, err
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
		pck, nss, err := loadRemotePack(req.pck.File)
		if err != nil {
			req.res <- packResp{err: err}
			return
		}

		packs[pck.Name] = pck
		docs[pck.Name] = nss
	}

	loadPack := func(req packReq) {
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
			case req.tpe == Load:
				loadPack(req)
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
				req.res <- packResp{err: fmt.Errorf("Unsupported req.tpe %+vMaster\n")}
				close(req.res)
			}
		}
	}
}

func loadInstalled() error {
	log.Printf("Loading installed packs from %v\n", config.packDir)

	if err := os.MkdirAll(config.packDir, 0755); err != nil {
		return err
	}

	dirs, err := ioutil.ReadDir(config.packDir)
	if err != nil {
		log.Printf("Failed to read contents of packDir %v: %v\n", config.packDir, err)
		return err
	}

loadingdirs:
	for _, dir := range dirs {
		if !dir.IsDir() {
			log.Printf("Skipping: Expected only directories, found %v.", dir.Name())
			continue loadingdirs
		}

		pack := dir.Name()
		log.Printf("Loading pack %v\n", pack)

		if err := loadFromPackDir(pack); err != nil {
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

func loadFromPackDir(name string) error {

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
	req := packReq{Load, pck, nss, res}
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

	return loadFromPackDir(pn)
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

func loadPack(pck shared.Pack, nss []shared.Namespace) {
	res := make(chan packResp)
	req := packReq{Load, pck, nss, res}
	log.Printf("Load pack: %v\n", req.pck.Name)
	global.packs <- req
	<-res
}
