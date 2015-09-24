package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"
)

var docs = map[string][]entry{}
var packDir string

type indexer func() ([]entry, error)
type downloader func(string) (*http.Response, error)
type parser func(string, io.Reader) []entry
type pack struct {
	name     string
	location string // can be URL http:// or local path
	indexer  indexer
}
type member struct {
	Name      string
	Signature string
	Target    string
}
type entry struct {
	Namespace []string
	Name      string
	Members   []member
}

func (m member) eq(other member) bool {
	return reflect.DeepEqual(m, other)
}

func (m member) String() string {
	return fmt.Sprintf(
		"member{Name: %v, Target: %v, Signature: %v}",
		m.Name,
		m.Target,
		m.Signature,
	)
}

func (e entry) eq(other entry) bool {
	return reflect.DeepEqual(e, other)
}

func (e entry) String() string {
	return fmt.Sprintf(
		"entry{Name: %v, Namespace: %v, Members: %v}",
		e.Name,
		e.Namespace,
		e.Members,
	)
}

func isSameEntry(a entry, b entry) bool {
	if a.Name != b.Name ||
		len(a.Namespace) != len(b.Namespace) {
		return false
	}
	for i := range a.Namespace {
		if a.Namespace[i] != b.Namespace[i] {
			return false
		}
	}

	return true
}

func mergeEntries(entries []entry) []entry {
	if len(entries) < 1 {
		return entries
	}

	unmerged := entries[1:]
	merged := []entry{entries[0]}

merging:
	for ui := range unmerged {
		for mi := range merged {
			if isSameEntry(unmerged[ui], merged[mi]) {
				merged[mi].Members = append(merged[mi].Members, unmerged[ui].Members...)
				continue merging
			}
		}
		merged = append(merged, unmerged[ui])
	}

	return merged
}

func unmarshalPack(pack pack, dataPath string) error {
	start := time.Now()

	data, err := ioutil.ReadFile(dataPath)
	if err != nil {
		return err
	}

	var es []entry
	err = json.Unmarshal(data, &es)
	if err != nil {
		return err
	}

	docs[pack.name] = es
	log.Printf(
		"Deserialized [%v] entries for [%v] in %v.",
		len(es),
		pack.name,
		time.Since(start),
	)

	return nil
}

func main() {
	flag.StringVar(&packDir, "packDir", "packs", "Path where packages will be installed")
	flag.Parse()

	install(
		pack{
			name:     "java",
			location: "/Users/fgeller/Downloads/jdk-8u60-docs-all.zip",
			indexer:  indexJavaApi("java"),
		},
	)

	install(
		pack{
			name:     "scala",
			location: "http://downloads.typesafe.com/scala/2.11.7/scala-docs-2.11.7.zip",
			indexer:  indexScalaApi("scala"),
		},
	)

	install(
		pack{
			name:     "akka",
			location: "/Users/fgeller/Downloads/akka-2.3.14/doc/akka-2.3.14.zip",
			indexer:  indexScalaApi("akka"),
		},
	)

	serve("0.0.0.0:3024")
}
