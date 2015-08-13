package main

import "runtime"
import "fmt"
import "os"
import "strings"
import "time"
import "regexp"
import "io"
import "io/ioutil"
import "encoding/xml"

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
	namespace []string
	entity    string
	function  string
	signature string
}

func (e entry) String() string {
	return fmt.Sprintf("entry{namespace: %v, entity: %v, function: %v, signature: %v}", e.namespace, e.entity, e.function, e.signature)
}

func (e entry) eq(other entry) bool {
	if len(e.namespace) != len(other.namespace) {
		return false
	}

	for i, n := range e.namespace {
		if other.namespace[i] != n {
			return false
		}
	}

	return e.entity == other.entity &&
		e.function == other.function &&
		e.signature == other.signature
}

func parseEntry(s string) (entry, error) {
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
			return e, fmt.Errorf("incorrect match for string [%v]", s)
		}
	}

	e.namespace = strings.Split(ms[0][1], ".")
	e.entity = ms[0][2]
	if len(ms[0]) == 5 {
		e.function = ms[0][3]
		e.signature = ms[0][4]
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
						e, err := parseEntry(subs[1])
						if err != nil {
							fmt.Println(err)
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

func scan(dir string) (int, []entry, error) {

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
				fs, es, _ := scan(path)
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

func main() {
	path := "./scala-docs-2.11.7/api/"
	start := time.Now()
	fc, es, _ := scan(path)
	elapsed := time.Now().Sub(start)
	fmt.Printf("found %v links (%.1ff/s).\n", len(es), float64(fc)/elapsed.Seconds())
}
