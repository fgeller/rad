package main

import (
	"../shared"

	"io"
	"log"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func parseD3DocFile(filePath string, r io.Reader) []shared.Namespace {

	var nss []shared.Namespace

	if filepath.Base(filePath) != "API-Reference.html" {
		return nss
	}

	z := html.NewTokenizer(r)

	var err error
	var pastH2 bool
	beforeFooter := true
	var inUl bool
	var inLi bool
	var inA bool
	var href string
	var name string

processing:
	for {
		t := z.Next()
		switch {
		case t == html.ErrorToken:
			log.Printf("Finished parsing %v with err=%v found %v entries.\n", filePath, z.Err(), len(nss))
			break processing

		case t == html.TextToken && pastH2 && inUl && inLi && inA:
			name += string(z.Text())

		case t == html.EndTagToken:
			bn, _ := z.TagName()
			tn := string(bn)
			switch {
			case tn == "h2":
				pastH2 = true
			case tn == "ul":
				inUl = false
			case tn == "li":
				inLi = false
			case pastH2 && beforeFooter && inUl && inLi && tn == "a":
				inA = false

				pth := ""
				n := name
				li := strings.LastIndex(n, ".")
				if li > 0 {
					pth = n[:li]
					n = n[li+1:]
				}

				tgt := filepath.Dir(filePath) + "/" + href

				ns := shared.Namespace{
					Path:    pth,
					Members: []shared.Member{{Name: n, Target: tgt}},
				}
				nss = append(nss, ns)

				name = ""
				href = ""

				inLi = false // we only want the first link in a li
			}

		case t == html.StartTagToken:
			bn, hasAttrs := z.TagName()
			tn := string(bn)
			readAttrs := func() map[string]string {
				as := map[string]string{}
				if !hasAttrs {
					return as
				}
				for {
					k, v, more := z.TagAttr()
					as[string(k)] = string(v)
					if !more {
						return as
					}
				}
			}
			attrs := readAttrs()

			// <h3>
			// <a id="user-content-selections" class="anchor" href="API-Reference.html#selections" aria-hidden="true"><span class="octicon octicon-link"></span></a><a href="Selections.html">Selections</a>
			// </h3>

			// <ul>
			// <li>
			// <a href="Selections.html#d3_event">d3.event</a> - access the current user event for interaction.</li>
			// <li>
			switch {
			case tn == "ul":
				inUl = true
			case tn == "li":
				inLi = true
			case tn == "div":
				class, _ := attrs["class"]
				if class == "site-footer" {
					beforeFooter = false
				}
			case pastH2 && beforeFooter && inUl && inLi && tn == "a":
				href, _ = attrs["href"]
				inA = true
			}

		default:
		}

	}

	if len(nss) == 0 {
		log.Printf("Found no entries for %v, err=%v\n", filePath, err)
	}

	rs := shared.Merge(nss)
	return rs
}
