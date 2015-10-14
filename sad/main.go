package main

import (
	"flag"
	"log"
	"path/filepath"
)

var config struct {
	packDir string
	sapAddr string
	addr    string
}

func main() {
	flag.StringVar(&config.packDir, "packdir", "packs", "Path where packages will be installed")
	flag.StringVar(&config.sapAddr, "sapaddr", "geller.io:3025", "Addr where sap serves")
	flag.StringVar(&config.addr, "addr", "0.0.0.0:3024", "Addr where sad should serve")
	flag.Parse()

	pd, err := filepath.Abs(config.packDir)
	if err != nil {
		log.Fatalf("Can't find absolute path for %v: %v\n", config.packDir, err)
	}
	config.packDir = pd

	setupGlobals()

	loadInstalled()
	registerAssets()
	serve(config.addr)
}
