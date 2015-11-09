package main

import (
	"../shared"
)

var global struct {
	packs        chan packReq
	master       bool // should make this a mutex
	assets       map[string]shared.Asset
	buildVersion string
}

func setupGlobals() {
	global.packs = make(chan packReq)
	global.assets = map[string]shared.Asset{}
	global.buildVersion = ""
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
