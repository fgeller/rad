package main

import (
	"../shared"

	"encoding/xml"
	"io"
	"log"
	"strings"
)

func parseNodeJsDocFile(filePath string, r io.Reader) []shared.Namespace {

	var nss []shared.Namespace
	d := xml.NewDecoder(r)
	d.Strict = false
	d.AutoClose = xml.HTMLAutoClose
	d.Entity = xml.HTMLEntity

	//       <div id="toc">
	//         <h2>Table of Contents</h2>
	//         <ul>
	// <li><a href="console.html#console_console">Console</a><ul>
	// <li><a href="console.html#console_console_1">console</a><ul>
	// <li><a href="console.html#console_console_log_data">console.log([data][, ...])</a></li>
	// <li><a href="console.html#console_console_info_data">console.info([data][, ...])</a></li>
	// <li><a href="console.html#console_console_error_data">console.error([data][, ...])</a></li>
	// <li><a href="console.html#console_console_warn_data">console.warn([data][, ...])</a></li>
	// <li><a href="console.html#console_console_dir_obj_options">console.dir(obj[, options])</a></li>
	// <li><a href="console.html#console_console_time_label">console.time(label)</a></li>
	// <li><a href="console.html#console_console_timeend_label">console.timeEnd(label)</a></li>
	// <li><a href="console.html#console_console_trace_message">console.trace(message[, ...])</a></li>
	// <li><a href="console.html#console_console_assert_value_message">console.assert(value[, message][, ...])</a></li>
	// </ul>
	// </li>
	// <li><a href="console.html#console_class_console">Class: Console</a><ul>
	// <li><a href="console.html#console_new_console_stdout_stderr">new Console(stdout[, stderr])</a></li>
	// </ul>
	// </li>
	// </ul>
	// </li>
	// </ul>

	var t xml.Token
	var err error
	var inToc bool
	var inA bool
	var href string
	var name string
	var module string

	for ; err == nil; t, err = d.Token() {
		se, gotStartElement := t.(xml.StartElement)
		ee, gotEndElement := t.(xml.EndElement)
		cd, gotCharData := t.(xml.CharData)

		switch {
		case gotCharData && inToc && inA:
			name += string(cd)

		case gotEndElement:
			switch {
			case ee.Name.Local == "div":
				inToc = false
			case inToc && ee.Name.Local == "a":
				inA = false

				name = strings.Replace(name, "[", "", -1)
				name = strings.Replace(name, "]", "", -1)

				pth := module
				n := name
				if strings.HasPrefix(n, "Class: ") {
					n = n[len("Class: "):]
				}

				if strings.Index(name, ".") > 0 {
					pth += "." + name[:strings.Index(name, ".")]
					n = name[strings.Index(name, ".")+1:]
				}

				tgt := filePath + href[strings.Index(href, "#"):]

				ns := shared.Namespace{
					Path:    pth,
					Members: []shared.Member{{Name: n, Target: tgt}},
				}
				nss = append(nss, ns)

				if module == "" {
					module = name
				}

				name = ""
				href = ""
			}

		case gotStartElement:
			switch {
			case se.Name.Local == "div" && hasAttr(se, "id", "toc"):
				inToc = true
			case inToc && se.Name.Local == "a":
				href, _ = attr(se, "href")
				inA = true
			}
		}
	}

	if len(nss) == 0 {
		log.Printf("Found no entries for %v, err=%v\n", filePath, err)
	}

	rs := shared.Merge(nss)
	return rs
}
