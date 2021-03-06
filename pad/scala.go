package main

import (
	"../shared"
	"encoding/xml"
	"io"
	"log"
	"regexp"
	"strings"
)

func parseScalaDocFile(f string, r io.Reader) []shared.Namespace {
	d := xml.NewDecoder(r)
	var t xml.Token
	var err error
	namespaces := []shared.Namespace{}

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
						n, err := parseNamespace(f, subs[0], subs[1])
						if err != nil {
							log.Println(err)
						} else {
							namespaces = append(namespaces, n)
						}
					}
				}
			}
		}
	}

	return namespaces
}

func parseNamespace(source string, target string, s string) (shared.Namespace, error) {
	namespace := shared.Namespace{}

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
		return namespace, err
	}

	// [[ns0.ns1.ns2.e1$entity ns0.ns1.ns2 e1$entity]]
	ms := entPat.FindAllStringSubmatch(fqEnt, -1)

	// [ns0 ns1 ns2]
	path := []string{}
	entity := ""
	if len(ms) == 0 {
		entity = fqEnt
	} else {
		path = strings.Split(ms[0][1], ".")
		// [e1 entity]
		obj := strings.Split(ms[0][2], "$")
		for i := len(obj) - 1; i >= 0; i-- {
			if len(obj[i]) > 0 {
				entity = obj[i]
				for _, p := range obj[:i] {
					path = append(path, p)
				}
				break
			}
		}
	}

	// name[A](...
	// name(...
	sigIdx := -1
	if len(meth) > 0 {
		sigIdx = strings.IndexAny(meth[1:], ":[(")
		if sigIdx >= 0 {
			sigIdx++
		}
	}
	tpeAlias := regexp.MustCompile("=[[:alpha:]]")
	tpeAliasIdx := tpeAlias.FindStringIndex(meth)
	if tpeAliasIdx != nil && (tpeAliasIdx[0] < sigIdx || sigIdx <= 0) {
		sigIdx = tpeAliasIdx[0]
	}
	extendsIdx := strings.Index(meth, "extends")

	// want: smaller one that's larger than 0
	if extendsIdx > 0 && sigIdx > 0 && extendsIdx < sigIdx {
		sigIdx = extendsIdx
	}
	if extendsIdx > 0 && sigIdx < 0 {
		sigIdx = extendsIdx
	}

	m := meth
	// signature := ""
	if sigIdx > 0 {
		m = meth[:sigIdx]
		// signature = meth[sigIdx:]
	}
	if len(m) == 0 {
		m = entity
	} else {
		path = append(path, entity)
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

	namespace.Path = strings.Join(path, ".")
	namespace.Members = []shared.Member{{Name: m, Target: newTarget}}

	return namespace, nil
}
