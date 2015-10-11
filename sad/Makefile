export SHELL:=/bin/bash -O extglob -c
export GO15VENDOREXPERIMENT:=1

build:
	go run asset.go global.go build.go -v

test:
	go test !(build*).go
