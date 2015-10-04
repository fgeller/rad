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

func TestServePackInfo(t *testing.T) {

	docs = map[string][]shared.Namespace{
		"x": []shared.Namespace{{Members: []shared.Member{{Name: "m1"}}}},
		"y": []shared.Namespace{{Members: []shared.Member{{Name: "m2"}}}},
	}
	packs = map[string]shared.Pack{
		"x": shared.Pack{Name: "x", Created: time.Now()},
		"y": shared.Pack{Name: "y", Created: time.Now()},
	}

	addr := "localhost:6048"

	go serve(addr)
	err := awaitPing(addr)
	if err != nil {
		t.Errorf("Error waiting for server to be up: %v", err)
		return
	}

	time.Sleep(500 * time.Millisecond)
	url := "http://" + addr + "/status/packs"
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
