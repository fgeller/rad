package main

import (
	"../shared"

	"encoding/xml"
	"io"
	"log"
	"strings"
)

func parsePy27DocFile(filePath string, r io.Reader) []shared.Namespace {

	var nss []shared.Namespace
	// TODO: pull this out
	d := xml.NewDecoder(r)
	d.Strict = false
	d.AutoClose = xml.HTMLAutoClose
	d.Entity = xml.HTMLEntity

	var t xml.Token
	var err error

	for ; err == nil; t, err = d.Token() {
		se, gotStartElement := t.(xml.StartElement)

		switch {
		case gotStartElement:
			switch {
			case se.Name.Local == "a" && hasAttr(se, "class", "headerlink"):
				href, _ := attr(se, "href") // TODO: err
				tgt := filePath + href
				ns := parseHref(href, tgt)
				nss = append(nss, ns)
			}
		}
	}

	if len(nss) == 0 {
		log.Printf("Found no entries for %v\n", filePath)
	}

	rs := shared.Merge(nss)
	return rs
}

func parseHref(href string, target string) shared.Namespace {

	var ns shared.Namespace

	// href = #datetime.timedelta.total_seconds
	// datetime.timedelta.total_seconds
	h := href[1:]

	// strips prefix for #module-datetime link
	if strings.HasPrefix(h, "module-") {
		h = h[len("module-"):]
	}

	// ["datetime" "timedelta" "total_seconds"]
	parts := strings.Split(h, ".")
	pth := strings.Join(parts[:len(parts)-1], ".")
	n := parts[len(parts)-1]
	ns.Path = pth
	ns.Members = []shared.Member{{Name: n, Target: target}}
	return ns
}
