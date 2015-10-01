package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

func queryHandler(w http.ResponseWriter, r *http.Request) {
	pack := r.FormValue("p")
	entity := r.FormValue("e")
	mem := r.FormValue("m")
	limit, err := strconv.ParseInt(r.FormValue("limit"), 10, 32)
	if err != nil {
		limit = 10
	}

	res, err := findEntityMember(pack, entity, mem, int(limit))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	log.Printf(
		"got request for p[%v] and e[%v] and m[%v], found [%v] entries.",
		pack,
		entity,
		mem,
		len(res),
	)

	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func serve(addr string) {
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/s", queryHandler)

	pd, err := filepath.Abs(packDir)
	if err != nil {
		log.Fatalf("Can't find absolute path to packDir %v: %v\n", packDir, err)
	}

	packs := http.FileServer(http.Dir(pd))
	http.Handle("/pack/", http.StripPrefix("/pack/", packs))

	ui := http.FileServer(http.Dir("./ui"))
	http.Handle("/ui/", http.StripPrefix("/ui/", ui))

	log.Printf("Serving on addr %v\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
