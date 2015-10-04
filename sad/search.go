package main

import (
	"../shared"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/websocket" // TODO push this back into server
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

func maybeInsensitive(pat string) string {
	if strings.ToLower(pat) == pat {
		return fmt.Sprintf("(?i)%v", pat)
	}
	return pat
}

type searchParams struct {
	pack   *regexp.Regexp
	path   *regexp.Regexp
	member *regexp.Regexp
	limit  int
}

func compileParams(pk, pt, m string, limit int) (searchParams, error) {
	var result searchParams
	var pats [3]*regexp.Regexp

	for i, p := range [3]string{pk, pt, m} {
		pat := maybeInsensitive(p)
		cp, err := regexp.Compile(pat)
		if err != nil {
			return result, err
		}
		pats[i] = cp
	}

	result.pack = pats[0]
	result.path = pats[1]
	result.member = pats[2]
	result.limit = limit

	return result, nil
}

func find(packPattern, pathPattern, memberPattern string, limit int) ([]searchResult, error) {
	var results []searchResult
	args, err := compileParams(packPattern, pathPattern, memberPattern, limit)
	if err != nil {
		return results, err
	}

	for pack, namespaces := range docs {
		if args.pack.MatchString(pack) {
			for _, namespace := range namespaces {
				normalizedPath := strings.Join(namespace.Path, ".")
				if args.path.MatchString(normalizedPath) {
					for mi, member := range namespace.Members {
						if args.member.MatchString(member.Name) {
							results = append(results, NewSearchResult(namespace, mi))
							if len(results) >= limit {
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

func streamFind(sock *websocket.Conn, packPattern, pathPattern, memberPattern string, limit int) {
	args, err := compileParams(packPattern, pathPattern, memberPattern, limit)
	if err != nil {
		log.Printf("Error while compiling params: %v\n", err)
		return
	}
	count := 0
	start := time.Now()

	for pack, namespaces := range docs {
		if args.pack.MatchString(pack) {
			for _, namespace := range namespaces {
				normalizedPath := strings.Join(namespace.Path, ".")
				if args.path.MatchString(normalizedPath) {
					for mi, member := range namespace.Members {
						if args.member.MatchString(member.Name) {
							count++
							err := sock.WriteJSON(NewSearchResult(namespace, mi))
							if err != nil {
								log.Printf("Error while writing result: %v\n", err)
								return
							}
							log.Printf("Found result #%v after %v\n", count, time.Since(start))

							if count >= limit {
								return
							}
						}
					}
				}
			}
		}
	}

	return
}
