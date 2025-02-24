
NAME := bmc-iconik-io

.PHONY: build

build:
	./bin/build ${NAME}

clean:
	./bin/clean

readme:
	./bin/readme

lint:
	./bin/lint