package main

import (
	"../shared"
	"encoding/json"
	"io/ioutil"
	//	"os"
	"fmt"
	"testing"
)

func TestMkPack(t *testing.T) {
	indexerName := "java"
	packName := "jdk"
	packSource := mkPath("testdata", "jdk")

	conf, err := mkConfig(indexerName, packName, packSource)
	if err != nil {
		t.Errorf("Unexpected error while creating config: %v", err)
		return
	}

	actual, err := mkPack(conf)
	if err != nil || !shared.FileExists(actual) {
		t.Errorf("Expected pack archive to be created, err: %v.", err)
		return
	}
	fmt.Printf("created sample: %v\n", actual)
	// defer os.RemoveAll(actual) // TODO: woa?

	tmpDir, err := ioutil.TempDir("", "pad-test-pack")
	if err != nil {
		t.Errorf("Unexpected error while creating temporary directory: %v", err)
		return
	}
	// defer os.RemoveAll(tmpDir)

	err = shared.Unzip(actual, tmpDir)
	if err != nil {
		t.Errorf("Error while unzipping archive %v: %v", actual, err)
		return
	}

	packConfigFile := mkPath(tmpDir, packName, "pack.json")
	if !shared.FileExists(packConfigFile) {
		t.Errorf("Expected config file %v.", packConfigFile)
	}

	packConfigStr, err := ioutil.ReadFile(packConfigFile)
	if err != nil {
		t.Errorf("Couldn't read pack config file: %v", err)
	}

	var pack shared.Pack
	err = json.Unmarshal(packConfigStr, &pack)
	if err != nil {
		t.Errorf("Couldn't unmarshall pack config file: %v", err)
	}

	packDataFile := mkPath(tmpDir, packName, "data.json")
	if !shared.FileExists(packDataFile) {
		t.Errorf("Expected data file %v.", packDataFile)
	}

	packDataStr, err := ioutil.ReadFile(packDataFile)
	if err != nil {
		t.Errorf("Couldn't read pack data file: %v", err)
	}

	var namespaces []shared.Namespace
	err = json.Unmarshal(packDataStr, &namespaces)
	if err != nil {
		t.Errorf("Couldn't unmarshall data file: %v", err)
	}

}

func testIndexer(path string) ([]shared.Namespace, error) {
	return []shared.Namespace{}, nil
}

func TestIsValidConfig(t *testing.T) {
	validName := "b"
	validType := "java"
	validSource := "/tmp"
	validIndexer := testIndexer

	invalidName := ""
	invalidType := ""
	invalidSource := "/xxxxxxxx"
	invalidIndexer := indexer(nil)

	conf := config{}
	actual := isValidConfig(conf)
	if actual {
		t.Errorf("Expected empty config to be invalid.\n")
	}

	conf = config{
		indexer: validIndexer,
		name:    validName,
		Type:    validType,
		source:  validSource,
	}
	actual = isValidConfig(conf)
	if !actual {
		t.Errorf("Expected proper config to be valid.\n")
	}

	conf = config{
		indexer: validIndexer,
		name:    invalidName,
		Type:    validType,
		source:  validSource,
	}
	actual = isValidConfig(conf)
	if actual {
		t.Errorf("Expected invalid name to be invalid.\n")
	}

	conf = config{validIndexer, validName, validType, invalidSource}
	actual = isValidConfig(conf)
	if actual {
		t.Errorf("Expected invalid source to be invalid.\n")
	}

	conf = config{
		indexer: invalidIndexer,
		name:    validName,
		Type:    validType,
		source:  validSource,
	}
	actual = isValidConfig(conf)
	if actual {
		t.Errorf("Expected invalid indexer to be invalid.\n")
	}

	conf = config{
		indexer: validIndexer,
		name:    validName,
		Type:    invalidType,
		source:  validSource,
	}
	actual = isValidConfig(conf)
	if actual {
		t.Errorf("Expected invalid type to be invalid.\n")
	}

}
