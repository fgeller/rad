package shared

import (
	"archive/zip"

	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type downloader func(string) (*http.Response, error)

func unzipFile(f *zip.File, dest string) error {

	path := filepath.Join(dest, f.Name)
	if f.FileInfo().IsDir() {
		return os.MkdirAll(path, f.Mode())
	}

	if !FileExists(filepath.Dir(path)) {
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return err
		}
	}
	fc, err := f.Open()
	if err != nil {
		return err
	}
	defer fc.Close()

	path = MaybeEscapeWinPath(path)
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, fc)
	return err
}

func Unzip(src string, dest string) error {
	err := os.MkdirAll(dest, 0755)
	if err != nil {
		return err
	}

	log.Printf("Unzipping %v to %v\n", src, dest)

	r, err := zip.OpenReader(src)
	if err != nil {
		log.Printf("failed to open zip: %v", src)
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		err = unzipFile(f, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func Download(remote string) (string, error) {
	local := remote[strings.LastIndex(remote, "/")+1:]
	if FileExists(local) {
		log.Printf("Already downloaded [%v].", local)
		return local, nil
	}

	out, err := os.Create(local)
	if err != nil {
		return "", err
	}
	defer out.Close()

	return local, DownloadToLocal(remote, out)
}

func DownloadToTemp(remote string) (string, error) {
	f, err := ioutil.TempFile("", "temp-download")
	if err != nil {
		return "", err
	}

	log.Printf("Downloading %v to temp: %v\n", remote, f.Name())

	return f.Name(), DownloadToLocal(remote, f)
}

func DownloadToLocal(remote string, local *os.File) error {
	resp, err := http.Get(remote)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to download, got status code: %v", resp.StatusCode)
	}

	log.Printf("Downloading [%v] to local [%v].\n", remote, local.Name())

	n, err := io.Copy(local, resp.Body)
	log.Printf("Downloaded %v bytes.\n", n)

	return err
}

func ZipDir(out *os.File, in string) error {

	w := zip.NewWriter(out)
	tld := filepath.Dir(in)

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		rel, err := filepath.Rel(tld, path)
		if err != nil {
			return err
		}

		f, err := w.Create(rel)
		if err != nil {
			return err
		}

		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		_, err = f.Write(contents)
		return err
	}

	err := filepath.Walk(in, walker)
	if err != nil {
		return err
	}

	err = w.Close()
	return err
}

func CopyDir(source string, dest string) (int, error) {

	// TODO: ensure src/dst are directories

	c := 0
	absSource, err := filepath.Abs(source)
	if err != nil {
		return c, err
	}
	source = absSource

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Ignoring err=%v.\n", err)
		}

		target := filepath.Join(dest, path[len(filepath.Dir(source)):])
		if err == nil && info.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		target = MaybeEscapeWinPath(target)
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		defer out.Close()

		path = MaybeEscapeWinPath(path)
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		_, err = io.Copy(out, in)
		c++
		return err
	}

	return c, filepath.Walk(source, walker)
}

func MaybeEscapeWinPath(path string) string {
	if runtime.GOOS == "windows" {
		return "\\\\.\\" + path
	}

	return path
}
