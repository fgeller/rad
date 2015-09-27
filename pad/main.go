package main

import (
	"../shared"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type indexer func(string) ([]shared.Entry, error)
type parser func(string, io.Reader) []shared.Entry
type downloader func(string) (*http.Response, error)

type config struct {
	indexer indexer
	name    string
	Type    string
	source  string
}

func isValidConfig(conf config) bool {
	return len(conf.name) > 0 &&
		len(conf.Type) > 0 &&
		fileExists(conf.source) &&
		conf.indexer != nil
}

func mkIndexer(name string, source string) indexer {
	mk := func(fn func(string, io.Reader) []shared.Entry) indexer {
		return func(path string) ([]shared.Entry, error) {
			return scan(path, fn)
		}
	}

	fmt.Printf("SOURCE: %v\n", source)

	switch name {
	case "java":
		return mk(parseJavaDocFile)
		// case "scala":
		// 	return mk(parseScalaDocFile)
	}

	return nil
}

func mkConfig(indexerName string, packName string, source string) (config, error) {
	source, err := filepath.Abs(source)
	if err != nil {
		return config{}, err
	}

	conf := config{
		indexer: mkIndexer(indexerName, source),
		Type:    indexerName,
		name:    packName,
		source:  source,
	}

	if !isValidConfig(conf) {
		return conf, fmt.Errorf("Invalid configuration.")
	}

	return conf, nil
}

func mkPack(conf config) (string, error) {
	packPrefix := fmt.Sprintf("pad-create-%v", conf.name)
	tmpDir, err := ioutil.TempDir("", packPrefix)
	if err != nil {
		log.Printf("Error while creating temp dir: %v\n", err)
		return "", err
	}

	targetDir := filepath.Join(tmpDir, conf.name)
	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		log.Printf("Failed to create directories for pack.\n")
		return "", err
	}

	// expected: name: scala, source: /some/path/to/scala-docs
	// created:
	//   /tmp-dir/scala/scala-docs/
	//   /tmp-dir/scala/pack.json
	//   /tmp-dir/scala/data.json
	err = copy(conf.source, targetDir)
	if err != nil {
		return "", err
	}
	log.Printf("Copied files over to %v.\n", targetDir)

	// 1. index
	entries, err := conf.indexer(targetDir)

	// TODO: push this into parser?
	// convert to relative Targets
	for e := range entries {
		for m := range entries[e].Members {
			target := entries[e].Members[m].Target
			rel, err := filepath.Rel(tmpDir, target)
			if err != nil {
				return "", err
			}

			entries[e].Members[m].Target = rel
		}
	}
	log.Printf("Make targets relative.\n")

	pack := shared.Pack{
		Name: conf.name,
		Type: conf.Type,
	}

	// 2. serialize conf
	jsonPack, err := json.Marshal(pack)
	if err != nil {
		return "", err
	}

	packFilePath := filepath.Join(tmpDir, conf.name, "pack.json")
	packFile, err := os.Create(packFilePath)
	if err != nil {
		log.Printf("Cannot create pack file.\n")
		return "", err
	}

	_, err = packFile.WriteString(string(jsonPack))
	if err != nil {
		log.Printf("Cannot write to pack file.\n")
		return "", err
	}

	err = packFile.Close()
	if err != nil {
		log.Printf("Cannot close pack file.\n")
		return "", err
	}

	// 3. serialize entries
	jsonEntries, err := json.Marshal(entries)
	if err != nil {
		return "", err
	}

	dataFilePath := filepath.Join(tmpDir, conf.name, "data.json")
	dataFile, err := os.Create(dataFilePath)
	if err != nil {
		return "", err
	}

	_, err = dataFile.WriteString(string(jsonEntries))
	if err != nil {
		return "", err
	}

	err = dataFile.Close()
	if err != nil {
		return "", err
	}

	log.Printf("Serialized pack and data files.\n")

	// 4. zip it all up
	outFile := filepath.Join(tmpDir, conf.name+".zip")
	out, err := os.Create(outFile)
	if err != nil {
		return "", err
	}
	defer out.Close()

	err = zipDir(out, filepath.Join(tmpDir, conf.name))
	if err != nil {
		return "", err
	}

	err = out.Close()
	if err != nil {
		return "", err
	}

	log.Printf("Zipped files into %v.\n", out.Name())

	return out.Name(), nil
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

	result, err := mkPack(conf)
	fmt.Printf("RESULT: %v\nERR: %v\n", result, err)

}
