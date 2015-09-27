package main

import (
	"../shared"
	"encoding/xml"
	"io"
	"log"
	"regexp"
	"strings"
)

func parseScalaDocFile(f string, r io.Reader) []shared.Entry {
	d := xml.NewDecoder(r)
	var t xml.Token
	var err error
	entries := []shared.Entry{}

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

func parseEntry(source string, target string, s string) (shared.Entry, error) {
	e := shared.Entry{}

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
	sigIdx := strings.IndexAny(meth, ":[(=")
	extendsIdx := strings.Index(meth, "extends")

	// want: smaller one that's larger than 0
	if extendsIdx > 0 && sigIdx > 0 && extendsIdx < sigIdx {
		sigIdx = extendsIdx
	}
	if extendsIdx > 0 && sigIdx < 0 {
		sigIdx = extendsIdx
	}

	m := meth
	signature := ""
	if sigIdx > 0 {
		m = meth[:sigIdx]
		signature = meth[sigIdx:]
	}

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

	e.Namespace = namespace
	e.Name = entity
	e.Members = []shared.Member{{Name: m, Signature: signature, Target: newTarget}}

	return e, nil
}
