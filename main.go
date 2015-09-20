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
type member struct {
	Name      string
	Signature string
	Target    string
	Source    string
}
type entry struct {
	Namespace []string
	Name      string
	Members   []member
	Source    string
}

func (m member) eq(other member) bool {
	return m.Name == other.Name &&
		m.Signature == other.Signature &&
		m.Target == other.Target &&
		m.Source == other.Source
}

func (m member) String() string {
	return fmt.Sprintf(
		"member{Name: %v, Target: %v, Signature: %v, Source: %v}",
		m.Name,
		m.Target,
		m.Signature,
		m.Source,
	)
}

func (e entry) eq(other entry) bool { // TODO: reflect.DeepEqual?
	if len(e.Namespace) != len(other.Namespace) || len(e.Members) != len(other.Members) {
		return false
	}

	for i, n := range e.Namespace {
		if other.Namespace[i] != n {
			return false
		}
	}

	for i, m := range e.Members {
		if !other.Members[i].eq(m) {
			return false
		}
	}

	return (e.Name == other.Name && e.Source == other.Source)
}

func (e entry) String() string {
	return fmt.Sprintf(
		"entry{Name: %v, Namespace: %v, Members: %v, Source: %v}",
		e.Name,
		e.Namespace,
		e.Members,
		e.Source,
	)
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

	install(
		pack{
			name:     "akka",
			location: "/Users/fgeller/Downloads/akka-2.3.14/doc/akka-2.3.14.zip",
			indexer:  indexScalaApi("akka"),
		},
	)

	serve("0.0.0.0:3024")
}
