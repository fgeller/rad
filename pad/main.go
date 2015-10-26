package main

import (
	"../shared"

	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type indexer func(string) ([]shared.Namespace, error)
type parser func(string, io.Reader) []shared.Namespace

type config struct {
	indexer     indexer
	name        string
	Type        string
	version     string
	source      string
	description string
	destination string
}

func mkIndexer(name string, source string) indexer {
	mk := func(fn func(string, io.Reader) []shared.Namespace) indexer {
		return func(path string) ([]shared.Namespace, error) {
			return scan(path, fn)
		}
	}

	switch name {
	case "java":
		return mk(parseJavaDocFile)
	case "scala":
		return mk(parseScalaDocFile)
	case "go":
		return mk(parseGoDocFile)
	case "clojure":
		return mk(parseClojureDocFile)
	case "py27":
		return mk(parsePy27DocFile)
	case "django":
		return mk(parseDjangoDocFile)
	case "react":
		return mk(parseReactDocFile)
	case "jquery":
		return mk(parseJQueryDocFile)
	case "lodash":
		return mk(parseLodashDocFile)
	default:
		log.Fatalf("Unsupported indexer name: %v\n", name)
	}

	return nil
}

func mkConfig(
	indexerName string,
	packName string,
	source string,
	version string,
	description string,
	dest string,
) (config, error) {
	var conf config
	source, err := filepath.Abs(source)
	if err != nil {
		log.Fatalf("Error finding absolute path: %v\n", err)
	}

	if dest != "" && shared.FileExists(dest) {
		fi, err := os.Stat(dest)
		if err != nil {
			log.Fatalf("Failed to stat destination %v: %v", dest, err)
		}
		if !fi.IsDir() {
			log.Fatalf("Destination should be a directory: %v", dest)
		}
	}

	if dest != "" && !shared.FileExists(dest) {
		os.MkdirAll(dest, 0755)
		log.Printf("Created destination %v\n", dest)
	}

	conf = config{
		indexer:     mkIndexer(indexerName, source),
		Type:        indexerName,
		name:        packName,
		source:      source,
		version:     version,
		description: description,
		destination: dest,
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

	log.Printf("Copying files to %v\n", targetDir)
	// expected: name: scala, source: /some/path/to/scala-docs
	// created:
	//   /tmp-dir/scala/scala-docs/
	//   /tmp-dir/scala/pack.json
	//   /tmp-dir/scala/data.json
	c, err := shared.CopyDir(conf.source, targetDir)
	if err != nil {
		return "", err
	}
	log.Printf("Copied %v files over.\n", c)

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
	log.Printf("Made targets relative to pack folder.\n")

	pack := shared.Pack{
		Name:        conf.name,
		Type:        conf.Type,
		Version:     conf.version,
		Created:     time.Now(),
		Description: conf.description,
	}

	// 2. serialize conf
	jsonPack, err := json.MarshalIndent(pack, "", "  ")
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
	jsonEntries, err := json.MarshalIndent(entries, "", "  ")
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
	dt := pack.Created.Format(time.RFC3339)[:len("9999-99-99")]
	fn := fmt.Sprintf("%v-%v-%v.zip", pack.Name, pack.Version, dt)
	outFile := filepath.Join(tmpDir, fn)
	out, err := os.Create(outFile)
	if err != nil {
		return "", err
	}
	defer out.Close()

	err = shared.ZipDir(out, filepath.Join(tmpDir, conf.name))
	if err != nil {
		return "", err
	}

	err = out.Close()
	if err != nil {
		return "", err
	}

	log.Printf("Zipped files into %v.\n", out.Name())

	if conf.destination == "" {
		return out.Name(), nil
	}

	dest := filepath.Join(conf.destination, filepath.Base(out.Name()))
	err = os.Link(out.Name(), dest)
	if err != nil {
		return "", err
	}

	log.Printf("Linked %v to %v.", out.Name(), dest)
	return dest, nil
}

func main() {
	var (
		indexerName = flag.String("indexer", "", "Indexer type for this pack (scala, java, go, clojure, py27)")
		packName    = flag.String("name", "", "Name for this pack")
		source      = flag.String("source", "", "Source directory for this pack")
		version     = flag.String("version", "", "Version string for this pack")
		description = flag.String("desc", "", "Description string for this pack")
		destination = flag.String("dest", "", "Destination folder (pack remains in tmp location when omitted).")
	)

	flag.Parse()
	conf, err := mkConfig(
		*indexerName,
		*packName,
		*source,
		*version,
		*description,
		*destination,
	)
	if err != nil {
		log.Fatalf("Invalid configuration %v", conf)
	}

	result, err := mkPack(conf)
	if err != nil {
		log.Fatalf("Failed to create pack: %v\n", err)
		return
	}

	log.Printf("Created pack: %v\n", result)
}
