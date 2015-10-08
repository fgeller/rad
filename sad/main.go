package main

import (
	"../shared"

	"flag"
)

var packs = map[string]shared.Pack{}
var docs = map[string][]shared.Namespace{}
var packDir string
var sapAddr string

func main() {
	flag.StringVar(&packDir, "packdir", "packs", "Path where packages will be installed")
	flag.StringVar(&sapAddr, "sapaddr", "localhost:3025", "Addr where sap is running")
	flag.Parse()

	loadInstalled()
	serve("0.0.0.0:3024")
}
