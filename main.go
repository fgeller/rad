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

func parseScalaEntry(s string) (entry, error) {
	e := entry{}
	pat, err := regexp.Compile("(.+)\\.([^@]+)@(.+?)(\\(.*)$")
	if err != nil {
		return e, nil
	}

	ms := pat.FindAllStringSubmatch(s, -1)
	if len(ms[0]) != 5 {
		return e, fmt.Errorf("incorrect match for string [%v]", s)
	}
	e.namespace = strings.Split(ms[0][1], ".")
	e.entity = ms[0][2]
	e.function = ms[0][3]
	e.signature = ms[0][4]

	return e, nil
}

func parse(f string, r io.Reader) int {
	d := xml.NewDecoder(r)
	hrefs := []string{}
	var t xml.Token
	var err error

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

						fmt.Printf("found fragment %v\n", subs[1])
						hrefs = append(hrefs, href)

					}
				}
			}
		}
	}

	return len(hrefs)
}

type empty struct{}

func scan(dir string) (int, int, error) {
	var fileCount, hrefCount int

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("can't read dir %v, err %v\n", dir, err)
		return 0, 0, err
	}

	files := []os.FileInfo{}
	dirs := []os.FileInfo{}

	for _, f := range fs {
		switch {
		case f.IsDir():
			dirs = append(dirs, f)
		case strings.HasSuffix(f.Name(), "html") || strings.HasSuffix(f.Name(), "xml"):
			files = append(files, f)
		}
	}

	runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	sem := make(chan empty, len(files))
	for _, f := range files {
		go func(f os.FileInfo) {
			path := dir + string(os.PathSeparator) + f.Name()

			r, err := os.Open(path)
			defer r.Close()
			if err != nil {
				fmt.Printf("can't open file %v, err %v\n", path, err)
				sem <- empty{}
				return
			}

			lc := parse(path, r)
			fileCount++
			hrefCount += lc
			sem <- empty{}
		}(f)
	}

	for i := 0; i < len(files); i++ {
		<-sem
	}

	for _, d := range dirs {
		path := dir + string(os.PathSeparator) + d.Name()
		fc, lc, err := scan(path)
		if err != nil {
			return 0, 0, err
		}
		fileCount += fc
		hrefCount += lc
	}

	return fileCount, hrefCount, nil
}

func main() {
	path := "./scala-docs-2.11.7/api/"
	start := time.Now()
	fc, lc, _ := scan(path)
	elapsed := time.Now().Sub(start)
	fmt.Printf("found %v links (%.1ff/s).\n", lc, float64(fc)/elapsed.Seconds())
}
