package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var docs = map[string][]entry{}

type indexer func() ([]entry, error)
type downloader func(string) (*http.Response, error)
type parser func(string, io.Reader) []entry
type pack struct {
	name     string
	location string // can be URL http:// or local path
	indexer  indexer
}
type entry struct {
	Namespace []string
	Entity    string
	Member    string
	Signature string
	Target    string // location relative to `packDir` where to find documentation
	Source    string
}

func (e entry) String() string {
	return fmt.Sprintf(
		"entry{Namespace: %v, Entity: %v, Member: %v, Signature: %v, Target: %v, Source: %v}",
		e.Namespace,
		e.Entity,
		e.Member,
		e.Signature,
		e.Target,
		e.Source,
	)
}

func (e entry) eq(other entry) bool { // TODO: reflect.DeepEqual?
	if len(e.Namespace) != len(other.Namespace) {
		return false
	}

	for i, n := range e.Namespace {
		if other.Namespace[i] != n {
			return false
		}
	}

	return (e.Entity == other.Entity &&
		e.Member == other.Member &&
		e.Signature == other.Signature &&
		e.Target == other.Target &&
		e.Source == other.Source)
}

func unmarshalPack(pack pack, dataPath string) error {
	start := time.Now()

	data, err := ioutil.ReadFile(dataPath)
	if err != nil {
		return err // TODO: or re-download?
	}

	var es []entry
	err = json.Unmarshal(data, &es)
	if err != nil {
		return err // TODO: or re-download?
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

	serve("0.0.0.0:3024")
}
