package main

import (
	"../shared"
	"flag"
)

var docs = map[string][]shared.Namespace{}
var packDir string

func main() {
	flag.StringVar(&packDir, "packDir", "packs", "Path where packages will be installed")
	flag.Parse()

	serve("0.0.0.0:3024")
}
