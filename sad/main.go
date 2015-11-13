package main

import (
	"flag"
	"log"
	"os/exec"
	"os/user"
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
	case "windows":
		exec.Command("cmd.exe", "/C", "start", url).Run()
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
	u, err := user.Current()
	if err != nil {
		log.Fatalf("Couldn't access user: %v", err)
	}
	pd := filepath.Join(u.HomeDir, ".rad", "sad-packs")

	flag.StringVar(&config.packDir, "packdir", pd, "Path where packages will be installed")
	flag.StringVar(&config.sapAddr, "sapaddr", "geller.io:3025", "Addr where sap serves")
	flag.StringVar(&config.addr, "addr", "localhost:3024", "Addr where sad should serve")
	flag.BoolVar(&config.readOnly, "readonly", false, "Whether to allow modifications of installed packs.")
	flag.Parse()

	pd, err = filepath.Abs(config.packDir)
	if err != nil {
		log.Fatalf("Can't find absolute path for %v: %v\n", config.packDir, err)
	}
	config.packDir = pd

	setupGlobals()
	loadInstalled()
	registerBuildVersion()
	registerAssets()
	go waitAndOpenUrl("http://" + config.addr)
	serve(config.addr)
}
