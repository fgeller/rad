package main

import (
	"../shared"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

var serving bool

func setup() string {
	tmp, err := ioutil.TempDir("", "sap-pack-dir")
	if err != nil {
		log.Fatalf("Could not create temporary directory: %v\n", err)
	}

	for _, fn := range []string{"go.zip", "java.zip"} {
		err = os.Link("testdata/"+fn, filepath.Join(tmp, fn))
		if err != nil {
			log.Fatalf("Could not link %v to temp directory: %v\n", fn, err)
		}
	}

	config.PackDir = tmp
	return tmp
}

func ensureServe(addr string) {
	if !serving {
		serving = true
		go serve(addr)
	}
}

func TestServePack(t *testing.T) {
	defer os.RemoveAll(setup())

	addr := "localhost:6050"
	ensureServe(addr)

	err := awaitPing(addr)
	if err != nil {
		t.Errorf("Unexpected error while waiting for ping: %v\n", err)
		return
	}

	res, err := http.Get("http://" + addr + "/pack/go.zip")
	if err != nil {
		t.Errorf("Unexpected error when asking for go.zip: %v\n", err)
		return
	}

	if res.Header["Content-Type"][0] != "application/zip" {
		t.Errorf("Expected application/zip but got: %v\n", res.Header["Content-Type"])
		return
	}

	if res.Header["Content-Length"][0] != "386" {
		t.Errorf("Expected application/zip but got: %v\n", res.Header["Content-Type"])
		return
	}

	res, err = http.Get("http://" + addr + "/pack/java.zip")
	if err != nil {
		t.Errorf("Unexpected error when asking for java.zip: %v\n", err)
		return
	}

	if res.Header["Content-Type"][0] != "application/zip" {
		t.Errorf("Expected application/zip but got: %v\n", res.Header["Content-Type"])
		return
	}

	if res.Header["Content-Length"][0] != "391" {
		t.Errorf("Expected application/zip but got: %v\n", res.Header["Content-Type"])
		return
	}
}

func TestServingPacks(t *testing.T) {
	defer os.RemoveAll(setup())

	addr := "localhost:6050"
	ensureServe(addr)

	err := awaitPing(addr)
	if err != nil {
		t.Errorf("Unexpected error while waiting for ping: %v\n", err)
		return
	}

	res, err := http.Get("http://" + addr + "/packs")
	if err != nil {
		t.Errorf("Unexpected error while asking for packs: %v\n", err)
		return
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Errorf("Unexpected error while reading response body: %v\n", err)
		return
	}

	var actual []packListing
	err = json.Unmarshal(data, &actual)
	if err != nil {
		t.Errorf("Unexpected error while unmarshalling response body: %v\n", err)
		return
	}

	creationTime := time.Date(2015, time.October, 4, 0, 0, 0, 0, time.UTC)
	expected := []packListing{
		{
			Pack: &shared.Pack{
				Name:    "go",
				Type:    "go",
				Version: "2015-10-04",
				Created: creationTime,
			},
			Path: "/pack/go.zip",
		},
		{
			Pack: &shared.Pack{
				Name:    "java",
				Type:    "java",
				Version: "jdk8",
				Created: creationTime,
			},
			Path: "/pack/java.zip",
		},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected packs:\n%v\nBut got:\n%v\n", expected, actual)
		return
	}
}

func awaitPing(addr string) error {
	for i := 0; i < 5; i++ {
		_, err := http.Get("http://" + addr + "/ping")
		if err == nil {
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}

	return fmt.Errorf(
		"Got no successful ping in %v",
		5*10*time.Millisecond,
	)
}
