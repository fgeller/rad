package main

import (
	"../shared"
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strings"
)

func parseGoDocFile(filePath string, r io.Reader) []shared.Namespace {

	var ns []shared.Namespace

	d := xml.NewDecoder(r)
	d.Strict = false
	d.AutoClose = xml.HTMLAutoClose
	d.Entity = xml.HTMLEntity

	var t xml.Token
	var err error

	var importPat = regexp.MustCompile(`^import "(.+)"$`)

	var inPkgIndex bool
	var inPkgOverview bool
	var charData []byte
	var path string

	for ; err == nil; t, err = d.Token() {
		se, gotStartElement := t.(xml.StartElement)
		ee, gotEndElement := t.(xml.EndElement)
		cd, gotCharData := t.(xml.CharData)

		switch {
		case gotCharData:
			charData = append(charData, cd...)
		case gotEndElement:
			switch {
			case ee.Name.Local == "code" && inPkgOverview:
				// import "archive/tar"
				ctnt := string(charData)
				// ["import \"archive/tar\"", "archive/tar"]
				matches := importPat.FindStringSubmatch(ctnt)
				if len(matches) == 2 {
					// ["archive", "tar"]
					path = strings.Join(strings.Split(matches[1], "/"), ".")
				}
			}
		case gotStartElement:
			charData = []byte{}
			switch {
			case se.Name.Local == "h2" && hasAttr(se, "id", "pkg-overview"):
				inPkgOverview = true
			case se.Name.Local == "h3" && hasAttr(se, "id", "pkg-index"):
				inPkgIndex = true
				inPkgOverview = false
			case se.Name.Local == "h4" && inPkgIndex:
				inPkgIndex = false
			case se.Name.Local == "a" && inPkgIndex:
				href, err := attr(se, "href")
				if err != nil {
					return ns // TODO log error
				}

				n, err := parseGoHref(filePath, path, href)
				if err == nil {
					ns = append(ns, n)
				}
			}
		}
	}

	return shared.Merge(ns)
}

func parseGoHref(filePath string, path string, href string) (shared.Namespace, error) {
	n := shared.Namespace{
		Path:    path,
		Members: []shared.Member{},
	}

	if strings.Index(href, "#pkg-") > 0 {
		return n, fmt.Errorf("Unapplicable href [%v]", href)
	}

	// tar.html#Header.FileInfo
	// Header.FileInfo
	frg := href[strings.Index(href, "#")+1:]

	var tpe string
	m := frg
	sep := strings.Index(frg, ".")
	withType := sep > 0
	if withType {
		tpe = frg[:sep]
		m = frg[sep+1:]
	}

	tgt := fmt.Sprintf("%v#%v", filePath, frg)
	n.Members = []shared.Member{{Name: m, Target: tgt}}
	if withType {
		n.Path = n.Path + "." + tpe
	}

	return n, nil
}
