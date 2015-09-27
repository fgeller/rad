package main

import (
	"encoding/xml"
	"fmt"
	"regexp"
)

func attr(se xml.StartElement, name string) (string, error) {
	for _, att := range se.Attr {
		if att.Name.Local == name {
			return att.Value, nil
		}
	}

	return "", fmt.Errorf("could not find attr %v", name)
}

func hasAttr(se xml.StartElement, name string, value string) bool {
	v, err := attr(se, name)
	return err == nil && v == value
}

func hasAttrMatches(se xml.StartElement, name string, pattern *regexp.Regexp) bool {
	v, err := attr(se, name)
	return err == nil && pattern.MatchString(v)
}
