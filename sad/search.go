package main

import (
	"../shared"
	"reflect"
	"strings"
)

type searchResult struct {
	Namespace []string
	Member    string
	Target    string
}

func (s searchResult) eq(o searchResult) bool {
	return reflect.DeepEqual(s, o)
}

func NewSearchResult(n shared.Namespace, memberIdx int) searchResult {

	if len(n.Members) == 0 { // TODO: do we need this guy?
		return searchResult{
			Namespace: n.Path,
		}
	}

	return searchResult{
		Namespace: n.Path,
		Member:    n.Members[memberIdx].Name,
		Target:    "/pack/" + n.Members[memberIdx].Target, // TODO: should we fix that here?
	}
}

func iPrefix(s string, pfx string) bool {
	return strings.HasPrefix(strings.ToLower(s), strings.ToLower(pfx))
}

func findEntityMember(pack string, entity string, fun string, limit int) ([]searchResult, error) {
	results := []searchResult{}

	for packName, ns := range docs {
		if iPrefix(packName, pack) {
			for _, n := range ns {
				if iPrefix(n.Last(), entity) {
					for mi, m := range n.Members {
						if iPrefix(m.Name, fun) {
							results = append(results, NewSearchResult(n, mi))
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
