package main

import (
	"encoding/xml"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func indexJavaApi(packName string) func() ([]entry, error) {
	return func() ([]entry, error) {
		path := packDir + "/" + packName
		log.Printf("About to index java api in [%v]\n", path)
		return scan(path, parseJavaDocFile)
	}
}

func parseHref(href string, path string) entry {
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

	return entry{
		Namespace: ns,
		Entity:    ent,
		Member:    fun,
		Target:    tgt,
		Source:    path,
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
	var inInheritedBlock bool
	entries := []entry{}
	inheritedBlockPattern, _ := regexp.Compile("^(methods|fields)\\.inherited.+") // TODO: err

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
			case inInheritedBlock && se.Name.Local == "a":
				// TODO: should be matching href
				href, err := attr(se, "href")
				if err != nil {
					return entries // TODO: , err
				}

				e := parseHref(href, path)
				entries = append(entries, e)

			case inMemberSummary && inMemberNameLink && se.Name.Local == "a":
				// href="../../../javax/xml/parsers/SAXParser.html#getParser--"
				href, err := attr(se, "href")
				if err != nil {
					return entries // TODO: , err
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

	return entries
}
