package main

import (
	"../shared"

	"io"
	"log"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func parseMDNDocFile(filePath string, r io.Reader) []shared.Namespace {

	var nss []shared.Namespace
	z := html.NewTokenizer(r)

	// <dt>
	// 	<a href="/en-US/docs/Web/JavaScript/Reference/Global_Objects/String/fromCharCode"
	//     title="The static String.fromCharCode() method returns a string created by using the specified sequence of Unicode values.">
	// 	  <code>String.fromCharCode()</code>
	//   </a>
	// </dt>

	var err error
	var inMethods bool
	var inProperties bool
	var inDt bool
	var inA bool
	var inCode bool
	var href string
	var name string

processing:
	for {
		t := z.Next()
		switch {
		case t == html.ErrorToken:
			log.Printf("Finished parsing %v with err=%v\n", filePath, z.Err())
			break processing

		case t == html.TextToken && (inMethods || inProperties) && inDt && inA && inCode:
			name += string(z.Text())

		case t == html.EndTagToken:
			bn, _ := z.TagName()
			tn := string(bn)
			switch {
			case tn == "dt":
				inDt = false
			case tn == "a":
				inA = false
			case (inMethods || inProperties) && inDt && inA && tn == "code":
				inCode = false
				pth := ""
				tgt := filepath.Join(filepath.Dir(filePath), href)
				n := ""

				parts := strings.Split(name, ".")
				pth = strings.Join(parts[:len(parts)-1], ".")
				n = parts[len(parts)-1]

				ns := shared.Namespace{
					Path:    pth,
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

			switch {
			case tn == "h2" || tn == "h3":
				id, ok := attrs["id"]
				inMethods = ok &&
					(strings.HasPrefix(id, "Methods") ||
						strings.HasPrefix(id, "Properties"))

			case (inMethods || inProperties) && tn == "dt":
				inDt = true
			case (inMethods || inProperties) && inDt && tn == "a":
				href, _ = attrs["href"]
				inA = true
			case (inMethods || inProperties) && inA && tn == "code":
				inCode = true
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
