package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var config struct {
	packDir  string
	sapAddr  string
	addr     string
	readOnly bool
}

func openURL(url string) {
	switch runtime.GOOS {
	case "darwin":
		exec.Command("open", url).Run()
	case "linux":
		exec.Command("xdg-open", url).Run()
	}
}

func waitAndOpenUrl(url string) {
	if err := awaitPing(config.addr); err == nil {
		openURL("http://" + config.addr)
		return
	}

	log.Printf("Couldn't get ping, slow startup?\n")
}

func main() {
	pd := filepath.Join(os.Getenv("HOME"), ".rad", "sad-packs")

	flag.StringVar(&config.packDir, "packdir", pd, "Path where packages will be installed")
	flag.StringVar(&config.sapAddr, "sapaddr", "geller.io:3025", "Addr where sap serves")
	flag.StringVar(&config.addr, "addr", "0.0.0.0:3024", "Addr where sad should serve")
	flag.BoolVar(&config.readOnly, "readonly", false, "Whether to allow modifications of installed packs.")
	flag.Parse()

	pd, err := filepath.Abs(config.packDir)
	if err != nil {
		log.Fatalf("Can't find absolute path for %v: %v\n", config.packDir, err)
	}
	config.packDir = pd

	setupGlobals()
	loadInstalled()
	registerAssets()
	go waitAndOpenUrl("http://" + config.addr)
	serve(config.addr)
}
