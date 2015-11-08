package main

import (
	"../shared"

	"io"
	"log"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func parsePHPDocFile(filePath string, r io.Reader) []shared.Namespace {

	var nss []shared.Namespace

	fn := filepath.Base(filePath)
	if fn != "indexes.functions.php.html" {
		return nss
	}
	z := html.NewTokenizer(r)

	// <ul class='gen-index index-for-refentry'><li class='gen-index index-for-a'>a<ul id='refentry-index-for-a'>
	// <li><a href="function.abs.php.html" class="index">abs</a> - Absolute value</li>
	// <li><a href="function.json-decode.php.html" class="index">json_decode</a> - Decodes a JSON string</li>
	// <li><a href="judy.bycount.php.html" class="index">Judy::byCount</a> - Locate the Nth index present in the Judy array</li>
	// <li><a href="zmqsocket.construct.php.html" class="index">ZMQSocket::__construct</a> - Construct a new ZMQSocket</li>
	// </ul></li>
	// </ul>

	var err error
	var inRefEntry bool
	var inLi bool
	var inA bool
	var ulCount int
	var href string
	var name string

processing:
	for {
		t := z.Next()
		switch {
		case t == html.ErrorToken:
			log.Printf("Finished parsing %v with err=%v\n", filePath, z.Err())
			break processing

		case t == html.TextToken && inRefEntry && inLi && inA:
			name += string(z.Text())

		case t == html.EndTagToken:
			bn, _ := z.TagName()
			tn := string(bn)
			switch {
			case tn == "ul" && ulCount > 0:
				ulCount--
			case tn == "ul" && ulCount == 0:
				inRefEntry = false
			case tn == "li":
				inLi = false
			case inRefEntry && inLi && inA && tn == "a":
				inA = false
				tgt := filepath.Dir(filePath) + "/" + href
				n := name
				pth := ""
				if strings.Index(n, "::") > 0 {
					pth = n[:strings.Index(n, "::")]
					n = n[strings.Index(n, "::")+2:]
				}
				ns := shared.Namespace{
					Path: pth,
					Members: []shared.Member{
						{
							Name:   n,
							Target: tgt,
						},
					},
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

			switch {
			case tn == "ul" && !inRefEntry:
				class, ok := attrs["class"]
				inRefEntry = ok && strings.Index(class, "index-for-refentry") >= 0

			case tn == "ul" && inRefEntry:
				ulCount++
			case tn == "li" && inRefEntry:
				inLi = true
			case tn == "a" && inRefEntry && inLi:
				inA = true
				h, ok := attrs["href"]
				if ok {
					href = h
				}
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
