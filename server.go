package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type searchResult struct {
	Entity    string
	Namespace []string
	Member    string
	Signature string
	Target    string
	Source    string
}

func (e entry) searchResult() searchResult {
	return searchResult{
		Entity:    e.Name,
		Namespace: e.Namespace,
		Member:    e.Members[0].Name,
		Signature: e.Members[0].Signature,
		Target:    e.Members[0].Target,
		Source:    e.Members[0].Source,
	}
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	pack := r.FormValue("p")
	entity := r.FormValue("e")
	fun := r.FormValue("m")
	limit, err := strconv.ParseInt(r.FormValue("limit"), 10, 32)
	if err != nil {
		limit = 10
	}

	res, _ := findEntityMember(pack, entity, fun, int(limit))
	log.Printf("got request for p[%v] and e[%v] and f[%v], found [%v] entries.", pack, entity, fun, len(res))

	results := []searchResult{}
	for _, sr := range res {
		results = append(results, sr.searchResult())
	}
	js, _ := json.Marshal(results) // TODO: return proper err

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
