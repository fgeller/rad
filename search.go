package main

import (
	"strings"
)

func iPrefix(s string, pfx string) bool {
	return strings.HasPrefix(strings.ToLower(s), strings.ToLower(pfx))
}

func findEntityMember(pack string, entity string, fun string, limit int) ([]entry, error) {
	results := []entry{}

	for packName, es := range docs {
		if iPrefix(packName, pack) {
			for _, e := range es {
				if iPrefix(e.Name, entity) {
					for _, m := range e.Members {
						if iPrefix(m.Name, fun) {
							results = append(results, e)
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
