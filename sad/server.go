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
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type controlRequest struct {
	Typ  string
	Data interface{}
}

type controlResponse struct {
	Typ  string
	Data interface{}
}

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

	if len(m) == 0 {
		m = pt
		pt = ".*"
	}

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

type PackInfo struct {
	Installed []shared.Pack
	Available []shared.Pack
}

type StatusInfo struct {
	Packs   PackInfo
	Version string
}

func availablePacks() []shared.Pack {
	availablePacks := []shared.Pack{}
	res, err := http.Get("http://" + config.sapAddr + "/packs")
	if err != nil {
		log.Printf("Error requesting available packs: %v\n", err)
	}

	if err == nil && res.StatusCode == 200 {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error reading data from response: %v\n", err)
		} else {
			json.Unmarshal(data, &availablePacks)
		}
	}

	return availablePacks
}

func dispatchControl(sock *websocket.Conn) {
	var req controlRequest
	defer sock.Close()

	for {
		err := sock.ReadJSON(&req)
		if err != nil {
			log.Printf("Closing control socket. err=%v", err)
			return
		}
		log.Printf("Received control request %v\n", req)

		switch req.Typ {
		case "StatusRequest":
			go streamStatus(sock)
		default:
			log.Printf("UNKNOWN REQUEST TYPE %v", req.Typ)
		}
	}
}

func streamStatus(sock *websocket.Conn) {
	installed := installedPacks()
	filteredInstalled := []shared.Pack{}
	available := availablePacks()

loopinstalled:
	for _, ip := range installed {
		for ai, ap := range available {
			if ip.Installing && ip.File == ap.File {
				available[ai].Installing = true
				continue loopinstalled
			}
		}
		filteredInstalled = append(filteredInstalled, ip)
	}

	info := StatusInfo{
		Packs: PackInfo{
			Installed: filteredInstalled,
			Available: available,
		},
		Version: global.buildVersion,
	}

	err := sock.WriteJSON(controlResponse{Typ: "StatusResponse", Data: info})
	if err != nil {
		log.Printf("Error while writing status info: %v\n", err)
		return
	}
}

func controlSocket(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/c" {
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

	go dispatchControl(c)
}

func query(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/q" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	query := r.URL.Query()
	l, err := strconv.Atoi(strings.Join(query["limit"], ""))
	if err != nil {
		l = 100
	}
	req := searchRequest{
		Pack:   strings.Join(query["pack"], ""),
		Path:   strings.Join(query["path"], ""),
		Member: strings.Join(query["member"], ""),
		Limit:  int(l),
	}
	log.Printf("Received search request %v\n", req)

	params, err := compileParams(req.Pack, req.Path, req.Member)
	if err != nil {
		log.Printf("Error while compiling params: %v\n", err)
		return
	}

	found := searchResults(params, req.Limit)
	js, err := json.Marshal(found)
	if err != nil {
		log.Printf("Failed to marshal err: %v\n", err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}

func searchResults(params searchParams, limit int) []searchResult {
	start := time.Now()
	results := make(chan searchResult)
	control := make(chan struct{}, 1)
	sr := []searchResult{}

	go find(results, control, params)
	for {
		res, ok := <-results
		if !ok {
			log.Printf("Finished request in %v\n", time.Since(start))
			return sr
		}

		sr = append(sr, res)
		log.Printf("Found result #%v after %v\n", len(sr), time.Since(start))
		if len(sr) >= limit {
			log.Printf("Finished request after hitting limit in %v\n", time.Since(start))
			control <- struct{}{}
			return sr
		}
	}
}

func searchSocket(w http.ResponseWriter, r *http.Request) {
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
	if config.readOnly {
		log.Printf("Failing install request, read only mode.\n")
		http.Error(w, "Forbidden", 403)
		return
	}

	fn := r.URL.Path[len("/install/"):]

	resp := make(chan packResp)
	req := packReq{tpe: Install, pck: shared.Pack{File: fn}, res: resp}

	global.packs <- req
	msg, ok := <-resp
	if ok && msg.err != nil {
		log.Printf("Error installing %v: %v\n", fn, msg.err)
		http.Error(w, "Internal server error", 500)
	}
}

func removeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got remove request %v\n", r.URL.Path)
	if config.readOnly {
		log.Printf("Failing remove request, read only mode.\n")
		http.Error(w, "Forbidden", 403)
		return
	}

	pn := r.URL.Path[len("/remove/"):]
	remove(pn)
}

func assetHandler(w http.ResponseWriter, r *http.Request) {
	an := r.URL.Path
	switch {
	case an == "/":
		an = "index.html"
	case strings.HasPrefix(an, "/a/"):
		an = an[len("/a/"):]
	case strings.HasPrefix(an, "/"):
		an = an[1:]
	}

	a, ok := global.assets[an]
	if !ok {
		log.Printf("Got asset request %v but not available.\n", r.URL.Path)
		http.Error(w, "404 page not found", 404)
		return
	}

	log.Printf("Serving asset request %v with %v.\n", r.URL.Path, a)
	w.Header().Set("Content-Type", a.ContentType)
	w.Write(a.Content)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got ping for %v\n", r.URL.Path)
	w.Write([]byte("pong"))
}

func awaitPing(addr string) error {
	limit := 50
	attempts := 0

	for {
		resp, err := http.Get("http://" + addr + "/ping")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		attempts++
		if attempts > limit {
			return fmt.Errorf("Got no ping on %v.", addr)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

type root string

func (r root) Open(n string) (http.File, error) {

	rn := strings.Replace(n, "\\", "%5C", -1)
	rn = strings.Replace(rn, "|", "%7C", -1)
	rn = strings.Replace(rn, "\"", "%22", -1)
	rn = strings.Replace(rn, "*", "%2A", -1)
	rn = strings.Replace(rn, "<", "%3C", -1)
	rn = strings.Replace(rn, ">", "%3E", -1)
	t := filepath.Join(string(r), rn)
	t = shared.MaybeEscapeWinPath(t)

	return os.Open(t)
}

func serve(addr string) {
	http.HandleFunc("/ping/", pingHandler)
	http.HandleFunc("/q", query)
	http.HandleFunc("/s", searchSocket)
	http.HandleFunc("/c", controlSocket)
	http.HandleFunc("/a/", assetHandler)
	if config.devMode {
		log.Printf("Serving assets from ui folder.\n")
		http.Handle("/", http.FileServer(http.Dir("ui")))
	} else {
		http.HandleFunc("/", assetHandler)
	}

	ps := http.FileServer(root(config.packDir))
	http.Handle("/pack/", http.StripPrefix("/pack/", ps))

	log.Printf("Serving on addr http://%v\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
