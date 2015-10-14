package main

import (
	"../shared"

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

var global struct {
	packs  chan packReq
	master bool // should make this a mutex
	assets map[string]shared.Asset
}

func setupGlobals() {
	global.packs = make(chan packReq)
	global.assets = map[string]shared.Asset{}
	go packMaster()
}

func resetGlobals() {
	if !global.master {
		global.master = true
		setupGlobals()
		return
	}

	global.assets = map[string]shared.Asset{}
	res := make(chan packResp)
	global.packs <- packReq{tpe: Reset, res: res}
	<-res
}
