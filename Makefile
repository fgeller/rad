export SHELL:=/bin/bash -O extglob -c

build:
	cd pad ; make build
	cd sad ; make build
	cd sap ; make build

test:
	cd shared ; make test
	cd pad ; make test
	cd sad ; make test
	cd sap ; make test
