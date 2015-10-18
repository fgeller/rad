package main

import (
	"../shared"

	"encoding/xml"
	"io"
	"log"
	"path/filepath"
	"strings"
)

func parseReactDocFile(filePath string, r io.Reader) []shared.Namespace {

	var nss []shared.Namespace
	d := xml.NewDecoder(r)
	d.Strict = false
	d.AutoClose = xml.HTMLAutoClose
	d.Entity = xml.HTMLEntity

	var t xml.Token
	var err error
	var inH3 bool
	var anchor string
	var desc string

	// <h3>
	// 	<a class="anchor" name="updating-shouldcomponentupdate"></a>
	// 	Updating: shouldComponentUpdate
	//   <a class="hash-link" href="component-specs.html#updating-shouldcomponentupdate">#</a>
	// </h3>

	for ; err == nil; t, err = d.Token() {
		se, gotStartElement := t.(xml.StartElement)
		ee, gotEndElement := t.(xml.EndElement)
		cd, gotCharData := t.(xml.CharData)

		switch {
		case gotCharData && inH3 && len(anchor) > 0:
			desc += string(cd)
		case gotEndElement:
			switch {
			case ee.Name.Local == "h3":
				inH3 = false
			}
		case gotStartElement:
			switch {
			case se.Name.Local == "h3":
				inH3 = true
			case inH3 && se.Name.Local == "a" && hasAttr(se, "class", "anchor"):
				anchor, _ = attr(se, "name")
			case inH3 && se.Name.Local == "a" && hasAttr(se, "class", "hash-link"):
				fn := filepath.Base(filePath)
				pth := fn[:strings.Index(fn, ".html")]
				n := strings.TrimSpace(string(desc))

				switch {
				case strings.IndexAny(n, ":.") > 0:
					// desc = "Mounting: componentWillMount"
					// desc = "React.Component"
					pth += "." + n[:strings.IndexAny(n, ":.")]
					n = strings.TrimSpace(n[strings.IndexAny(n, ":.")+1:])
				}

				m := shared.Member{Name: n, Target: filePath + "#" + anchor}
				ns := shared.Namespace{Path: pth, Members: []shared.Member{m}}
				nss = append(nss, ns)
				anchor = ""
				desc = ""
			}
		}
	}

	if len(nss) == 0 {
		log.Printf("Found no entries for %v\n", filePath)
	}

	rs := shared.Merge(nss)
	return rs
}

func parseReactHref(href string, path string) shared.Namespace {

	var ns shared.Namespace
	// href="component-specs.html#unmounting-componentwillunmount"
	// href="component-specs.html#render"

	// pth = "component-specs"
	pth := href[:strings.Index(href, ".html")]
	ns.Path = pth

	// frag="render"
	frag := href[strings.Index(href, "#")+1:]
	m := shared.Member{Name: frag, Target: path + "#" + frag}
	ns.Members = []shared.Member{m}
	return ns
}
