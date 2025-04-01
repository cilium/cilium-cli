#!/usr/bin/env bash

set -ex

CILIUM_CLI_IMAGE_REPO=${CILIUM_CLI_IMAGE_REPO:-quay.io/cilium/cilium-cli-ci}
CILIUM_CLI_IMAGE_TAG=${CILIUM_CLI_IMAGE_TAG:-latest}
KUBECONFIG=${KUBECONFIG:-~/.kube/config}
KUBECONFIG_MOUNT_PATH=${KUBECONFIG_MOUNT_PATH:-/root/.kube/config}

docker run \
  --network host \
  -v "$KUBECONFIG":"$KUBECONFIG_MOUNT_PATH" \
  -v "$(pwd)":/root/app \
  -e GITHUB_WORKFLOW_REF="$GITHUB_WORKFLOW_REF" \
  "$CILIUM_CLI_IMAGE_REPO":"$CILIUM_CLI_IMAGE_TAG" cilium "$@"
