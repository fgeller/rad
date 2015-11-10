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
	packDir string
	addr    string
}

type packListing struct {
	File string
	*shared.Pack
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got ping request.")
}

func packHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got pack request for %v", r.URL.Path)
	http.ServeFile(w, r, filepath.Join(config.packDir, r.URL.Path[6:]))
}

func packsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got packs request.")

	fs, err := ioutil.ReadDir(config.packDir)
	if err != nil {
		log.Printf("Could not read packs directory: %v\n", err)
		http.Error(w, "Internal server error", 500)
		return
	}

	log.Printf("Found %v files in %v\n", len(fs), config.packDir)

	var packs []packListing
	for _, fi := range fs {
		if strings.HasSuffix(fi.Name(), ".zip") {
			log.Printf("Found zip file: %v", fi.Name())
		}

		r, err := zip.OpenReader(filepath.Join(config.packDir, fi.Name()))
		if err != nil {
			log.Printf("Error opening archive %v: %v", fi.Name(), err)
			continue
		}
		defer r.Close()

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

				packs = append(packs, packListing{Pack: &pack, File: fi.Name()})
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

func serve(addr string) error {
	http.HandleFunc("/packs", packsHandler)
	http.HandleFunc("/pack/", packHandler)
	http.HandleFunc("/ping", pingHandler)
	log.Printf("Serving on %v\n", addr)
	return http.ListenAndServe(addr, nil)
}

func main() {
	flag.StringVar(&config.packDir, "packdir", "packs", "Path where to find packs")
	flag.StringVar(&config.addr, "addr", "0.0.0.0:3025", "Where to listen.")
	flag.Parse()

	log.Printf("Stopped, err: %v\n", serve(config.addr))
}
