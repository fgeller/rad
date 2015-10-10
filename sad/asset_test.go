package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestLoadAssets(t *testing.T) {
	global.assets = map[string]asset{}
	dir, err := filepath.Abs("ui")
	if err != nil {
		t.Errorf("Unexpected error while turning into absolute path: %v", err)
		return
	}

	ui, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Errorf("Unexpected error while reading ui folder's contents: %v", err)
		return
	}
	if len(ui) == 0 {
		t.Errorf("Expected contents in ui folder, but got nothing.")
		return
	}

	err = loadAssets(dir)
	if err != nil {
		t.Errorf("Unexpected error while loading assets: %v", err)
		return
	}

	isPresent := func(n string) bool {
		actual, ok := global.assets[n]
		if !ok {
			t.Errorf("Expected %v to be present in loaded assets", n)
			fmt.Printf("assets:\n")
			for k, v := range global.assets {
				fmt.Printf("%v: %v\n", k, v)
			}
			return false
		}

		if len(actual.contentType) == 0 || len(actual.content) == 0 {
			return false
		}

		return true
	}

loop:
	for _, fi := range ui {
		if !fi.IsDir() {
			if !isPresent(fi.Name()) {
				return
			}
			continue loop
		}

		d := filepath.Join(dir, fi.Name())
		nfs, err := ioutil.ReadDir(d)
		if err != nil {
			t.Errorf("Unexpected error reading dir %v: %v", d, err)
			return
		}

		for _, nf := range nfs {
			if !nf.IsDir() {
				rel, err := filepath.Rel(dir, filepath.Join(d, nf.Name()))
				if err != nil {
					t.Errorf("Unexpected error finding rel path: %v", err)
					return
				}
				if !isPresent(rel) {
					return
				}
			}
		}
	}
}
