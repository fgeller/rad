package main

import (
	"../shared"
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

type scanResult struct {
	entries        []shared.Entry
	processedFiles int
}

func mkPath(parts ...string) string {
	result := ""
	for i, p := range parts {
		if i != len(parts)-1 {
			result += p + string(os.PathSeparator)
		} else {
			result += p
		}
	}
	return result
}

func scan(path string, p parser) ([]shared.Entry, error) {
	start := time.Now()
	fc, es, err := scanDir(path, p)
	elapsed := time.Now().Sub(start)
	log.Printf("found %v links (%.1ff/s).\n", len(es), float64(fc)/elapsed.Seconds())

	return es, err
}

func scanDir(dir string, p parser) (int, []shared.Entry, error) {

	files, err := findDirsAndMarkupFiles(dir)
	if err != nil {
		fmt.Printf("can't read dir %v, err %v\n", dir, err)
		return 0, []shared.Entry{}, err
	}

	rc := make(chan scanResult)
	runtime.GOMAXPROCS(runtime.NumCPU())

	for _, fi := range files {
		go func(dir string, f os.FileInfo, c chan scanResult) {
			path := dir + string(os.PathSeparator) + f.Name()
			switch {
			case f.IsDir():
				fs, es, _ := scanDir(path, p)
				c <- scanResult{es, fs}
			default:
				c <- scanResult{scanFile(path, p), 1}
			}
		}(dir, fi, rc)
	}

	results := []shared.Entry{}
	fc := 0
	for i := 0; i < len(files); i++ {
		r := <-rc
		fc += r.processedFiles
		results = append(results, r.entries...)
	}

	return fc, results, nil
}

func scanFile(path string, p parser) []shared.Entry {
	r, err := os.Open(path)
	defer r.Close()
	if err != nil {
		fmt.Printf("can't open file %v, err %v\n", path, err)
		return []shared.Entry{}
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

func unzip(src string, dest string) error {
	err := os.MkdirAll(dest, 0755)
	if err != nil {
		return err
	}

	r, err := zip.OpenReader(src)
	if err != nil {
		log.Printf("failed to open zip: %v", src)
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		path := dest + string(os.PathSeparator) + f.Name
		if f.FileInfo().IsDir() {
			os.Mkdir(path, f.Mode())
			continue
		}

		fc, err := f.Open()

		if err != nil {
			return err
		}

		dst, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		_, err = io.Copy(dst, fc)
		if err != nil {
			return err
		}
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func download(d downloader, remote string) (string, error) {
	local := remote[strings.LastIndex(remote, "/")+1:]
	if fileExists(local) {
		log.Printf("Already downloaded [%v].", local)
		return local, nil
	}

	out, err := os.Create(local)
	if err != nil {
		return "", err
	}
	defer out.Close()

	resp, err := http.Get(remote)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	log.Printf("Downloading [%v] to local [%v].\n", remote, local)

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	log.Printf("Downloaded %v bytes.\n", n)

	return local, nil
}
