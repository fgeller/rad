package main

import (
	"../shared"

	"encoding/xml"
	"io"
	"log"
	"strings"
)

func parseDjangoDocFile(filePath string, r io.Reader) []shared.Namespace {

	var nss []shared.Namespace
	// TODO: pull this out
	d := xml.NewDecoder(r)
	d.Strict = false
	d.AutoClose = xml.HTMLAutoClose
	d.Entity = xml.HTMLEntity

	var t xml.Token
	var err error
	var inDt bool
	var inH1 bool

	for ; err == nil; t, err = d.Token() {
		se, gotStartElement := t.(xml.StartElement)
		ee, gotEndElement := t.(xml.EndElement)

		switch {
		case gotEndElement:
			switch {
			case ee.Name.Local == "dt":
				inDt = false
			case ee.Name.Local == "h1":
				inH1 = false
			}
		case gotStartElement:
			switch {
			case se.Name.Local == "dt":
				inDt = true
			case se.Name.Local == "h1":
				inH1 = true
			case (inDt || inH1) && se.Name.Local == "a" && hasAttr(se, "class", "headerlink"):
				// request-response.html#django.http.HttpRequest
				href, _ := attr(se, "href") // TODO: err
				frag := href[strings.Index(href, "#"):]
				tgt := filePath + frag
				ns := parseFrag(frag, tgt)
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

func parseFrag(frag string, target string) shared.Namespace {

	var ns shared.Namespace

	// frag = #django.http.HttpRequest
	// f = django.http.HttpRequest
	f := frag[1:]

	// ["django" "http" "HttpRequest"]
	parts := strings.Split(f, ".")
	pth := strings.Join(parts[:len(parts)-1], ".")
	n := parts[len(parts)-1]
	ns.Path = pth
	ns.Members = []shared.Member{{Name: n, Target: target}}
	return ns
}
