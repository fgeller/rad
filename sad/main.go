package main

import (
	"flag"
)

var config struct {
	packDir string
	sapAddr string
	addr    string
}

func main() {
	flag.StringVar(&config.packDir, "packdir", "packs", "Path where packages will be installed")
	flag.StringVar(&config.sapAddr, "sapaddr", "localhost:3025", "Addr where sap serves")
	flag.StringVar(&config.addr, "addr", "0.0.0.0:3024", "Addr where sad should serve")
	flag.Parse()

	resetGlobals()

	loadInstalled()
	registerAssets()
	serve(config.addr)
}
