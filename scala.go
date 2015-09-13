package main

import (
	"encoding/xml"
	"io"
	"log"
	"regexp"
	"strings"
)

func indexScalaApi(packName string) func() ([]entry, error) {
	return func() ([]entry, error) {
		path := packDir + "/" + packName
		log.Printf("About to index scala api in [%v]\n", path)
		return scan(path, parseScalaDocFile)
	}
}

func parseScalaDocFile(f string, r io.Reader) []entry {
	d := xml.NewDecoder(r)
	var t xml.Token
	var err error
	entries := []entry{}

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
						e, err := parseEntry(f, subs[0], subs[1])
						if err != nil {
							log.Println(err)
						} else {
							entries = append(entries, e)
						}

					}
				}
			}
		}
	}

	return entries
}

func parseEntry(source string, target string, s string) (entry, error) {
	e := entry{Source: source}

	// ns0.ns1.ns2.e1$entity @ method
	splits := strings.Split(s, "@")
	fqEnt := splits[0]
	meth := ""
	if len(splits) > 0 {
		meth = strings.Join(splits[1:], "")
	}

	// ns0.ns1.ns2 . e1$entity
	entPat, err := regexp.Compile("(.+)\\.(.+)")
	if err != nil {
		return e, err
	}

	// [[ns0.ns1.ns2.e1$entity ns0.ns1.ns2 e1$entity]]
	ms := entPat.FindAllStringSubmatch(fqEnt, -1)

	// [ns0 ns1 ns2]
	namespace := []string{}
	entity := ""
	if len(ms) == 0 {
		entity = fqEnt
	} else {
		namespace = strings.Split(ms[0][1], ".")
		// [e1 entity]
		obj := strings.Split(ms[0][2], "$")
		for i := len(obj) - 1; i >= 0; i-- {
			if len(obj[i]) > 0 {
				entity = obj[i]
				for _, e := range obj[:i] {
					namespace = append(namespace, e)
				}
				break
			}
		}
	}

	// name[A](...
	// name(...
	sigIdx := strings.IndexAny(meth, ":[(")
	member := meth
	signature := ""
	if sigIdx > 0 {
		member = meth[:sigIdx]
		signature = meth[sigIdx:]
	}

	e.Namespace = namespace
	e.Entity = entity
	e.Member = member
	e.Signature = signature

	// find target link
	targetSplits := strings.Split(target, "/")
	upCount := 0
	for i, v := range targetSplits {
		if v != ".." {
			upCount = i
			break
		}
	}

	// TODO: clarify this a bit
	sourceSplits := strings.Split(source, "/")
	newSplits := sourceSplits[:len(sourceSplits)-(upCount+1)]
	newSplits = append(newSplits, targetSplits[len(targetSplits)-1]+s)
	newTarget := strings.Join(newSplits, "/")
	e.Target = newTarget

	return e, nil
}
