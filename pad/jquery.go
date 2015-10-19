package main

import (
	"../shared"

	"encoding/xml"
	"io"
	"log"
	// "path/filepath"
	"strings"
)

func parseJQueryDocFile(filePath string, r io.Reader) []shared.Namespace {

	var nss []shared.Namespace
	d := xml.NewDecoder(r)
	d.Strict = false
	d.AutoClose = xml.HTMLAutoClose
	d.Entity = xml.HTMLEntity

	var t xml.Token
	var err error
	var inSignature bool
	var inA bool
	var inIconLink bool
	var afterIconLink bool
	var name string
	var href string

	//
	// <li class="signature">
	// <h4 class="name">
	// <span class="version-details">version added: <a href="../category/version/1.5/index.html">1.5</a></span>
	// <a id="jQuery-ajax-url-settings" href="index.html#jQuery-ajax-url-settings">
	//   <span class="icon-link"></span>jQuery.ajax( url [, settings ] )
	// </a>
	//

	for ; err == nil; t, err = d.Token() {
		se, gotStartElement := t.(xml.StartElement)
		ee, gotEndElement := t.(xml.EndElement)
		cd, gotCharData := t.(xml.CharData)

		switch {
		case gotCharData && inSignature && inA && afterIconLink:
			name += string(cd)

		case gotEndElement:
			switch {
			case ee.Name.Local == "li":
				inSignature = false
			case inSignature && afterIconLink && ee.Name.Local == "a":
				inA = false
				afterIconLink = false

				// name [jQuery.ajax( url [, settings ] )]
				// href [index.html#jQuery-ajax-url-settings]
				tgt := filePath + href[strings.Index(href, "#"):]
				nm := "name"
				pth := ""
				if strings.Index(name, ".") > 0 {
					pth = name[:strings.Index(name, ".")]
					nm = name[strings.Index(name, ".")+1:]
					nm = strings.Replace(nm, " ", "", -1)
					nm = strings.Replace(nm, "[", "", -1)
					nm = strings.Replace(nm, "]", "", -1)
				}

				ns := shared.Namespace{
					Path:    pth,
					Members: []shared.Member{{Name: nm, Target: tgt}},
				}
				nss = append(nss, ns)
				name = ""
				href = ""
			case inSignature && inIconLink && ee.Name.Local == "span":
				inIconLink = false
				afterIconLink = true
			}

		case gotStartElement:
			switch {
			case se.Name.Local == "li" && hasAttr(se, "class", "signature"):
				inSignature = true
			case inSignature && se.Name.Local == "a":
				inA = true
				href, _ = attr(se, "href")
			case inSignature && inA && se.Name.Local == "span" && hasAttr(se, "class", "icon-link"):
				inIconLink = true
			}
		}
	}

	if len(nss) == 0 {
		log.Printf("Found no entries for %v, err=%v\n", filePath, err)
	}

	rs := shared.Merge(nss)
	return rs
}
