package main

import (
	"../shared"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMkPack(t *testing.T) {
	indexerName := "java"
	packName := "jdk"
	packSource := filepath.Join("testdata", "jdk")
	packVersion := "1.2.3"

	conf, err := mkConfig(indexerName, packName, packSource, packVersion)
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
	defer os.RemoveAll(actual)

	tmpDir, err := ioutil.TempDir("", "pad-test-pack")
	if err != nil {
		t.Errorf("Unexpected error while creating temporary directory: %v", err)
		return
	}
	defer os.RemoveAll(tmpDir)

	err = shared.Unzip(actual, tmpDir)
	if err != nil {
		t.Errorf("Error while unzipping archive %v: %v", actual, err)
		return
	}

	packConfigFile := filepath.Join(tmpDir, packName, "pack.json")
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

	if pack.Version != packVersion ||
		time.Now().Before(pack.Created) {
		t.Errorf("Unexpected pack parameters: %v\n", pack)
	}

	packDataFile := filepath.Join(tmpDir, packName, "data.json")
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
