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
var sapServing bool

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
			data := `[{"File":"go.zip","Name":"go","Type":"go","Version":"2015-10-08","Created":"2015-10-08T00:00:00.0+00:00"}]`
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
	resetGlobals()
	tmp, err := ioutil.TempDir("", "sad-main-test-pack-dir")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
	}
	config.packDir = tmp
	return tmp
}

// a bit annoying, but this has to go early as it relies on the location of the
// packDir being passed to the http.FileServer -- which happens only once when
// these tests are run
func TestServeEscapedPackFile(t *testing.T) {
	pd := setup()
	defer os.RemoveAll(pd)

	fn := "%3C!!.html"
	_, err := os.Create(filepath.Join(pd, fn))
	if err != nil {
		t.Errorf("Error creating escaped file: %v", err)
		return
	}

	addr := ensureServe()
	err = awaitPing(addr)
	if err != nil {
		t.Errorf("Error waiting for server to be up: %v", err)
		return
	}

	r, err := http.Get("http://" + addr + "/pack/" + fn)
	if err != nil {
		t.Errorf("Error getting escaped pack file: %v\n", err)
		return
	}

	if r.StatusCode != 200 {
		t.Errorf("Expected 200 got: %v\n", r)
		return
	}
}

func TestServeInstalledPackInfo(t *testing.T) {
	resetGlobals()

	loadPack(
		shared.Pack{Name: "x", Created: time.Now()},
		[]shared.Namespace{{Members: []shared.Member{{Name: "m1"}}}},
	)

	loadPack(
		shared.Pack{Name: "y", Created: time.Now()},
		[]shared.Namespace{{Members: []shared.Member{{Name: "m2"}}}},
	)

	addr := ensureServe()
	err := awaitPing(addr)
	if err != nil {
		t.Errorf("Error waiting for server to be up: %v", err)
		return
	}

	url := "http://" + addr + "/status"
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
		t.Errorf("Error while reading response body: %s", err)
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Errorf("Error while closing response body: %v", err)
		return
	}

	var actual StatusInfo
	err = json.Unmarshal(data, &actual)
	if err != nil {
		t.Errorf("Error while unmarshalling pack info [%s]: %v", data, err)
		return
	}

	installed := installedPacks()

comparing:
	for _, p := range installed {
		for _, ip := range actual.Packs.Installed {
			if reflect.DeepEqual(ip, p) {
				continue comparing
			}
		}
		t.Errorf("Expected %v to be installed, got %v.", p, actual.Packs.Installed)
		return
	}
}

func TestServeHTTPSearch(t *testing.T) {
	resetGlobals()

	loadPack(
		shared.Pack{Name: "x", Created: time.Now()},
		[]shared.Namespace{{Members: []shared.Member{{Name: "m1"}, {Name: "m2"}}}},
	)

	addr := ensureServe()
	err := awaitPing(addr)
	if err != nil {
		t.Errorf("Error waiting for server to be up: %v", err)
		return
	}

	url := "http://" + addr + "/q?pack=&path=&member=m&limit=1"
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("Error querying: %v", err)
		return
	}
	if resp.StatusCode != 200 {
		t.Errorf("Error querying got status code: %v", resp.StatusCode)
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %s", err)
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Errorf("Error closing response body: %v", err)
		return
	}

	var actual []searchResult
	err = json.Unmarshal(data, &actual)
	if err != nil {
		t.Errorf("Error while unmarshalling search results [%s]: %v", data, err)
		return
	}

	if len(actual) != 1 {
		t.Errorf("Expected one search result but got %v.\n", actual)
		return
	}
	var expected = []searchResult{
		{
			Namespace: "",
			Member:    "m1",
			Target:    "/pack",
		},
	}
	if reflect.DeepEqual(expected, actual) {
		t.Errorf("Unexpected search result, expected:\n%v\nbut got:\n%v\n", expected, actual)
		return
	}
}

func TestServeAvailablePacksInfo(t *testing.T) {
	resetGlobals()
	addr := ensureServe()
	ensureSap()

	err := awaitPing(addr)
	if err != nil {
		t.Errorf("Error waiting for server to be up: %v", err)
		return
	}

	url := "http://" + addr + "/status"
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
		t.Errorf("Error reading response body: %v", err)
		return
	}
	err = resp.Body.Close()
	if err != nil {
		t.Errorf("Error closing response body: %v", err)
		return
	}

	var actual StatusInfo
	err = json.Unmarshal(data, &actual)
	if err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
		return
	}

	// expected := `[{"Path":"/pack/go.zip","Name":"go","Type":"go","Version":"2015-10-08","Created":"2015-10-08T00:00:0.0+00:00"}]`
	installed := []shared.Pack{}
	if err != nil {
		t.Errorf("Failed to parse date: %v", err)
		return
	}
	// "2015-10-08T00:00:0.0+00:00"
	available := []shared.Pack{
		{
			File:    "go.zip",
			Name:    "go",
			Type:    "go",
			Version: "2015-10-08",
			Created: time.Date(2015, time.October, 8, 0, 0, 0, 0, time.UTC),
		},
	}

	expected := StatusInfo{
		Packs: PackInfo{
			Installed: installed,
			Available: available,
		},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"Retrieved available pack info was not the same. Expected:\n%+v\nbut got:\n%+v\n",
			expected,
			actual,
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

	docs := installedDocs()
	if len(docs) == 0 || len(docs["scala"]) == 0 {
		t.Errorf("Expected to find installed scala docs, but got: %v", docs)
		return
	}
}

