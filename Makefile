.PHONY: all help build dep

all: help

help:				## Show this help
	@scripts/help.sh

build:				## Run docker-compose to build image
	@scripts/build.sh

dep: 				## Get the dependencies
	@dep status
