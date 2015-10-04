package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type searchRequest struct {
	Pack   string
	Path   string
	Member string
	Limit  int
}

func socket(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ws" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	start := time.Now()

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error while upgrading request: %v\n", err)
		return
	}
	defer c.Close()

	var req searchRequest
	err = c.ReadJSON(&req)
	if err != nil {
		log.Printf("Failed to read request: %v", err)
		return
	}
	log.Printf("Received search request %v\n", req)

	streamFind(c, req.Pack, req.Path, req.Member, req.Limit)
	log.Printf("Finished request %v in %v\n", req, time.Since(start))
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	pack := r.FormValue("pk")
	ns := r.FormValue("ns")
	mem := r.FormValue("m")
	limit, err := strconv.ParseInt(r.FormValue("limit"), 10, 32)
	if err != nil {
		limit = 10
	}

	start := time.Now()
	res, err := find(pack, ns, mem, int(limit))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	log.Printf(
		"Request pk[%v] and ns[%v] and m[%v] found [%v] entries in %v.",
		pack,
		ns,
		mem,
		len(res),
		time.Since(start),
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
	http.HandleFunc("/ws", socket)

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
