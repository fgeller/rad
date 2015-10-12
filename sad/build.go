package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func findSources() ([]string, error) {
	var sources []string

	fs, err := ioutil.ReadDir(".")
	if err != nil {
		return sources, err
	}

	for _, f := range fs {
		n := f.Name()
		if strings.HasSuffix(n, ".go") &&
			strings.Index(n, "_test") < 0 &&
			n != "build.go" {
			sources = append(sources, f.Name())
		}
	}

	return sources, nil
}

func build(out string, assets string) error {
	start := time.Now()

	err := loadAssets(assets)
	if err != nil {
		return err
	}
	log.Printf("Assets loaded in %v\n", time.Since(start))

	last := time.Now()
	asset, err := writeAssets()
	if err != nil {
		return err
	}
	log.Printf("Wrote assets out in %v\n", time.Since(last))

	err = os.RemoveAll("generated_assets.go")
	if err != nil {
		return err
	}

	err = os.Link(asset, "generated_assets.go")
	if err != nil {
		return err
	}

	sources, err := findSources()
	if err != nil {
		return err
	}
	args := []string{"build", "-v", "-o", out}
	args = append(args, sources...)
	cmd := exec.Command("go", args...)
	env := os.Environ()
	env = append(env, "GO15VENDOREXPERIMENT=1")
	cmd.Env = env

	if buildConfig.verbose {
		log.Printf("Building [%v] from %v.", out, sources)
		log.Printf("Build command:\n%v", cmd.Args)
		log.Printf("Build env:\n%v", cmd.Env)
	}

	log.Printf("Compiling... Get a coffee, or read http://xkcd.com/ :P\n")
	last = time.Now()
	output, err := cmd.CombinedOutput()
	log.Printf("Finished compiling after %v", time.Since(last))
	if buildConfig.verbose {
		log.Printf("Build combined output:\n%s", output)
	}
	if err != nil {
		return err
	}

	resetGeneratedAssets()
	log.Printf("Done, happy serving! (after %v)", time.Since(start))
	return nil
}

var buildConfig struct {
	out     string
	assets  string
	verbose bool
}

func main() {
	flag.StringVar(&buildConfig.out, "out", "sad", "Name of generated binary")
	flag.StringVar(&buildConfig.assets, "assets", "assets", "Location of assets")
	flag.BoolVar(&buildConfig.verbose, "v", false, "Verbose output")
	flag.Parse()

	resetGlobals()

	log.Printf("Read config: %+v\n", buildConfig)
	err := build(buildConfig.out, buildConfig.assets)
	if err != nil {
		log.Fatalf("Error during build: %v", err)
	}
}
