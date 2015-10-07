package main

import (
	"../shared"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type scanResult struct {
	namespaces     []shared.Namespace
	processedFiles int
}

func scan(path string, p parser) ([]shared.Namespace, error) {
	start := time.Now()
	fc, ns, err := scanDir(path, p)
	elapsed := time.Now().Sub(start)
	log.Printf("Found %v entries (%.1ff/s).\n", len(ns), float64(fc)/elapsed.Seconds())

	return shared.Merge(ns), err
}

// TODO: use filepath.Walk
func scanDir(dir string, p parser) (int, []shared.Namespace, error) {

	files, err := findDirsAndMarkupFiles(dir)
	if err != nil {
		fmt.Printf("can't read dir %v, err %v\n", dir, err)
		return 0, []shared.Namespace{}, err
	}

	rc := make(chan scanResult)
	runtime.GOMAXPROCS(runtime.NumCPU())

	for _, fi := range files {
		go func(dir string, f os.FileInfo, c chan scanResult) {
			path := dir + string(os.PathSeparator) + f.Name() // TODO
			switch {
			case f.IsDir():
				fs, ns, _ := scanDir(path, p)
				c <- scanResult{ns, fs}
			default:
				c <- scanResult{scanFile(path, p), 1}
			}
		}(dir, fi, rc)
	}

	results := []shared.Namespace{}
	fc := 0
	for i := 0; i < len(files); i++ {
		r := <-rc
		fc += r.processedFiles
		results = append(results, r.namespaces...)
	}

	return fc, results, nil
}

func scanFile(path string, p parser) []shared.Namespace {
	r, err := os.Open(path)
	defer r.Close()
	if err != nil {
		fmt.Printf("can't open file %v, err %v\n", path, err)
		return []shared.Namespace{}
	}

	return p(path, r)
}

func findDirsAndMarkupFiles(dir string) ([]os.FileInfo, error) {
	files := []os.FileInfo{}

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("can't read dir %v, err %v\n", dir, err)
		return files, err
	}

	for _, f := range fs {
		if f.IsDir() ||
			strings.HasSuffix(f.Name(), "html") ||
			strings.HasSuffix(f.Name(), "xml") {
			files = append(files, f)
		}
	}

	return files, nil
}
