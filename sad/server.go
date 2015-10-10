package main

import (
	"../shared"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
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

type searchParams struct {
	pack   *regexp.Regexp
	path   *regexp.Regexp
	member *regexp.Regexp
}

type searchResult struct {
	Namespace string
	Member    string
	Target    string
}

func (s searchResult) eq(o searchResult) bool {
	return reflect.DeepEqual(s, o)
}

func NewSearchResult(n shared.Namespace, memberIdx int) searchResult {

	if len(n.Members) == 0 { // TODO: do we need this guy?
		return searchResult{
			Namespace: n.Path,
		}
	}

	return searchResult{
		Namespace: n.Path,
		Member:    n.Members[memberIdx].Name,
		Target:    "/pack/" + n.Members[memberIdx].Target, // TODO: should we fix that here?
	}
}

func maybeInsensitive(pat string) string {
	if strings.ToLower(pat) == pat {
		return fmt.Sprintf("(?i)%v", pat)
	}
	return pat
}

func compileParams(pk, pt, m string) (searchParams, error) {
	var result searchParams
	var pats [3]*regexp.Regexp

	for i, p := range [3]string{pk, pt, m} {
		pat := maybeInsensitive(p)
		cp, err := regexp.Compile(pat)
		if err != nil {
			return result, err
		}
		pats[i] = cp
	}

	result.pack = pats[0]
	result.path = pats[1]
	result.member = pats[2]

	return result, nil
}

func status(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got status request for %v\n", r.URL.Path)

	if r.URL.Path == "/status/packs/installed" {
		js, err := json.Marshal(global.packs)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}

	if r.URL.Path == "/status/packs/available" {
		res, err := http.Get("http://" + config.sapAddr + "/packs")
		if err != nil {
			log.Printf("Error while requesting available packs: %v\n", err)
			http.Error(w, err.Error(), 500)
			return
		}

		if res.StatusCode != 200 {
			log.Printf("Expected 200 but got status code: %v\n", res.StatusCode)
			http.Error(w, err.Error(), 500)
			return
		}

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error reading data from response: %v\n", err)
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	http.Error(w, "Not found", 404)
}

func socket(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/s" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

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

	params, err := compileParams(req.Pack, req.Path, req.Member)
	if err != nil {
		log.Printf("Error while compiling params: %v\n", err)
		return
	}

	streamResults(c, params, req.Limit)
}

func streamResults(sock *websocket.Conn, params searchParams, limit int) {
	start := time.Now()
	count := 0
	results := make(chan searchResult)
	control := make(chan struct{}, 1)

	go find(results, control, params)
	for {
		res, ok := <-results
		if !ok {
			log.Printf("Finished request in %v\n", time.Since(start))
			return
		}

		count++
		log.Printf("Found result #%v after %v\n", count, time.Since(start))

		err := sock.WriteJSON(res)
		if err != nil {
			log.Printf("Error while writing result: %v\n", err)
			control <- struct{}{}
			return
		}

		if count >= limit {
			log.Printf("Finished request after hitting limit in %v\n", time.Since(start))
			control <- struct{}{}
			return
		}
	}
}

func installHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got install request %v\n", r.URL.Path)
	fn := r.URL.Path[len("/install/"):]
	path, err := shared.DownloadToTemp("http://" + config.sapAddr + "/pack/" + fn)
	if err != nil {
		log.Printf("Error downloading pack: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}
	defer os.RemoveAll(path)

	err = install(path)
	if err != nil {
		log.Printf("Error installing pack: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func removeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got remove request %v\n", r.URL.Path)
	pn := r.URL.Path[len("/remove/"):]

	err := remove(pn)
	if err != nil {
		log.Printf("Error removing pack %v: %v\n", pn, err)
		http.Error(w, err.Error(), 500)
		return
	}

	resetGlobal()
	err = loadInstalled()
	if err != nil {
		log.Printf("Error loading installed packs: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func assetHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got asset request %v\n", r.URL.Path)
	an := r.URL.Path[len("/a/"):]

	a, ok := global.assets[an]
	if !ok {
		http.Error(w, "Not found", 404)
		return
	}

	w.Header().Set("Content-Type", a.contentType)
	w.Write(a.content)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got ping for %v\n", r.URL.Path)
	w.Write([]byte("pong"))
}

func serve(addr string) {
	http.HandleFunc("/ping/", pingHandler)
	http.HandleFunc("/s", socket)
	http.HandleFunc("/status/", status)
	http.HandleFunc("/install/", installHandler)
	http.HandleFunc("/remove/", removeHandler)
	http.HandleFunc("/a/", assetHandler)

	pd, err := filepath.Abs(config.packDir)
	if err != nil {
		log.Fatalf("Can't find absolute path to packDir %v: %v\n", config.packDir, err)
	}

	ps := http.FileServer(http.Dir(pd))
	http.Handle("/pack/", http.StripPrefix("/pack/", ps))

	ui := http.FileServer(http.Dir("./ui"))
	http.Handle("/ui/", http.StripPrefix("/ui/", ui))

	log.Printf("Serving on addr %v\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
