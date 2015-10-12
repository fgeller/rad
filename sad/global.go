package main

import (
	"../shared"
)

var global struct {
	packs  map[string]shared.Pack
	docs   map[string][]shared.Namespace
	assets map[string]shared.Asset
}

func resetGlobals() {
	global.packs = map[string]shared.Pack{}
	global.docs = map[string][]shared.Namespace{}
	global.assets = map[string]shared.Asset{}
}
