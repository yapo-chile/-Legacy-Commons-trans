#!/usr/bin/env bash
echo "Publishing helm package to Artifactory"

export CHART_DIR=k8s/trans
export CHART_VERSION=$(grep version $CHART_DIR/Chart.yaml | awk '{print $2}')

helm lint ${CHART_DIR}
helm package ${CHART_DIR} --version ${CHART_VERSION}
jfrog rt u "*.tgz" "helm-local/yapo/" || true
