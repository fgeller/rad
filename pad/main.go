package main

import (
	"../shared"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

type indexer func(string) ([]shared.Entry, error)
type parser func(string, io.Reader) []shared.Entry
type downloader func(string) (*http.Response, error)

type config struct {
	indexer indexer
	name    string
	source  string
}

func isValidConfig(conf config) bool {
	return len(conf.name) > 0 &&
		fileExists(conf.source) &&
		conf.indexer != nil
}

func mkIndexer(name string, source string) indexer {
	mk := func(fn func(string, io.Reader) []shared.Entry) indexer {
		return func(path string) ([]shared.Entry, error) {
			return scan(path, fn)
		}
	}

	switch name {
	case "java":
		return mk(parseJavaDocFile)
	case "scala":
		return mk(parseScalaDocFile)
	}

	return nil
}

func mkConfig(indexerName string, packName string, source string) (config, error) {
	conf := config{
		indexer: mkIndexer(indexerName, source),
		name:    packName,
		source:  source,
	}

	if !isValidConfig(conf) {
		return conf, fmt.Errorf("Invalid configuration.")
	}

	return conf, nil
}

func mkPack(conf config) {
	// 0. copy all into temp location
	// 1. index
	// 2. serialize conf
	// 3. serialize entries
	// 4. zip it all up

}

func main() {
	var (
		indexerName = flag.String("indexer", "", "Indexer type for this pack (scala, java)")
		packName    = flag.String("name", "", "Name for this pack")
		source      = flag.String("source", "", "Source directory for this pack")
	)

	flag.Parse()
	conf, err := mkConfig(
		*indexerName,
		*packName,
		*source,
	)
	if err != nil {
		log.Fatal("Invalid configuration %v", conf)
	}

	mkPack(conf)

}
