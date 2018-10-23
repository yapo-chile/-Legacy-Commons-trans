#!/usr/bin/env bash
export UNAMESTR = $(uname)
export GO_FILES = $(shell find . -iname '*.go' -type f | grep -v vendor | grep -v pact) # All the .go files, excluding vendor/ and pact/
GENPORTOFF?=0
genport = $(shell expr ${GENPORTOFF} + \( $(shell id -u) - \( $(shell id -u) / 100 \) \* 100 \) \* 200 + 30100 + $(1))

# BRANCH info from travis
export BUILD_BRANCH=$(shell if [ "${TRAVIS_PULL_REQUEST}" = "false" ]; then echo "${TRAVIS_BRANCH}"; else echo "${TRAVIS_PULL_REQUEST_BRANCH}"; fi)

# GIT variables
export BRANCH=$(shell git branch | sed -n 's/^\* //p')
export GIT_BRANCH=$(shell if [ -n "${BUILD_BRANCH}" ]; then echo "${BUILD_BRANCH}"; else echo "${BRANCH}"; fi;)
export GIT_COMMIT=$(shell git rev-parse HEAD)
export GIT_COMMIT_DATE=$(shell TZ="America/Santiago" git show --quiet --date='format-local:%d-%m-%Y_%H:%M:%S' --format="%cd")
export BUILD_CREATOR=$(shell git log --format=format:%ae | head -n 1)

# REPORT_ARTIFACTS should be in sync with `RegexpFilePathMatcher` in
# `reports-publisher/config.json`
export REPORT_ARTIFACTS=reports

# APP variables
# This variables are for the use of your microservice. This variables must be updated each time you are creating a new microservice
export APPNAME=goms
export YO=`whoami`
export SERVICE_HOST=:$(call genport,2)
export SERVER_ROOT=${PWD}
export BASE_URL="http://${SERVICE_HOST}"
export MAIN_FILE=cmd/${APPNAME}/main.go
export LOGGER_SYSLOG_ENABLED=false
export LOGGER_STDLOG_ENABLED=true
export LOGGER_LOG_LEVEL=0

# Pact test variables
export PACT_MAIN_FILE=cmd/${APPNAME}-pact/main.go
export PACT_BINARY=${APPNAME}-pact

# DOCKER variables
export DOCKER_REGISTRY=containers.schibsted.io
export DOCKER_IMAGE=${DOCKER_REGISTRY}/yapo/${APPNAME}
export DOCKER_PORT=$(call genport,1)

# Documentation variables
export DOCS_DIR=docs
export DOCS_HOST=localhost:$(call genport,3)
export DOCS_PATH=github.schibsted.io/Yapo/${APPNAME}
export DOCS_COMMIT_MESSAGE=Generate updated documentation
