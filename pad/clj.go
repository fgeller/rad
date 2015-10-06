package main

import (
	"../shared"

	"encoding/xml"
	"io"
	"log"
)

func parseClojureDocFile(filePath string, r io.Reader) []shared.Namespace {

	var ns []shared.Namespace
	// TODO: pull this out
	d := xml.NewDecoder(r)
	d.Strict = false
	d.AutoClose = xml.HTMLAutoClose
	d.Entity = xml.HTMLEntity

	var t xml.Token
	var err error

	var inMeta bool
	var inMetaA bool
	var inName bool
	var charData []byte
	var name string

	for ; err == nil; t, err = d.Token() {
		se, gotStartElement := t.(xml.StartElement)
		ee, gotEndElement := t.(xml.EndElement)
		cd, gotCharData := t.(xml.CharData)

		switch {
		case gotCharData:
			charData = append(charData, cd...)
		case gotEndElement:
			switch {
			case inName && ee.Name.Local == "h1":
				// &lt;!!
				name = string(charData)
				inName = false
			case inMeta && ee.Name.Local == "div":
				inMeta = false
			case inMetaA && ee.Name.Local == "a":
				// clojure.core.async
				path := string(charData)
				tgt := filePath
				n := shared.Namespace{
					Path:    path,
					Members: []shared.Member{{Name: name, Target: tgt}},
				}
				ns = append(ns, n)
				inMetaA = false
			}
		case gotStartElement:
			charData = []byte{}
			switch {
			// <h1 class="var-name">&lt;!!</h1>
			case se.Name.Local == "h1" && hasAttr(se, "class", "var-name"):
				inName = true
			case se.Name.Local == "div" && hasAttr(se, "class", "var-meta"):
				// <div class="var-meta"><h4><a href="../clojure.core.async.1.html">clojure.core.async</a></h4><span>Available in 1.6</span></div>
				inMeta = true
			case inMeta && se.Name.Local == "a":
				inMetaA = true
			}
		}
	}

	if len(ns) == 0 {
		log.Printf("Found no entries for %v\n", filePath)
	}

	return shared.Merge(ns)
}
