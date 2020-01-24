#!/usr/bin/env bash
echo "Publishing helm package to Artifactory"

export CHART_DIR=k8s/trans

helm lint ${CHART_DIR}
helm package ${CHART_DIR}
jfrog rt u "*.tgz" "helm-local/yapo/" || true
