package main

import (
	"../shared"

	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type scanResult struct {
	namespaces     []shared.Namespace
	processedFiles int
}

func scan(path string, p parser) ([]shared.Namespace, error) {
	start := time.Now()
	var ns []shared.Namespace
	ap, err := filepath.Abs(path)
	if err != nil {
		return ns, err
	}

	fc, ns, err := scanDir(ap, p)
	elapsed := time.Now().Sub(start)
	log.Printf("Found %v entries (%.1ff/s).\n", len(ns), float64(fc)/elapsed.Seconds())

	return shared.Merge(ns), err
}

func isMarkupFile(p string) bool {
	lp := strings.ToLower(p)
	return strings.HasSuffix(lp, "html") || strings.HasSuffix(lp, "xml")
}

func scanDir(dir string, p parser) (int, []shared.Namespace, error) {
	var namespaces []shared.Namespace
	results := make(chan []shared.Namespace)
	counts := make(chan int)
	var wg sync.WaitGroup

	walker := func(path string, info os.FileInfo, err error) error {
		wg.Add(1)
		if err != nil || info.IsDir() || !isMarkupFile(info.Name()) {
			wg.Done()
			return err
		}

		go func() {
			results <- scanFile(path, p)
			counts <- 1
			wg.Done()
		}()

		return err
	}

	err := filepath.Walk(dir, walker)
	if err != nil {
		return -1, namespaces, err
	}

	go func() { wg.Wait(); close(results) }()

	count := 0
wait:
	for {
		select {
		case new, ok := <-results:
			if !ok {
				break wait // we're done
			}
			namespaces = append(namespaces, new...)
		case inc := <-counts:
			count += inc
		}
	}

	return count, namespaces, err
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
