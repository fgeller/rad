export SHELL:=/bin/bash -O extglob -c
export GO15VENDOREXPERIMENT:=1

build:
	go build -o sad +([^_]).go

test:
	go test *.go
