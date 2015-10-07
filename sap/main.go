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

var packDir string

func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got ping request.")
}

func packsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got packs request.")

	fs, err := ioutil.ReadDir(packDir)
	if err != nil {
		log.Printf("Could not read packs directory: %v\n", err)
		http.Error(w, "Internal server error", 500)
		return
	}

	var packs []shared.Pack
	for _, fi := range fs {
		if strings.HasSuffix(fi.Name(), ".zip") {
			log.Printf("Found zip file: %v", fi.Name())
		}

		r, err := zip.OpenReader(filepath.Join(packDir, fi.Name()))
		if err != nil {
			log.Printf("Error opening archive %v: %v", fi.Name(), err)
			continue
		}

		defer r.Close() // TODO

		for _, f := range r.File {
			if f.Name == "pack.json" {
				log.Printf("Found file in zip: %v\n", f.Name)
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

				packs = append(packs, pack)
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
	flag.StringVar(&packDir, "packdir", "packs", "Path where to find packs")
	flag.Parse()

	serve("0.0.0.0:3025")
}
