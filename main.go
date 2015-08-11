package main

import "fmt"
import "os"
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

func parse(f string, r io.Reader) {
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
}

func main() {
	path := "./scala-docs-2.11.7/api/scala-library/scala/collection/"

	fs, _ := ioutil.ReadDir(path)

	fmt.Printf("found %v files\n", len(fs))
	for _, f := range fs {
		if !f.IsDir() {
			r, err := os.Open(path + f.Name())
			if err != nil {
				fmt.Printf("can't open file %v, err %v\n", f.Name(), err)
			}
			parse(f.Name(), r)
		}
	}
}
