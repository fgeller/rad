package main

import (
	"../shared"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

var serving bool
var sapServing bool

func awaitPing(addr string) error {
	limit := 10
	attempts := 0

	for {
		resp, err := http.Get("http://" + addr + "/ping")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		attempts++
		if attempts > limit {
			return fmt.Errorf("Got no ping on %v.", addr)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func ensureServe() string {
	addr := "localhost:6048"
	if !serving {
		serving = true
		go serve(addr)
	}
	return addr
}

func ensureSap() {
	if !sapServing {
		sapServing = true
		config.sapAddr = "localhost:6050"
		packsHandler := func(w http.ResponseWriter, r *http.Request) {
			data := `[{"Path":"/pack/go.zip","Name":"go","Type":"go","Version":"2015-10-08","Created":"2015-10-08T00:00:0.0+00:00"}]`
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(data))
		}
		packHandler := func(w http.ResponseWriter, r *http.Request) {
			log.Printf("test sap: Serving testdata/scala.zip\n")
			http.ServeFile(w, r, "testdata/scala.zip")
		}

		testSap := func() {
			http.HandleFunc("/pack/scala.zip", packHandler)
			http.HandleFunc("/packs", packsHandler)
			http.ListenAndServe(config.sapAddr, nil)
		}
		go testSap()
	}
}

func setup() string {
	global.packs = map[string]shared.Pack{}
	global.docs = map[string][]shared.Namespace{}
	tmp, err := ioutil.TempDir("", "sad-main-test-pack-dir")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
	}
	config.packDir = tmp
	return tmp
}

func TestServeInstalledPackInfo(t *testing.T) {

	global.docs = map[string][]shared.Namespace{
		"x": []shared.Namespace{{Members: []shared.Member{{Name: "m1"}}}},
		"y": []shared.Namespace{{Members: []shared.Member{{Name: "m2"}}}},
	}
	global.packs = map[string]shared.Pack{
		"x": shared.Pack{Name: "x", Created: time.Now()},
		"y": shared.Pack{Name: "y", Created: time.Now()},
	}

	addr := ensureServe()
	err := awaitPing(addr)
	if err != nil {
		t.Errorf("Error waiting for server to be up: %v", err)
		return
	}

	url := "http://" + addr + "/status/packs/installed"
	fmt.Printf("asking for url %v\n", url)
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("Error while querying for packs: %v", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Errorf("Error while querying for packs got status code: %v", resp.StatusCode)
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error while reading response body: %v", err)
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Errorf("Error while closing response body: %v", err)
		return
	}

	var actual map[string]shared.Pack
	err = json.Unmarshal(data, &actual)
	if err != nil {
		t.Errorf("Error while unmarshalling pack info [%s]: %v", data, err)
		return
	}

	if !reflect.DeepEqual(global.packs, actual) {
		t.Errorf(
			"Retrieved pack info was not the same. Expected:\n%v\nbut got:\n%v\n",
			global.packs,
			actual,
		)
		return
	}

}

func TestServeAvailablePacksInfo(t *testing.T) {

	global.docs = map[string][]shared.Namespace{}
	global.packs = map[string]shared.Pack{}
	addr := ensureServe()
	ensureSap()

	err := awaitPing(addr)
	if err != nil {
		t.Errorf("Error waiting for server to be up: %v", err)
		return
	}

	url := "http://" + addr + "/status/packs/available"
	fmt.Printf("asking for url %v\n", url)
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("Error while querying for packs: %v", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Errorf("Error while querying for packs got status code: %v\n%v\n", resp.StatusCode, resp)
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error while reading response body: %v", err)
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Errorf("Error while closing response body: %v", err)
		return
	}

	expected := `[{"Path":"/pack/go.zip","Name":"go","Type":"go","Version":"2015-10-08","Created":"2015-10-08T00:00:0.0+00:00"}]`

	if expected != string(data) {
		t.Errorf(
			"Retrieved available pack info was not the same. Expected:\n%v\nbut got:\n%v\n",
			expected,
			string(data),
		)
	}

}

func TestInstallAvailablePack(t *testing.T) {
	os.RemoveAll(setup())

	addr := ensureServe()
	err := awaitPing(addr)
	if err != nil {
		t.Errorf("Error while waiting for server to come up: %v", err)
		return
	}

	ensureSap()
	err = awaitPing(config.sapAddr)
	if err != nil {
		t.Errorf("Error while waiting for sap to come up: %v", err)
		return
	}

	_, err = http.Get("http://" + addr + "/install/scala.zip")
	if err != nil {
		t.Errorf("Unexpected error while trying to install pack: %v", err)
		return
	}

	if len(global.docs) == 0 || len(global.docs["scala"]) == 0 {
		t.Errorf("Expected to find installed scala docs, but got: %v", global.docs)
		return
	}
}
