package main

import (
	"../shared"

	"io"
	"log"
	"path/filepath"
	"regexp"

	"golang.org/x/net/html"
)

func parseManDocFile(filePath string, r io.Reader) []shared.Namespace {

	var nss []shared.Namespace

	bn := filepath.Base(filePath)
	dir := filepath.Dir(filePath)
	if bn != "dir_all_alphabetic.html" {
		log.Printf("Ignoring non-index %v\n", filePath)
		return nss
	}

	z := html.NewTokenizer(r)

	// <p><a id="letter_a" href="dir_all_alphabetic.html#top">top</a>
	// &nbsp; &nbsp; <a href="man3/a64l.3.html">a64l(3)</a> - convert between long and base-64
	// &nbsp; &nbsp; <a href="man3/abort.3.html">abort(3)</a> - cause abnormal process termination

	var err error
	var inA bool
	var href string
	var name string

processing:
	for {
		t := z.Next()
		switch {
		case t == html.ErrorToken:
			log.Printf("Finished parsing %v with err=%v\n", filePath, z.Err())
			break processing

		case t == html.TextToken && inA:
			name += string(z.Text())

		case t == html.EndTagToken:
			bn, _ := z.TagName()
			tn := string(bn)
			hrefMatches, err := regexp.Match("^man\\d/", []byte(href))
			switch {
			case tn == "a" && len(name) > 3 && hrefMatches && err == nil:
				inA = false
				p := href[:len("manD")]
				n := name[:len(name)-3]
				if n == "ChangeLog" {
					continue processing
				}
				tgt := dir + "/" + href
				ns := shared.Namespace{
					Path:    p,
					Members: []shared.Member{{Name: n, Target: tgt}},
				}
				nss = append(nss, ns)

				name = ""
				href = ""
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

			// <a href="man3/a64l.3.html">a64l(3)</a> - convert between long and base-64
			h, hasHref := attrs["href"]
			if tn == "a" && hasHref {
				matched, err := regexp.Match("^man\\d/", []byte(h))
				href = h
				inA = err == nil && matched
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
