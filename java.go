package main

import (
	"encoding/xml"
	"io"
	"log"
	"os"
	"strings"
)

func indexJavaApi(packName string) func() ([]entry, error) {
	return func() ([]entry, error) {
		path := packDir + "/" + packName
		log.Printf("About to index java api in [%v]\n", path)
		return scan(path, parseJavaDocFile)
	}
}

func parseJavaDocFile(path string, r io.Reader) []entry {

	d := xml.NewDecoder(r)
	d.Strict = false
	d.AutoClose = xml.HTMLAutoClose
	d.Entity = xml.HTMLEntity

	var t xml.Token
	var err error
	var inMemberNameLink bool
	var inMemberSummary bool
	entries := []entry{}

	for ; err == nil; t, err = d.Token() {

		// <span class="memberNameLink">
		//   <a href="../../../javax/xml/parsers/SAXParser.html#getParser--">getParser</a>
		// </span>
		if se, ok := t.(xml.StartElement); ok {
			switch {
			case se.Name.Local == "table" && hasAttr(se, "class", "memberSummary"):
				inMemberSummary = true
			case se.Name.Local == "span" && hasAttr(se, "class", "memberNameLink"):
				inMemberNameLink = true
			case inMemberSummary && inMemberNameLink && se.Name.Local == "a":
				// href="../../../javax/xml/parsers/SAXParser.html#getParser--"
				href, err := attr(se, "href")
				if err != nil {
					return entries // TODO: , err
				}

				// ds=["..", "..", "..", "javax", "xml", "parsers", "SAXParser.html#getParser--"]
				ds := strings.Split(href, "/")
				ns := []string{}
				lvls := 0
				for i, d := range ds {
					if i == len(ds)-1 {
						break // last one is the actual entity
					}
					if d == ".." {
						lvls++
						continue
					}
					ns = append(ns, d)
				}

				// "SAXParser.html#getParser--"
				last := ds[len(ds)-1]

				// "SAXParser"
				ent := last[0:strings.Index(last, ".")]

				// "getParser"
				fstart := strings.Index(last, "#") + 1
				fend := strings.Index(last, "-")
				if fend < 0 {
					fend = len(last)
				}
				fun := last[fstart:fend]

				// ["../", "../", "../", "javax/xml/parsers/SAXParser.html#getParser--"]
				// ls := strings.SplitAfterN(href, "/", lvls+1)
				// ["packs", "java", "docs", "api", "javax", "xml", "parsers", "SAXParser.html"]
				ps := strings.Split(path, string(os.PathSeparator))
				// ["..", "..", "..", "javax", "xml", "parsers", "SAXParser.html#getParser--"]
				hs := strings.Split(href, "/")
				// "SAXParser.html#getParser--"
				fl := hs[len(hs)-1]
				// "packs/java/docs/api/javax/xml/parsers/SAXParser.html#getParser--"
				tgt := strings.Join(append(ps[:len(ps)-1], fl), "/")

				e := entry{
					Namespace: ns,
					Entity:    ent,
					Member:    fun,
					Target:    tgt,
					Source:    path,
				}
				entries = append(entries, e)
			}
		}

		if se, ok := t.(xml.EndElement); ok {
			switch {
			case se.Name.Local == "span": // assumes no span nested in memberNameLink
				inMemberNameLink = false
			case se.Name.Local == "table": // assumes no table nested in memberSummary
				inMemberSummary = false
			}
		}
	}

	return entries
}
