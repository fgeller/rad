package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

func detectContentType(p string) (string, error) {
	content, err := os.Open(p)
	if err != nil {
		return "", err
	}
	defer content.Close()

	var buf [512]byte
	n, _ := io.ReadFull(content, buf[:])
	ctype := http.DetectContentType(buf[:n])

	// rewind to output whole file
	if _, err := content.Seek(0, os.SEEK_SET); err != nil {
		return ctype, err
	}
	return ctype, nil
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

		ctype, err := detectContentType(p)
		ctnt, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}

		log.Printf("Compiled %v with content type [%v]\n", rel, ctype)
		global.assets[rel] = asset{
			content:     ctnt,
			contentType: ctype,
		}

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
