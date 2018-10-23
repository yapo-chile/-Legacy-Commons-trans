include scripts/commands/vars.mk

## Run tests and generate quality reports
test:
	@scripts/commands/test.sh

## Run tests and output coverage reports
cover:
	@scripts/commands/test_cover.sh cli

## Run tests and open report on default web browser
coverhtml:
	@scripts/commands/test_cover.sh html

## Run gometalinter and output report as text
checkstyle:
	@scripts/commands/test_style.sh display

## Install golang system level dependencies
setup:
	@scripts/commands/setup.sh

## Compile the code
build:
	@scripts/commands/build.sh

## Execute the service
run:
	@./${APPNAME}

## Compile and start the service
start: build run

## Compile and start the service using docker
docker-start: build docker-build docker-compose-up info

## Stop docker containers
docker-stop: docker-compose-down

## Setup a new service repository based on goms
clone:
	@scripts/commands/clone.sh

## Run gofmt to reindent source
fix-format:
	@scripts/commands/fix-format.sh

## Display basic service info
info:
	@echo "YO           : ${YO}"
	@echo "ServerRoot   : ${SERVER_ROOT}"
	@echo "API Base URL : ${BASE_URL}"
	@echo "Healthcheck  : curl ${BASE_URL}/api/v1/healthcheck"

include docs.mk
include docker.mk
include help.mk
