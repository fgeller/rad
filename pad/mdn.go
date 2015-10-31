package main

import (
	"../shared"

	"encoding/xml"
	"io"
	"log"
	"path/filepath"
	"strings"
)

func parseMDNDocFile(filePath string, r io.Reader) []shared.Namespace {

	var nss []shared.Namespace
	d := xml.NewDecoder(r)
	d.Strict = false
	d.AutoClose = xml.HTMLAutoClose
	d.Entity = xml.HTMLEntity

	var t xml.Token
	var err error
	var inMethods bool
	var inProperties bool
	var inDt bool
	var inA bool
	var inCode bool
	var href string
	var name string

	// <dt>
	// 	<a href="/en-US/docs/Web/JavaScript/Reference/Global_Objects/String/fromCharCode"
	//     title="The static String.fromCharCode() method returns a string created by using the specified sequence of Unicode values.">
	// 	  <code>String.fromCharCode()</code>
	//   </a>
	// </dt>

	for ; err == nil; t, err = d.Token() {
		se, gotStartElement := t.(xml.StartElement)
		ee, gotEndElement := t.(xml.EndElement)
		cd, gotCharData := t.(xml.CharData)

		switch {
		case gotCharData && (inMethods || inProperties) && inDt && inA && inCode:
			name += string(cd)

		case gotEndElement:
			switch {
			case ee.Name.Local == "dt":
				inDt = false
			case ee.Name.Local == "a":
				inA = false
			case (inMethods || inProperties) && inDt && inA && ee.Name.Local == "code":
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

			// <h2 id="Methods">Methods</h2>
			// <h3 id="Methods_2">Methods</h3>

		case gotStartElement:
			switch {
			case se.Name.Local == "h2" || se.Name.Local == "h3":
				id, err := attr(se, "id")
				inMethods = err == nil &&
					(strings.HasPrefix(id, "Methods") ||
						strings.HasPrefix(id, "Properties"))
			case (inMethods || inProperties) && se.Name.Local == "dt":
				inDt = true
			case (inMethods || inProperties) && inDt && se.Name.Local == "a":
				href, _ = attr(se, "href")
				inA = true
			case (inMethods || inProperties) && inA && se.Name.Local == "code":
				inCode = true
			}
		}
	}

	if len(nss) == 0 {
		log.Printf("Found no entries for %v, err=%v\n", filePath, err)
	}

	rs := shared.Merge(nss)
	return rs
}
