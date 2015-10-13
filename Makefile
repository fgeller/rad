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

builddocker:
	docker build -t rad/build-sap -f ./Dockerfile.build .
	docker run -t rad/build-sap /bin/true
	docker cp `docker ps -q -n=1`:/sap .
	chmod 755 ./sap
	docker build --rm=true --tag=rad/sap -f Dockerfile.static .

rundocker: builddocker
	docker run -p 3025:3025 -v ${HOME}/sap-packs:/packs -t rad/sap
