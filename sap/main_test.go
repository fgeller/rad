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

	packDir = tmp
	return tmp
}

func TestServingPacks(t *testing.T) {
	defer os.RemoveAll(setup())

	addr := "localhost:6048"
	go serve(addr)

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
	log.Printf("Read response %s\n", data)

	var actual []shared.Pack
	err = json.Unmarshal(data, &actual)
	if err != nil {
		t.Errorf("Unexpected error while unmarshalling response body: %v\n", err)
		return
	}

	creationTime := time.Date(2015, time.October, 4, 0, 0, 0, 0, time.UTC)
	expected := []shared.Pack{
		{Name: "go", Type: "go", Version: "2015-10-04", Created: creationTime},
		{Name: "java", Type: "java", Version: "jdk8", Created: creationTime},
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
