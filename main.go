package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var docs = map[string][]entry{}
var packDir = "packs"

type indexer func() ([]entry, error)
type downloader func(string) (*http.Response, error)
type parser func(string, io.Reader) []entry
type pack struct {
	name    string
	url     string
	indexer indexer
}
type entry struct {
	Namespace []string
	Entity    string
	Function  string
	Signature string
	Target    string // location relative to `packDir` where to find documentation
	source    string
}

func (e entry) String() string {
	return fmt.Sprintf("entry{Namespace: %v, Entity: %v, Function: %v, Signature: %v}", e.Namespace, e.Entity, e.Function, e.Signature)
}

func (e entry) eq(other entry) bool {
	if len(e.Namespace) != len(other.Namespace) {
		return false
	}

	for i, n := range e.Namespace {
		if other.Namespace[i] != n {
			return false
		}
	}

	return e.Entity == other.Entity &&
		e.Function == other.Function &&
		e.Signature == other.Signature // TODO: expand
}

func install(pack pack) error {
	dataPath := packDir + string(os.PathSeparator) +
		pack.name + string(os.PathSeparator) +
		"rad-data.json"

	if fileExists(dataPath) {
		log.Printf("Already installed pack [%v], deserializing entries.", pack.name)
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

	local, err := download(http.Get, pack.url)
	if err != nil {
		log.Fatalf("Failed to download [%v] err: %v.\n", pack.url, err)
		return err
	}
	defer os.Remove(local)

	err = unzip(local, packDir+string(os.PathSeparator)+pack.name)
	if err != nil {
		log.Fatalf("Failed to unzip archive [%v], err: %v", local, err)
		return err
	}

	docs[pack.name], err = pack.indexer()
	if err != nil {
		return err
	}

	datPath := packDir + string(os.PathSeparator) +
		pack.name + string(os.PathSeparator) +
		"rad-data.json"

	data, err := json.Marshal(docs[pack.name])
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(datPath, data, 0644)
	if err != nil {
		return err
	}

	log.Printf("Installed [%v] entries for pack [%v].", len(docs[pack.name]), pack.name)
	return nil
}

func main() {
	install(
		pack{
			name:    "scala",
			url:     "http://downloads.typesafe.com/scala/2.11.7/scala-docs-2.11.7.zip",
			indexer: indexScalaApi("scala"),
		},
	)

	serve("0.0.0.0:3024")
}
