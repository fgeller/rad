package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type asset struct {
	contentType string
	content     []byte
}

func (a asset) String() string {
	return fmt.Sprintf(
		"asset{contentType: %v, len(content)=%v}",
		a.contentType,
		len(a.content),
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

func loadAssets(dir string) error {
	log.Printf("Loading assets in %v\n", dir)

	fi, err := os.Stat(dir)
	if err != nil || !fi.IsDir() {
		return fmt.Errorf("Expected directory, err: %v", err)
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

		a := asset{
			content:     ctnt,
			contentType: ctype,
		}
		global.assets[rel] = a
		log.Printf("Loading asset [%v] as %v.", rel, a)

		return nil
	}

	return filepath.Walk(dir, walker)
}

// TODO: expect out
func writeAssets() (string, error) {
	tmp, err := ioutil.TempFile("", "assets")
	if err != nil {
		return "", err
	}
	defer tmp.Close()
	log.Printf("Generating assets in %v\n", tmp.Name())

	tmpl := `package main

func registerAssets() {
`
	for rel, a := range global.assets {
		tmpl += `
	global.assets["` + rel + `"] = asset{
		contentType: "` + a.contentType + `",
		content:     `
		tmpl += fmt.Sprintf("%#v", a.content)
		tmpl += `,
	}
`
	}

	tmpl += `
}
`
	_, err = io.WriteString(tmp, tmpl)
	return tmp.Name(), err
}

func resetGeneratedAssets() {
	err := ioutil.WriteFile(
		"generated_assets.go",
		[]byte(`package main
func registerAssets() {}
`),
		0755,
	)
	log.Printf("Reset generated assets file (err: %v).", err)
}
