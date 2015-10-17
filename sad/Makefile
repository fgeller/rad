export SHELL:=/bin/bash -O extglob -c
export GO15VENDOREXPERIMENT:=1

build:
	rm -rf assets
	mkdir -p assets
	cp -r ui/!(rad.js) assets/
	babel --out-file assets/rad.js ui/rad.js
	go run build.go -v

test: buildtest apptest

apptest:
	go test !(build*).go

buildtest:
	go test build*.go

release: build
	tar Jcf sad-`date -I`-`git rev-parse HEAD | cut -c-7`-`uname -s -m | tr ' ' '-'`.xz sad

run:
	./sad
