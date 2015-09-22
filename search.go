package main

import (
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
	if len(s.Namespace) != len(o.Namespace) {
		return false
	}

	for i := range s.Namespace {
		if s.Namespace[i] != o.Namespace[i] {
			return false
		}
	}

	return s.Entity == o.Entity &&
		s.Member == o.Member &&
		s.Signature == o.Signature &&
		s.Target == o.Target &&
		s.Source == o.Source
}

// TODO: always a member available?
func NewSearchResult(e entry, memberIdx int) searchResult {
	return searchResult{
		Entity:    e.Name,
		Namespace: e.Namespace,
		Member:    e.Members[memberIdx].Name,
		Signature: e.Members[memberIdx].Signature,
		Target:    e.Members[memberIdx].Target,
		Source:    e.Members[memberIdx].Source,
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
