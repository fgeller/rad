export SHELL:=/bin/bash -O extglob -c

build:
	cd sad ; make build
	cd pad ; make build

test:
	cd shared ; make test
	cd pad ; make test
	cd sad ; make test
