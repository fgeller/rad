package main

import (
	"../shared"

	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func fileHeader(imports []string) string {
	out := `package main

`
	if len(imports) > 0 {
		qi := []string{}
		for _, i := range imports {
			qi = append(qi, fmt.Sprintf("\"%v\"", i))
		}
		out += `import (
	` + strings.Join(qi, `	
`) + `)
`
	}

	return out
}

func registerBuildVersion() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error identifying version err=%v.\n", err)
	}

	v := out.String()[:7]
	log.Printf("Identified build version [%v].\n", v)

	str := `

func registerBuildVersion() {
	global.buildVersion = "` + v + `"
}
`

	return str
}

func writeGeneratedInfo(contents string) error {
	log.Printf("Generating assets in %v\n", config.assetsOut)
	return ioutil.WriteFile(config.assetsOut, []byte(contents), 0755)
}

func registerAssets(assets map[string]shared.Asset) string {

	tmpl := `

func registerAssets() {
`
	for rel, a := range assets {
		tmpl += `
	global.assets["` + rel + `"] = shared.Asset{
		ContentType: "` + a.ContentType + `",
		Content:     `
		tmpl += fmt.Sprintf("%#v", a.Content)
		tmpl += `,
	}
`
	}

	tmpl += `
}
`

	return tmpl
}

func resetGeneratedInfo() {
	err := ioutil.WriteFile(
		config.assetsOut,
		[]byte(`package main
func registerBuildVersion() {}
func registerAssets() {}
`),
		0755,
	)
	log.Printf("Reset generated file (err: %v).", err)
}

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
			(n != "prof.go" || config.prof) &&
			n != "build.go" {
			sources = append(sources, f.Name())
		}
	}

	return sources, nil
}

func build(out string, assetsDir string) error {
	start := time.Now()
	resetGeneratedInfo()

	assets, err := shared.LoadAssets(assetsDir)
	if err != nil {
		return err
	}
	log.Printf("Assets loaded in %v\n", time.Since(start))

	last := time.Now()
	contents := fileHeader([]string{"../shared"})
	contents += registerBuildVersion()
	contents += registerAssets(assets)
	if err = writeGeneratedInfo(contents); err != nil {
		return err
	}
	log.Printf("Generated info in %v\n", time.Since(last))

	sources, err := findSources()
	if err != nil {
		return err
	}
	args := []string{"build", "-v", "-o", out}
	args = append(args, sources...)
	cmd := exec.Command("go", args...)
	env := os.Environ()
	env = append(env, "GO15VENDOREXPERIMENT=1")
	env = append(env, "GOOS="+config.os)
	env = append(env, "GOARCH="+config.arch)
	cmd.Env = env

	if config.verbose {
		log.Printf("Building [%v] from %v.", out, sources)
		log.Printf("Build command:\n%v", cmd.Args)
		log.Printf("Build env:\n%v", cmd.Env)
	}

	log.Printf("Compiling... Get a coffee, or read http://xkcd.com/ :P\n")
	last = time.Now()
	output, err := cmd.CombinedOutput()
	log.Printf("Finished compiling after %v", time.Since(last))
	if config.verbose {
		log.Printf("Build combined output:\n%s", output)
	}
	if err != nil {
		return err
	}

	resetGeneratedInfo()
	log.Printf("Done, happy serving! (after %v)", time.Since(start))
	return nil
}

var config struct {
	out       string
	assets    string
	assetsOut string
	os        string
	arch      string
	verbose   bool
	prof      bool
}

func main() {
	flag.StringVar(&config.out, "out", "sad", "Name of generated binary")
	flag.StringVar(&config.assets, "assets", "assets", "Location of assets")
	flag.StringVar(&config.assetsOut, "assetsOut", "generated_assets.go", "File where assets are compiled.")
	flag.StringVar(&config.os, "os", "darwin", "GOOS to compile binary for.")
	flag.StringVar(&config.arch, "arch", "amd64", "GOARCH to compile binary for.")
	flag.BoolVar(&config.prof, "prof", false, "Enable prof output at /debug/pprof")
	flag.BoolVar(&config.verbose, "v", false, "Verbose output")
	flag.Parse()

	if config.os == "" {
		config.os = "darwin"
	}

	if config.arch == "" {
		config.arch = "amd64"
	}

	if config.os == "windows" && !strings.HasSuffix(config.out, ".exe") {
		config.out += ".exe"
	}

	log.Printf("Read config: %+v\n", config)
	err := build(config.out, config.assets)
	if err != nil {
		log.Fatalf("Error during build: %v", err)
	}
}
