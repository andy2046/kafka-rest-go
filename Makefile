.PHONY: all help test

all: help

help:				## Show this help
	@scripts/help.sh

test:				## Test potential bugs and race conditions
	@scripts/test.sh
