#!/bin/bash

export CHART_DIR=k8s/trans

helm lint ${CHART_DIR}
helm package ${CHART_DIR} --version 0.1.${TRAVIS_BUILD_NUMBER}
jfrog rt u "*.tgz" "helm-local/yapo/"
