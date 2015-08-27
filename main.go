package main

import "runtime"
import "fmt"
import "os"
import "strings"
import "net/http"
import "log"
import "time"
import "regexp"
import "io"
import "io/ioutil"
import "encoding/xml"
import "encoding/json"

import "archive/zip"

var docs = map[string][]entry{}

type indexer func() []entry
type packConfig struct {
	name    string
	url     string
	indexer indexer
}

func attr(se xml.StartElement, name string) (string, error) {
	for _, att := range se.Attr {
		if att.Name.Local == name {
			return att.Value, nil
		}
	}

	return "", fmt.Errorf("could not find attr %v", name)
}

func hasAttr(se xml.StartElement, name string, value string) bool {
	v, err := attr(se, name)
	return err == nil && v == value
}

type entry struct {
	Namespace []string
	Entity    string
	Function  string
	Signature string
	Target    string
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
		e.Signature == other.Signature
}

func parseEntry(source string, target string, s string) (entry, error) {
	e := entry{}
	funPat, err := regexp.Compile("(.+)\\.([^@]+)@(.+?)(\\(.*)?$")
	if err != nil {
		return e, err
	}
	entPat, err := regexp.Compile("(.+)\\.(.+)$")
	if err != nil {
		return e, err
	}

	ms := funPat.FindAllStringSubmatch(s, -1)
	if len(ms) < 1 || len(ms[0]) != 5 {
		ms = entPat.FindAllStringSubmatch(s, -1)
		if len(ms) < 1 || len(ms[0]) != 3 {
			return e, fmt.Errorf("couldn't match Scala entry [%v].", s)
		}
	}

	e.Namespace = strings.Split(ms[0][1], ".")
	e.Entity = ms[0][2]

	targetSplits := strings.Split(target, "/")
	upCount := 0
	for i, v := range targetSplits {
		if v != ".." {
			upCount = i
			break
		}
	}
	sourceSplits := strings.Split(source, "/")
	newSplits := sourceSplits[1 : len(sourceSplits)-(upCount+1)]
	newSplits = append(newSplits, targetSplits[len(targetSplits)-1]+s)
	newTarget := strings.Join(newSplits, "/")
	e.Target = newTarget
	e.source = source
	if len(ms[0]) == 5 {
		e.Function = ms[0][3]
		e.Signature = ms[0][4]
	}

	return e, nil
}

func parse(f string, r io.Reader) []entry {
	d := xml.NewDecoder(r)
	var t xml.Token
	var err error
	entries := []entry{}

	for ; err == nil; t, err = d.Token() {
		if se, ok := t.(xml.StartElement); ok {
			switch {
			case se.Name.Local == "a" && hasAttr(se, "title", "Permalink"):
				// <a href="../../index.html#scala.collection.TraversableLike@WithFilterextendsFilterMonadic[A,Repr]"
				//    title="Permalink"
				//    target="_top">
				//    <img src="../../../lib/permalink.png" alt="Permalink" />
				//  </a>

				href, err := attr(se, "href")
				if err == nil {
					subs := strings.SplitAfterN(href, "#", 2)
					if len(subs) > 1 {
						// fmt.Printf("found fragment %v\n", subs[1])
						e, err := parseEntry(f, subs[0], subs[1])
						if err != nil {
							log.Println(err)
						} else {
							entries = append(entries, e)
						}

					}
				}
			}
		}
	}

	return entries
}

func scanFile(path string) []entry {
	r, err := os.Open(path)
	defer r.Close()
	if err != nil {
		fmt.Printf("can't open file %v, err %v\n", path, err)
		return []entry{}
	}

	return parse(path, r)
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

func scanDir(dir string) (int, []entry, error) {

	files, err := findDirsAndMarkupFiles(dir)
	if err != nil {
		fmt.Printf("can't read dir %v, err %v\n", dir, err)
		return 0, []entry{}, err
	}

	rc := make(chan scanResult)
	runtime.GOMAXPROCS(runtime.NumCPU())

	for _, p := range files {
		go func(dir string, f os.FileInfo, c chan scanResult) {
			path := dir + string(os.PathSeparator) + f.Name()
			switch {
			case f.IsDir():
				fs, es, _ := scanDir(path)
				c <- scanResult{es, fs}
			default:
				c <- scanResult{scanFile(path), 1}
			}
		}(dir, p, rc)
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

func scan(path string) ([]entry, error) {
	start := time.Now()
	fc, es, err := scanDir(path)
	elapsed := time.Now().Sub(start)
	fmt.Printf("found %v links (%.1ff/s).\n", len(es), float64(fc)/elapsed.Seconds())

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

func download(remote string) (string, error) {
	local := remote[strings.LastIndex(remote, "/")+1:]
	fmt.Printf("trying to download to local [%v]\n", local)

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

	// // TODO: identify type
	// unzip(local, )
	// os.Rename("scala-docs-2.11.7", "pkgs/scala")

	return local, nil
}

func indexScalaApi() {
	path := "./pkgs/scala"
	log.Printf("about to index scala api in [%v]\n", path)
	es, err := scan(path)
	if err != nil {
		log.Fatalf("Encountered error while indexing Scala api [%v].", err)
		return
	}

	docs["scala"] = es
}

func index() {
	indexScalaApi()
}

func install(pack packConfig) error {
	local, err := download(pack.url)
	if err != nil {
		log.Fatalf("Failed to download [%v] err: %v.\n", pack.url, err)
		return err
	}
	defer os.Remove(local)

	err = unzip(local, "packs"+string(os.PathSeparator)+pack.name)
	if err != nil {
		log.Fatalf("Failed to unzip archive [%v], err: %v", local, err)
		return err
	}

	docs[pack.name] = pack.indexer()
	log.Printf("Installed [%v] entries for pack [%v].", len(docs[pack.name]), pack.name)
	return nil
}

func findEntries(pkg string, name string) ([]entry, error) {

	es, ok := docs[pkg]
	if !ok {
		return es, fmt.Errorf("Package [%v] not installed.", pkg)
	}

	results := []entry{}
	for _, e := range es {
		if e.Entity == name {
			results = append(results, e)
		}
	}

	return results, nil
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	pkg := r.FormValue("p")
	entity := r.FormValue("e")
	res, _ := findEntries(pkg, entity)
	log.Printf("got request for p[%v] and e[%v], found [%v] entries.", pkg, entity, len(res))

	js, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func serve() {
	http.HandleFunc("/s", queryHandler)

	pkgs := http.FileServer(http.Dir("./pkgs"))
	http.Handle("/pkgs/", http.StripPrefix("/pkgs/", pkgs))

	ui := http.FileServer(http.Dir("./ui"))
	http.Handle("/ui/", http.StripPrefix("/ui/", ui))

	addr := ":3024"

	log.Printf("serving on addr %v\n", addr)
	log.Fatal(http.ListenAndServe(":3024", nil))
}

func main() {
	// download("http://downloads.typesafe.com/scala/2.11.7/scala-docs-2.11.7.zip")
	// index()
	serve()
}
