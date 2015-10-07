package shared

import (
	"archive/zip"

	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

func Download(d downloader, remote string) (string, error) {
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

// TODO: CopyDir?
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
			return err
		}

		target := filepath.Join(dest, path[len(filepath.Dir(source)):])
		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		out, err := os.Create(target)
		if err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(out, in)
		c++
		return err
	}

	return c, filepath.Walk(source, walker)
}
