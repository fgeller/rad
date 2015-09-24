package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func queryHandler(w http.ResponseWriter, r *http.Request) {
	pack := r.FormValue("p")
	entity := r.FormValue("e")
	fun := r.FormValue("f")
	limit, err := strconv.ParseInt(r.FormValue("limit"), 10, 32)
	if err != nil {
		limit = 10
	}

	res, err := findEntityMember(pack, entity, fun, int(limit))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	log.Printf("got request for p[%v] and e[%v] and f[%v], found [%v] entries.", pack, entity, fun, len(res))

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

	packs := http.FileServer(http.Dir("./" + packDir))
	http.Handle(fmt.Sprintf("/%v/", packDir), http.StripPrefix(fmt.Sprintf("/%v/", packDir), packs))

	ui := http.FileServer(http.Dir("./ui"))
	http.Handle("/ui/", http.StripPrefix("/ui/", ui))

	log.Printf("Serving on addr %v\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