func TestRemoveInstalledPack(t *testing.T) {
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

	if _, err = http.Get("http://" + addr + "/install/scala.zip"); err != nil {
		t.Errorf("Unexpected error when installing scala.zip: %v", err)
		return
	}

	docs := installedDocs()
	if len(docs) == 0 {
		t.Errorf("Expected to find scala docs installed, but got: %v", docs)
		return
	}

	_, err = http.Get("http://" + addr + "/remove/scala")
	if err != nil {
		t.Errorf("Unexpected error while trying to remove pack: %v", err)
		return
	}

	ps, err := ioutil.ReadDir(config.packDir)
	if err != nil {
		t.Errorf("Unexpected error while trying to read contents of pack dir: %v", err)
		return
	}

	if len(ps) != 0 {
		t.Errorf("Expected pd to be empty, but got: %s", ps)
		return
	}
}

func TestDontRemoveInstalledPackWhenReadOnly(t *testing.T) {
	os.RemoveAll(setup())

	config.readOnly = true
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

	p := shared.Pack{Name: "scala"}
	nss := []shared.Namespace{
		shared.Namespace{
			Path:    "a",
			Members: []shared.Member{{Name: "b"}},
		},
	}
	loadPack(p, nss)

	docs := installedDocs()
	if len(docs) != 1 {
		t.Errorf("Expected to find one pack installed, but got: %v", docs)
		return
	}

	resp, err := http.Get("http://" + addr + "/remove/scala")
	if err != nil {
		t.Errorf("Unexpected error while trying to remove pack: %v", err)
		return
	}
	if resp.StatusCode != 403 {
		t.Errorf("Expected forbidden status code (403), got: %v", resp.StatusCode)
		return
	}

	if len(docs) != 1 {
		t.Errorf("Expected to find one pack installed, but got: %v", docs)
		return
	}
}

func TestDontInstallPackWhenReadOnly(t *testing.T) {
	os.RemoveAll(setup())

	config.readOnly = true
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

	docs := installedDocs()
	if len(docs) != 0 {
		t.Errorf("Expected to find no pack installed, but got: %v", docs)
		return
	}

	resp, err := http.Get("http://" + addr + "/install/scala.zip")
	if err != nil {
		t.Errorf("Unexpected error while trying to remove pack: %v", err)
		return
	}
	if resp.StatusCode != 403 {
		t.Errorf("Expected forbidden status code (403), got: %v", resp.StatusCode)
		return
	}

	if len(docs) != 0 {
		t.Errorf("Expected to find no pack installed, but got: %v", docs)
		return
	}
}

func TestServeAsset(t *testing.T) {
	os.RemoveAll(setup())

	dir := "testdata/assets"
	assets, err := shared.LoadAssets(dir)
	if err != nil {
		t.Errorf("Error while loading assets from %v: %v", dir, err)
		return
	}
	global.assets = assets

	addr := ensureServe()
	err = awaitPing(addr)
	if err != nil {
		t.Errorf("Error while waiting for server to come up: %v", err)
		return
	}

	walker := func(p string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			return err
		}

		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}

		res, err := http.Get("http://" + addr + "/a/" + rel)
		if err != nil {
			t.Errorf("Unexpected error while trying to requesting asset %v: %v", rel, err)
			return err
		}

		if res.StatusCode != 200 {
			t.Errorf("Expected 200 when accessing asset %v, got %+v", rel, res)
			return fmt.Errorf("Expected 200 when accessing asset %v, got %+v", rel, res)
		}

		return nil
	}

	err = filepath.Walk(dir, walker)
	if err != nil {
		t.Errorf("Error walking %v: %v", dir, err)
	}
}

func TestServeAsset404(t *testing.T) {
	os.RemoveAll(setup())

	global.assets = map[string]shared.Asset{}

	addr := ensureServe()

	err := awaitPing(addr)
	if err != nil {
		t.Errorf("Error while waiting for server to come up: %v", err)
		return
	}

	res, err := http.Get("http://" + addr + "/a/anything")
	if err != nil {
		t.Errorf("Unexpected error while trying to requesting missing asset: %v", err)
		return
	}
	if res.StatusCode != 404 {
		t.Errorf("Expected 404 when accessing missing asset got %+v", res)
		return
	}
}

func TestServeRootFromAsset(t *testing.T) {
	os.RemoveAll(setup())

	dir := "testdata/assets"
	assets, err := shared.LoadAssets(dir)
	if err != nil {
		t.Errorf("Error while loading assets from %v: %v", dir, err)
		return
	}
	global.assets = assets

	addr := ensureServe()
	err = awaitPing(addr)
	if err != nil {
		t.Errorf("Error while waiting for server to come up: %v", err)
		return
	}

	res, err := http.Get("http://" + addr + "/")
	if err != nil {
		t.Errorf("Error requesting root: %v", err)
		return
	}

	if res.StatusCode != 200 {
		t.Errorf("Expected 200 when accessing root got %+v", res)
		return
	}
}
