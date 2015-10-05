package main

import (
	"../shared"

	"log"
	"math"
	"runtime"
	"strings"
)

func findNamespace(
	results chan searchResult,
	end chan bool,
	namespaces []shared.Namespace,
	params searchParams) {
	defer func() { end <- true }()

	for _, namespace := range namespaces {
		normalizedPath := strings.Join(namespace.Path, ".")
		if params.path.MatchString(normalizedPath) {
			for mi, member := range namespace.Members {
				if params.member.MatchString(member.Name) {
					select {
					case <-end:
						log.Printf("Got poison pill in findNamespace, stopping search.\n")
						return
					case results <- NewSearchResult(namespace, mi):
					}
				}
			}
		}
	}
}

func find(results chan searchResult, end chan bool, params searchParams) {
	for pack, namespaces := range docs {
		if params.pack.MatchString(pack) {
			cpus := int64(runtime.NumCPU())

			if int64(len(namespaces)) < cpus {
				findNamespace(results, end, namespaces, params)
				end <- true
				return
			}

			// example for cpus = 4; len(namespaces) = 11
			// partitionSize = ceil(11 / 4) = 3
			// 0:4  - 0 1 2 3
			// 4:8  - 4 5 6 7
			// 8:12 - 8 9 10

			ps := int64(math.Ceil(float64(len(namespaces)) / float64(cpus)))
			for i := int64(0); i < cpus; i++ {
				ns := namespaces[i*ps : (i+1)*ps]
				go findNamespace(results, end, ns, params)
			}
		}
	}
}
