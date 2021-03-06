export SHELL:=/bin/bash -O extglob -c
export GO15VENDOREXPERIMENT:=1

build: BUILD_OS ?= darwin
build: BUILD_ARCH ?= amd64
build: BUILD_PROF ?= false
build:
	rm -vf sad
	rm -rf assets
	mkdir -p assets
	cp -r ui/* assets/
	go run build.go -v -prof=${BUILD_PROF} -os ${BUILD_OS} -arch ${BUILD_ARCH}

test: buildtest apptest

apptest:
	go test !(build*).go

buildtest:
	go test build*.go

snapshot:
	BUILD_PROF=true $(MAKE) build
	tar Jcf sad-`date -I`-`git rev-parse HEAD | cut -c-7`-`uname -s -m | tr ' ' '-'`.xz sad

release-windows:
	BUILD_OS=windows BUILD_ARCH=amd64 $(MAKE) build
	tar Jcf sad-`git describe --abbrev=0 --tags`-win-amd64.xz sad.exe

release-linux:
	BUILD_OS=linux BUILD_ARCH=amd64 $(MAKE) build
	tar Jcf sad-`git describe --abbrev=0 --tags`-linux-amd64.xz sad

release-darwin:
	BUILD_OS=darwin BUILD_ARCH=amd64 $(MAKE) build
	tar Jcf sad-`git describe --abbrev=0 --tags`-darwin-amd64.xz sad

release: release-linux release-darwin release-windows

run:
	./sad -sapaddr 0.0.0.0:3025 -devmode=true
