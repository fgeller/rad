package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var docs = map[string][]entry{}
var packDir = "packs"

type indexer func() ([]entry, error)
type parser func(string, io.Reader) []entry
type pack struct {
	name    string
	url     string
	indexer indexer
}
type entry struct {
	Namespace []string
	Entity    string
	Function  string
	Signature string
	Target    string // location relative to `packDir` where to find documentation
	source    string
}

func (e entry) String() string {
	return fmt.Sprintf("entry{Namespace: %v, Entity: %v, Function: %v, Signature: %v}", e.Namespace, e.Entity, e.Function, e.Signature)
}

func (e entry) eq(other entry) bool {
	if len(e.Namespace) != len(other.Namespace) {
		return false
	}

	for i, n := range e.Namespace {
		if other.Namespace[i] != n {
			return false
		}
	}

	return e.Entity == other.Entity &&
		e.Function == other.Function &&
		e.Signature == other.Signature // TODO: expand
}

func scanFile(path string, p parser) []entry {
	r, err := os.Open(path)
	defer r.Close()
	if err != nil {
		fmt.Printf("can't open file %v, err %v\n", path, err)
		return []entry{}
	}

	return p(path, r)
}

func findDirsAndMarkupFiles(dir string) ([]os.FileInfo, error) {
	files := []os.FileInfo{}

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("can't read dir %v, err %v\n", dir, err)
		return files, err
	}

	for _, f := range fs {
		if f.IsDir() ||
			strings.HasSuffix(f.Name(), "html") ||
			strings.HasSuffix(f.Name(), "xml") {
			files = append(files, f)
		}
	}

	return files, nil
}

type scanResult struct {
	entries        []entry
	processedFiles int
}

func scanDir(dir string, p parser) (int, []entry, error) {

	files, err := findDirsAndMarkupFiles(dir)
	if err != nil {
		fmt.Printf("can't read dir %v, err %v\n", dir, err)
		return 0, []entry{}, err
	}

	rc := make(chan scanResult)
	runtime.GOMAXPROCS(runtime.NumCPU())

	for _, fi := range files {
		go func(dir string, f os.FileInfo, c chan scanResult) {
			path := dir + string(os.PathSeparator) + f.Name()
			switch {
			case f.IsDir():
				fs, es, _ := scanDir(path, p)
				c <- scanResult{es, fs}
			default:
				c <- scanResult{scanFile(path, p), 1}
			}
		}(dir, fi, rc)
	}

	results := []entry{}
	fc := 0
	for i := 0; i < len(files); i++ {
		r := <-rc
		fc += r.processedFiles
		results = append(results, r.entries...)
	}

	return fc, results, nil
}

func scan(path string, p parser) ([]entry, error) {
	start := time.Now()
	fc, es, err := scanDir(path, p)
	elapsed := time.Now().Sub(start)
	log.Printf("found %v links (%.1ff/s).\n", len(es), float64(fc)/elapsed.Seconds())

	return es, err
}

func unzip(src string, dest string) error {
	err := os.MkdirAll(dest, 0755)
	if err != nil {
		return err
	}

	r, err := zip.OpenReader(src)
	if err != nil {
		log.Fatal("failed to open zip: %v", src)
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		path := dest + string(os.PathSeparator) + f.Name
		if f.FileInfo().IsDir() {
			os.Mkdir(path, f.Mode())
			continue
		}

		fc, err := f.Open()

		if err != nil {
			return err
		}

		dst, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		_, err = io.Copy(dst, fc)
		if err != nil {
			return err
		}
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

type downloader func(string) (*http.Response, error)

func download(d downloader, remote string) (string, error) {
	local := remote[strings.LastIndex(remote, "/")+1:]
	if fileExists(local) {
		log.Printf("Already downloaded [%v].", local)
		return local, nil
	}

	out, err := os.Create(local)
	if err != nil {
		return "", err
	}
	defer out.Close()

	resp, err := http.Get(remote)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	log.Printf("Downloading [%v] to local [%v].\n", remote, local)

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	log.Printf("Downloaded %v bytes.\n", n)

	return local, nil
}

func install(pack pack) error {
	dataPath := packDir + string(os.PathSeparator) +
		pack.name + string(os.PathSeparator) +
		"rad-data.json"

	if fileExists(dataPath) {
		log.Printf("Already installed pack [%v], deserializing entries.", pack.name)
		start := time.Now()

		data, err := ioutil.ReadFile(dataPath)
		if err != nil {
			return err // TODO: or re-download?
		}

		var es []entry
		err = json.Unmarshal(data, &es)
		if err != nil {
			return err // TODO: or re-download?
		}

		docs[pack.name] = es
		log.Printf(
			"Deserialized [%v] entries for [%v] in %v.",
			len(es),
			pack.name,
			time.Since(start),
		)

		return nil
	}

	local, err := download(http.Get, pack.url)
	if err != nil {
		log.Fatalf("Failed to download [%v] err: %v.\n", pack.url, err)
		return err
	}
	defer os.Remove(local)

	err = unzip(local, packDir+string(os.PathSeparator)+pack.name)
	if err != nil {
		log.Fatalf("Failed to unzip archive [%v], err: %v", local, err)
		return err
	}

	docs[pack.name], err = pack.indexer()
	if err != nil {
		return err
	}

	datPath := packDir + string(os.PathSeparator) +
		pack.name + string(os.PathSeparator) +
		"rad-data.json"

	data, err := json.Marshal(docs[pack.name])
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(datPath, data, 0644)
	if err != nil {
		return err
	}

	log.Printf("Installed [%v] entries for pack [%v].", len(docs[pack.name]), pack.name)
	return nil
}

func findEntityFunction(pack string, entity string, fun string, limit int) ([]entry, error) {
	es, ok := docs[pack]
	if !ok {
		return es, fmt.Errorf("Package [%v] not installed.", pack)
	}

	results := []entry{}

	for _, e := range es {
		if strings.HasPrefix(strings.ToLower(e.Entity), strings.ToLower(entity)) &&
			strings.HasPrefix(strings.ToLower(e.Function), strings.ToLower(fun)) {
			results = append(results, e)
			if len(results) == limit {
				return results, nil
			}
		}
	}

	return results, nil
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	pack := r.FormValue("p")
	entity := r.FormValue("e")
	fun := r.FormValue("f")
	limit, err := strconv.ParseInt(r.FormValue("limit"), 10, 32)
	if err != nil {
		limit = 10
	}

	res, _ := findEntityFunction(pack, entity, fun, int(limit))
	log.Printf("got request for p[%v] and e[%v] and f[%v], found [%v] entries.", pack, entity, fun, len(res))

	js, _ := json.Marshal(res) // TODO: return proper err

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

func main() {
	install(
		pack{
			name:    "scala",
			url:     "http://downloads.typesafe.com/scala/2.11.7/scala-docs-2.11.7.zip",
			indexer: indexScalaApi("scala"),
		},
	)

	serve("0.0.0.0:3024")
}
