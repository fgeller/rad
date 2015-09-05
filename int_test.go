package main

import "io/ioutil"
import "fmt"
import "os"
import "testing"
import "time"
import "reflect"
import "net/http"
import "encoding/json"

func TestInstallPack(t *testing.T) {

	e := entry{[]string{"main"}, "Entity", "Function", "Signature", "Target", "source"}
	es := []entry{e}
	indexer := func() ([]entry, error) { return es, nil }
	serveZip := func(addr string) {
		http.HandleFunc("/test.zip", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./testdata/test.zip")
		})
		http.ListenAndServe(addr, nil)
	}
	conf := pack{
		name:    "blubb",
		url:     "http://localhost:8881/test.zip",
		indexer: indexer,
	}

	// possible race condition here?
	go serveZip(":8881")
	defer os.RemoveAll("packs/" + conf.name)

	install(conf)

	actual, err := findEntityFunction(conf.name, "Entity", "", 10)
	if err != nil || len(actual) < 1 || !reflect.DeepEqual(e, actual[0]) {
		t.Errorf("Expected to find sample entry, got %v [err: %v].\n", actual, err)
	}

	dat, err := ioutil.ReadFile("packs/" + conf.name + "/rad-data.json")
	if err != nil {
		t.Errorf("unexpected error while reading serialized data: %v", err)
		return
	}

	exDat, err := json.Marshal(es)
	if string(exDat) != string(dat) {
		t.Errorf("expected serialized data\n%v\nbut got\n%v", err, exDat, string(dat))
	}

}

func TestFindEntityFunctions(t *testing.T) {
	pack := "test"
	entity := "Entity"
	samples := []entry{
		{[]string{"main"}, "AnotherEntity", "abc", "Signature", "Target", "source"},
		{[]string{"main"}, entity, "a", "Signature", "Target", "source"},
		{[]string{"main"}, entity, "ab", "Signature", "Target", "source"},
		{[]string{"main"}, entity, "abc", "Signature", "Target", "source"},
		{[]string{"main"}, entity, "d", "Signature", "Target", "source"},
		{[]string{"main"}, entity + "suffix", "x", "Signature", "Target", "source"},
		{[]string{"main"}, entity + "suffix", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_zzz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_zyz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_zzy", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yzz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_zyy", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yyz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yyy", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yzy", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yxz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_xyz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yyx", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_xxz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yxz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_xxx", "a", "Signature", "Target", "source"},
	}
	docs = map[string][]entry{}
	docs[pack] = samples

	addr := "0.0.0.0:3025"
	go serve(addr)
	time.Sleep(200)

	// find all for entity
	res, err := http.Get("http://" + addr + "/s?p=" + pack + "&e=" + entity)
	if err != nil {
		t.Errorf("unexpected error while finding entries: %v", err)
		return
	}

	byts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unexpected error while reading response body: %v", err)
		return
	}

	expectedEntries := []entry{
		{[]string{"main"}, entity, "a", "Signature", "Target", "source"},
		{[]string{"main"}, entity, "ab", "Signature", "Target", "source"},
		{[]string{"main"}, entity, "abc", "Signature", "Target", "source"},
		{[]string{"main"}, entity, "d", "Signature", "Target", "source"},
		{[]string{"main"}, entity + "suffix", "x", "Signature", "Target", "source"},
		{[]string{"main"}, entity + "suffix", "a", "Signature", "Target", "source"},
	}
	expected, err := json.Marshal(expectedEntries)
	if err != nil {
		t.Errorf("unexpected error while marshaling to json: %v", err)
		return
	}

	if string(byts) != string(expected) {
		t.Errorf("unexpected response, got\n%v\nbut expected\n%v\n", string(byts), string(expected))
	}

	// find all with given function prefix

	res, err = http.Get("http://" + addr + "/s?p=" + pack + "&e=" + entity + "&f=a")
	if err != nil {
		t.Errorf("unexpected error while finding entries: %v", err)
		return
	}

	byts, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unexpected error while reading response body: %v", err)
		return
	}

	expectedEntries = []entry{
		{[]string{"main"}, entity, "a", "Signature", "Target", "source"},
		{[]string{"main"}, entity, "ab", "Signature", "Target", "source"},
		{[]string{"main"}, entity, "abc", "Signature", "Target", "source"},
		{[]string{"main"}, entity + "suffix", "a", "Signature", "Target", "source"},
	}
	expected, err = json.Marshal(expectedEntries)
	if err != nil {
		t.Errorf("unexpected error while marshaling to json: %v", err)
		return
	}

	if string(byts) != string(expected) {
		t.Errorf("unexpected response, got \n%v\nbut expected\n%v\n", string(byts), string(expected))
	}

	// limit the results by default to 10
	res, err = http.Get("http://" + addr + "/s?p=" + pack + "&e=_")
	if err != nil {
		t.Errorf("unexpected error while finding entries: %v", err)
		return
	}

	byts, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unexpected error while reading response body: %v", err)
		return
	}

	expectedEntries = []entry{
		{[]string{"main"}, "_zzz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_zyz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_zzy", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yzz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_zyy", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yyz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yyy", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yzy", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_yxz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_xyz", "a", "Signature", "Target", "source"},
	}
	expected, err = json.Marshal(expectedEntries)
	if err != nil {
		t.Errorf("unexpected error while marshaling to json: %v", err)
		return
	}

	if string(byts) != string(expected) {
		t.Errorf("unexpected response, got \n%v\nbut expected\n%v\n", string(byts), string(expected))
	}

	// allow custom limit
	limit := 2
	res, err = http.Get(fmt.Sprintf("http://%v/s?p=%v&e=%v&limit=%v", addr, pack, "_", limit))
	if err != nil {
		t.Errorf("unexpected error while finding entries: %v", err)
		return
	}

	byts, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unexpected error while reading response body: %v", err)
		return
	}

	expectedEntries = []entry{
		{[]string{"main"}, "_zzz", "a", "Signature", "Target", "source"},
		{[]string{"main"}, "_zyz", "a", "Signature", "Target", "source"},
	}
	expected, err = json.Marshal(expectedEntries)
	if err != nil {
		t.Errorf("unexpected error while marshaling to json: %v", err)
		return
	}

	if string(byts) != string(expected) {
		t.Errorf("unexpected response, got \n%v\nbut expected\n%v\n", string(byts), string(expected))
	}
}
