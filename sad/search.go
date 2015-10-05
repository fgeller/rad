package main

import (
	"log"
	"strings"
)

func find(results chan searchResult, end chan bool, params searchParams) {
	defer func() { end <- true }()
	for pack, namespaces := range docs {
		if params.pack.MatchString(pack) {
			for _, namespace := range namespaces {
				normalizedPath := strings.Join(namespace.Path, ".")
				if params.path.MatchString(normalizedPath) {
					for mi, member := range namespace.Members {
						if params.member.MatchString(member.Name) {
							select {
							case <-end:
								log.Printf("Got poison pill in find, stopping search.\n")
								return
							case results <- NewSearchResult(namespace, mi):
							}
						}
					}
				}
			}
		}
	}
}
