package main

import "fmt"
import "os"
import "time"
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
					hrefs = append(hrefs, href)
				}
			}
		}
	}

	fmt.Printf("found %v hrefs in %v.\n", len(hrefs), f)
	return len(hrefs)
}

func scan(dir string) (int, int, error) {
	var fileCount, hrefCount int

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("can't read dir %v, err %v\n", dir, err)
		return 0, 0, err
	}

	for _, f := range fs {
		path := dir + string(os.PathSeparator) + f.Name()

		if f.IsDir() {
			fc, lc, err := scan(path)
			if err != nil {
				return fileCount, hrefCount, err
			}
			fileCount += fc
			hrefCount += lc
			continue
		}

		r, err := os.Open(path)
		if err != nil {
			fmt.Printf("can't open file %v, err %v\n", path, err)
			continue
		}

		lc := parse(path, r)
		fileCount++
		hrefCount += lc
	}

	return fileCount, hrefCount, nil
}

func main() {
	path := "./scala-docs-2.11.7/api/scala-library/scala/"
	start := time.Now()
	fc, lc, _ := scan(path)
	elapsed := time.Now().Sub(start)
	fmt.Printf("found %v links (%.1ff/s).\n", lc, float64(fc)/elapsed.Seconds())
}
