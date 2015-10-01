package main

import (
	"../shared"
	"reflect"
	"strings"
)

type searchResult struct {
	Entity    string
	Namespace []string
	Member    string
	Signature string
	Target    string
	Source    string
}

func (s searchResult) eq(o searchResult) bool {
	return reflect.DeepEqual(s, o)
}

func NewSearchResult(e shared.Entry, memberIdx int) searchResult {

	if len(e.Members) == 0 {
		return searchResult{
			Entity:    e.Name,
			Namespace: e.Namespace,
		}
	}

	return searchResult{
		Entity:    e.Name,
		Namespace: e.Namespace,
		Member:    e.Members[memberIdx].Name,
		Signature: e.Members[memberIdx].Signature,
		Target:    "/pack/" + e.Members[memberIdx].Target,
	}
}

func iPrefix(s string, pfx string) bool {
	return strings.HasPrefix(strings.ToLower(s), strings.ToLower(pfx))
}

func findEntityMember(pack string, entity string, fun string, limit int) ([]searchResult, error) {
	results := []searchResult{}

	for packName, es := range docs {
		if iPrefix(packName, pack) {
			for _, e := range es {
				if iPrefix(e.Name, entity) {
					for mi, m := range e.Members {
						if iPrefix(m.Name, fun) {
							results = append(results, NewSearchResult(e, mi))
							if len(results) == limit {
								return results, nil
							}
						}
					}
				}
			}
		}
	}

	return results, nil
}
