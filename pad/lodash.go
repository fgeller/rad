package main

import (
	"../shared"

	"encoding/xml"
	"io"
	"log"
	"strings"
)

func parseLodashDocFile(filePath string, r io.Reader) []shared.Namespace {

	var nss []shared.Namespace
	d := xml.NewDecoder(r)
	d.Strict = false
	d.AutoClose = xml.HTMLAutoClose
	d.Entity = xml.HTMLEntity

	var t xml.Token
	var err error
	var inToc bool
	var inH2 bool
	var inCode bool
	var inLi bool
	var divCount int
	var group string

	for ; err == nil; t, err = d.Token() {
		se, gotStartElement := t.(xml.StartElement)
		ee, gotEndElement := t.(xml.EndElement)
		cd, gotChardata := t.(xml.CharData)

		switch {
		case gotChardata && inToc && inH2 && inCode:
			group += string(cd)
		case gotEndElement:
			switch {
			case ee.Name.Local == "div" && divCount == 0:
				inToc = false
			case ee.Name.Local == "div" && divCount > 0:
				divCount--
				group = ""
			case ee.Name.Local == "li":
				inLi = false
			case ee.Name.Local == "h2":
				inH2 = false
			case ee.Name.Local == "code":
				inCode = false
			}
		case gotStartElement:
			switch {
			case se.Name.Local == "div" && inToc:
				divCount++
			case se.Name.Local == "li":
				inLi = true
			case se.Name.Local == "div" && hasAttr(se, "class", "toc-container"):
				inToc = true
			case se.Name.Local == "h2":
				inH2 = true
			case se.Name.Local == "code":
				inCode = true
			case inToc && inLi && se.Name.Local == "a":
				// docs.html#support-ownLast
				href, _ := attr(se, "href") // TODO: err
				frag := href[strings.Index(href, "#")+1:]

				name := frag
				pth := group
				if strings.Index(frag, "-") > 0 {
					// ["support", "ownLast"]
					parts := strings.Split(frag, "-")
					pthParts := append([]string{group}, parts[:len(parts)-1]...)
					pth = strings.Join(pthParts, ".")
					name = parts[len(parts)-1]
				}

				tgt := filePath + "#" + frag
				ns := shared.Namespace{
					Path: pth,
					Members: []shared.Member{
						{
							Name:   name,
							Target: tgt,
						},
					},
				}
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
