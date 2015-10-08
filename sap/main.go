package main

import (
	"../shared"

	"archive/zip"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var config struct {
	PackDir string
}

type packListing struct {
	Path string
	*shared.Pack
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got ping request.")
}

func packsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got packs request.")

	fs, err := ioutil.ReadDir(config.PackDir)
	if err != nil {
		log.Printf("Could not read packs directory: %v\n", err)
		http.Error(w, "Internal server error", 500)
		return
	}

	var packs []packListing
	for _, fi := range fs {
		if strings.HasSuffix(fi.Name(), ".zip") {
			log.Printf("Found zip file: %v", fi.Name())
		}

		r, err := zip.OpenReader(filepath.Join(config.PackDir, fi.Name()))
		if err != nil {
			log.Printf("Error opening archive %v: %v", fi.Name(), err)
			continue
		}

		defer r.Close() // TODO

		for _, f := range r.File {
			if filepath.Base(f.Name) == "pack.json" {
				fh, err := f.Open()
				if err != nil {
					log.Printf("Error opening %v: %v", f.Name, err)
					continue
				}

				data, err := ioutil.ReadAll(fh)
				if err != nil {
					log.Printf("Error reading from %v: %v", f.Name, err)
					continue
				}

				var pack shared.Pack
				err = json.Unmarshal(data, &pack)
				if err != nil {
					log.Printf("Error unmarshaling data from %v: %v", f.Name, err)
					continue
				}

				packs = append(packs, packListing{Pack: &pack, Path: "/pack/" + fi.Name()})
			}
		}
	}

	data, err := json.Marshal(packs)
	if err != nil {
		log.Printf("Could not marshal packs: %v\n", err)
		http.Error(w, "Internal server error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func serve(addr string) {
	http.HandleFunc("/packs", packsHandler)
	http.HandleFunc("/ping", pingHandler)
	log.Printf("Serving on %v\n", addr)
	http.ListenAndServe(addr, nil)
}

func main() {
	flag.StringVar(&config.PackDir, "packdir", "packs", "Path where to find packs")
	flag.Parse()

	serve("0.0.0.0:3025")
}
