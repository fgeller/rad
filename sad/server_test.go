package main

import (
	"../shared"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

var serving bool
var sapServing bool

func ensureServe(addr string) {
	if !serving {
		serving = true
		go serve(addr)
	}
}

func ensureSap(addr string) {
	if !sapServing {
		sapServing = true
		packHandler := func(w http.ResponseWriter, r *http.Request) {
			data := `[{"Path":"/pack/go.zip","Name":"go","Type":"go","Version":"2015-10-08","Created":"2015-10-08T00:00:0.0+00:00"}]`
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(data))
		}
		testSap := func() {
			http.HandleFunc("/packs", packHandler)
			http.ListenAndServe(addr, nil)
		}
		go testSap()
	}
}

func TestServeInstalledPackInfo(t *testing.T) {

	docs = map[string][]shared.Namespace{
		"x": []shared.Namespace{{Members: []shared.Member{{Name: "m1"}}}},
		"y": []shared.Namespace{{Members: []shared.Member{{Name: "m2"}}}},
	}
	packs = map[string]shared.Pack{
		"x": shared.Pack{Name: "x", Created: time.Now()},
		"y": shared.Pack{Name: "y", Created: time.Now()},
	}

	addr := "localhost:6048"

	ensureServe(addr)
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

	if !reflect.DeepEqual(packs, actual) {
		t.Errorf(
			"Retrieved pack info was not the same. Expected:\n%v\nbut got:\n%v\n",
			packs,
			actual,
		)
		return
	}

}

func TestServeAvailablePacksInfo(t *testing.T) {

	docs = map[string][]shared.Namespace{}
	packs = map[string]shared.Pack{}
	addr := "localhost:6048"
	sapAddr = "localhost:6050"

	ensureServe(addr)
	ensureSap(sapAddr)

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

	expected := `[{"Path":"/pack/go.zip","Name":"go","Type":"go","Version":"2015-10-08","Created":"2015-10-08T00:00:0.0+00:00"}]`

	if expected != string(data) {
		t.Errorf(
			"Retrieved available pack info was not the same. Expected:\n%v\nbut got:\n%v\n",
			expected,
			string(data),
		)
	}

}

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
