package main

import (
	"../shared"

	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func awaitPing(addr string) error {
	limit := 10
	attempts := 0

	for {
		resp, err := http.Get("http://" + addr + "/ping")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		attempts++
		if attempts > limit {
			return fmt.Errorf("Got no ping on %v.", addr)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func TestGenerateAssetsSourceFile(t *testing.T) {
	defer resetGeneratedInfo()
	dir := "testdata/assets"
	config.assetsOut = "generated_assets.go"

	assets, err := shared.LoadAssets(dir)
	if err != nil {
		t.Errorf("Error while loading assets: %v", err)
		return
	}

	if err = generate(assets); err != nil {
		t.Errorf("Error while generating assets source: %v", err)
		return
	}

	cmd := exec.Command("gofmt", "-e", "-d", config.assetsOut)
	out, err := cmd.CombinedOutput()
	if err != nil || len(out) != 0 {
		t.Errorf("Expected generated file to be well-formed, err %v, out:\n%s", err, out)
		return
	}
}

func TestBuildBinaryWithAssets(t *testing.T) {
	defer resetGeneratedInfo()
	out, err := filepath.Abs("test-sad")

	if err != nil {
		t.Errorf("Error finding absolute path: %v", err)
		return
	}
	defer os.RemoveAll(out)
	assets := "testdata/assets"

	err = build(out, assets)
	if err != nil {
		t.Errorf("Error building binary: %v", err)
		return
	}

	tmpPackDir, err := ioutil.TempDir("", "sad-pack-dir")
	if err != nil {
		t.Errorf("Error creating temporary directory: %v", err)
		return
	}

	addr := "localhost:6072"
	cmd := exec.Command(out, "-packdir", tmpPackDir, "-addr", addr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Errorf("Error accessing stdout: %v", err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Errorf("Error accessing stderr: %v", err)
		return
	}
	defer func() {
		cmd.Process.Kill()
		output, _ := ioutil.ReadAll(stdout)
		fmt.Printf("stdout:\n%s\n", output)
		errout, _ := ioutil.ReadAll(stderr)
		fmt.Printf("stderr:\n%s\n", errout)
	}()

	err = cmd.Start()
	if err != nil {
		t.Errorf("Error starting binary: %v", err)
		return
	}

	err = awaitPing(addr)
	if err != nil {
		t.Errorf("Error waiting for server to come up: %v", err)
		return
	}

	walker := func(p string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			return err
		}

		rel, err := filepath.Rel(assets, p)
		if err != nil {
			return err
		}

		res, err := http.Get("http://" + addr + "/a/" + rel)
		if err != nil {
			t.Errorf("Error requesting asset %v: %v", rel, err)
			return err
		}

		if res.StatusCode != 200 {
			t.Errorf("Non 200 status code when requesting asset %v: %+v", rel, res)
			return err
		}

		return nil
	}

	err = filepath.Walk(assets, walker)
	if err != nil {
		t.Errorf("Error walking asset files: %v", err)
		return
	}

}
