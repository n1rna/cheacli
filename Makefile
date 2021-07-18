# Makefile for go program
.PHONY: build usage install test

PROGRAM=cheat
DEFAULT_CONFIG_DIRECTORY = ${SUDO_USER}/.config/cheat

build:
	go build .

usage:
	@echo "make [build|run|kill|test]"
	@echo "   - build : compile and build binary"
	@echo "   - run   : start the server and exec client"
	@echo "   - kill  : stop the server"

install:
	echo $(DEFAULT_CONFIG_DIRECTORY)
	mkdir -p $(DEFAULT_CONFIG_DIRECTORY)/repos
	sudo cp cheat /usr/local/bin
	cp .cheatclirc $(DEFAULT_CONFIG_DIRECTORY)
	sed -i 's|\.CONFIGDIR|$(DEFAULT_CONFIG_DIRECTORY)|g' $(DEFAULT_CONFIG_DIRECTORY)/.cheatclirc
	cheat init --config $(DEFAULT_CONFIG_DIRECTORY)/.cheatclirc

test:
	cd test && go test -v
