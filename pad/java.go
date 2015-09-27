package main

import (
	"../shared"
	"encoding/xml"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func parseJavaDocFile(path string, r io.Reader) []shared.Entry {

	d := xml.NewDecoder(r)
	d.Strict = false
	d.AutoClose = xml.HTMLAutoClose
	d.Entity = xml.HTMLEntity

	var t xml.Token
	var err error
	var inMemberNameLink bool
	var inMemberSummary bool
	var inInheritedBlock bool
	var inheritedBlock string
	entries := []shared.Entry{}
	inheritedBlockPattern, err := regexp.Compile("^(methods|fields)\\.inherited.+")
	if err != nil {
		log.Fatal("Can't compile pattern for Java doc parsing: %v\n", err)
	}

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
			case se.Name.Local == "a" && hasAttrMatches(se, "name", inheritedBlockPattern):
				inInheritedBlock = true
				inheritedBlock, err = attr(se, "name")
				if err != nil {
					log.Fatal("Unexpected error while accessing attr 'name': %v\n", err)
				}

			case inInheritedBlock && se.Name.Local == "a":
				href, err := attr(se, "href")
				if err != nil {
					log.Fatal("Unexpected error while accessing attr 'href: %v\n", err)
				}

				// ["testdata", "ActionEvent.html"]
				ps := strings.Split(path, string(os.PathSeparator))
				// "ActionEvent.html"
				ef := ps[len(ps)-1]
				// "ActionEvent"
				ent := ef[:strings.Index(ef, ".")]

				// testdata/ActionEvent.html#inheritedBlock
				tgt := strings.Join(ps, "/") + "#" + inheritedBlock

				e := parseHref(href, path)
				e.Name = ent
				for i := range e.Members {
					e.Members[i].Target = tgt
				}
				entries = append(entries, e)

			case inMemberSummary && inMemberNameLink && se.Name.Local == "a":
				// href="../../../javax/xml/parsers/SAXParser.html#getParser--"
				href, err := attr(se, "href")
				if err != nil {
					log.Fatal("Unexpected error while accessing attr 'href': %v\n", err)
				}

				e := parseHref(href, path)
				entries = append(entries, e)
			}
		}

		if se, ok := t.(xml.EndElement); ok {
			switch {
			case se.Name.Local == "span": // assumes no span nested in memberNameLink
				inMemberNameLink = false
			case se.Name.Local == "table": // assumes no table nested in memberSummary
				inMemberSummary = false
			case se.Name.Local == "li":
				inInheritedBlock = false
			}
		}
	}

	return shared.MergeEntries(entries)

}

func parseHref(href string, path string) shared.Entry {
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
	fend := strings.Index(last[fstart:], "-")
	if fend < 0 {
		fend = len(last)
	} else {
		fend += fstart
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

	return shared.Entry{
		Namespace: ns,
		Name:      ent,
		Members: []shared.Member{
			shared.Member{
				Name:   fun,
				Target: tgt,
			},
		},
	}
}
