package main

import (
	"../shared"
	"flag"
)

var packs = map[string]shared.Pack{}
var docs = map[string][]shared.Namespace{}
var packDir string

func main() {
	flag.StringVar(&packDir, "packdir", "packs", "Path where packages will be installed")
	flag.Parse()

	loadInstalled()
	serve("0.0.0.0:3024")
}
