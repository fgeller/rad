package main

import (
	"../shared"

	"log"
	"math"
	"runtime"
	"sync"
)

func findNamespace(
	results chan searchResult,
	end chan struct{},
	namespaces []shared.Namespace,
	params searchParams) {

	for _, namespace := range namespaces {
		if params.path.MatchString(namespace.Path) {
			for mi, member := range namespace.Members {
				if len(end) > 0 {
					log.Printf("Got poison pill in findNamespace, stopping search.\n")
					return
				}

				if params.member.MatchString(member.Name) {
					select {
					case results <- NewSearchResult(namespace, mi):
					case <-end:
						end <- struct{}{}
						log.Printf("Got poison pill in findNamespace, stopping search.\n")
						return
					}
				}
			}
		}
	}
}

// expects end to be buffered
// will close results channel when done.
func find(results chan searchResult, end chan struct{}, params searchParams) {

	res := make(chan packResp)
	req := packReq{tpe: Read, res: res}
	global.packs <- req

	for {
		select {
		case resp, ok := <-res:
			if !ok {
				close(results)
				return
			}

			namespaces := resp.nss
			if params.pack.MatchString(resp.pck.Name) {
				cpus := runtime.NumCPU()
				if len(namespaces) < cpus {
					findNamespace(results, end, namespaces, params)
					continue
				}

				// example for cpus = 4; len(namespaces) = 11
				// partitionSize = ceil(11 / 4) = 3
				// 0:4  - 0 1 2 3
				// 4:8  - 4 5 6 7
				// 8:12 - 8 9 10
				ps := int64(math.Ceil(float64(len(namespaces)) / float64(cpus)))
				var wg sync.WaitGroup
				wg.Add(cpus)
				for i := int64(0); i < int64(cpus); i++ {
					ns := namespaces[i*ps : (i+1)*ps]
					go func() {
						findNamespace(results, end, ns, params)
						wg.Done()
					}()
				}
				wg.Wait()
			}

		}
	}

	close(results)
}
