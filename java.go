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

	return entry{
		Namespace: ns,
		Name:      ent,
		Members:   []member{member{Name: fun, Target: tgt, Source: path}},
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
	var inheritedBlock string
	entries := []entry{}
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

	return mergeEntries(entries)
}

func isSameEntry(a entry, b entry) bool {
	if a.Name != b.Name ||
		len(a.Namespace) != len(b.Namespace) {
		return false
	}
	for i := range a.Namespace {
		if a.Namespace[i] != b.Namespace[i] {
			return false
		}
	}

	return true
}

func mergeEntries(entries []entry) []entry {
	if len(entries) < 1 {
		return entries
	}

	unmerged := entries[1:]
	merged := []entry{entries[0]}

merging:
	for ui := range unmerged {
		for mi := range merged {
			if isSameEntry(unmerged[ui], merged[mi]) {
				merged[mi].Members = append(merged[mi].Members, unmerged[ui].Members...)
				continue merging
			}
		}
		merged = append(merged, unmerged[ui])
	}

	return merged
}
