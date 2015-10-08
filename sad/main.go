package main

import (
	"../shared"

	"flag"
)

var global struct {
	packs map[string]shared.Pack
	docs  map[string][]shared.Namespace
}
var config struct {
	packDir string
	sapAddr string
}

func main() {
	flag.StringVar(&config.packDir, "packdir", "packs", "Path where packages will be installed")
	flag.StringVar(&config.sapAddr, "sapaddr", "localhost:3025", "Addr where sap is running")
	flag.Parse()

	loadInstalled()
	serve("0.0.0.0:3024")
}
