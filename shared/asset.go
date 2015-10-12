package shared

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Asset struct {
	ContentType string
	Content     []byte
}

func (a Asset) String() string {
	return fmt.Sprintf(
		"Asset{ContentType: %v, len(Content)=%v}",
		a.ContentType,
		len(a.Content),
	)
}

func detectContentType(p string) string {
	if strings.HasSuffix(p, ".css") {
		return "text/css; charset=utf-8"
	}

	if strings.HasSuffix(p, ".js") {
		return "application/javascript"
	}

	if strings.HasSuffix(p, ".html") {
		return "text/html; charset=utf-8"
	}

	return "application/octet-stream"
}

func LoadAssets(dir string) (map[string]Asset, error) {
	log.Printf("Loading assets in %v\n", dir)
	assets := map[string]Asset{}

	fi, err := os.Stat(dir)
	if err != nil || !fi.IsDir() {
		return assets, fmt.Errorf("Expected directory, err: %v", err)
	}

	walker := func(p string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			return err
		}

		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}

		ctype := detectContentType(p)
		ctnt, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}

		a := Asset{
			Content:     ctnt,
			ContentType: ctype,
		}
		assets[rel] = a
		log.Printf("Loaded asset [%v] as %v.", rel, a)

		return nil
	}

	return assets, filepath.Walk(dir, walker)
}
