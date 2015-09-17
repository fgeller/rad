package main

import (
	"fmt"
	"strings"
)

func findEntityMember(pack string, entity string, fun string, limit int) ([]entry, error) {
	es, ok := docs[pack]
	if !ok {
		return es, fmt.Errorf("Package [%v] not installed.", pack)
	}

	results := []entry{}

	for _, e := range es {
		if strings.HasPrefix(strings.ToLower(e.Name), strings.ToLower(entity)) {
			for _, m := range e.Members {
				if strings.HasPrefix(strings.ToLower(m.Name), strings.ToLower(fun)) {
					results = append(results, e)
					if len(results) == limit {
						return results, nil
					}
				}
			}
		}
	}

	return results, nil
}
