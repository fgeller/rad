export SHELL:=/bin/bash -O extglob -c
export GO15VENDOREXPERIMENT:=1

build:
	rm -rf assets
	mkdir -p assets
	cp -r ui/!(rad.js) assets/
	babel --out-file assets/rad.js ui/rad.js
	go run asset.go global.go build.go

test:
	go test !(build*).go
