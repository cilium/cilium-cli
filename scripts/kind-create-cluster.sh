#!/bin/bash

# Copyright 2017-2020 Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

# TODO: move this script to the kube-test image

set -o errexit
set -o pipefail
set -o nounset

TESTER_IMAGE="${TESTER_IMAGE:-docker.io/cilium/kube-test:35e73180328e7a2a5679a14527d6826eb60a6574}"

if [ "$#" -lt 3 ] ; then
  echo "$0 requires at least 3 arguments - version, config flavour & cluster name (any additionall arguments passed to 'kind create cluster' directly)"
  exit 1
fi

root_dir="$(git rev-parse --show-toplevel)"

if [ -z "${TESTER_CONTAINER+x}" ] ; then
   exec docker run --rm --volume /var/run/docker.sock:/var/run/docker.sock --volume "${root_dir}:/src" --workdir /src --env KUBECONFIG="/src/kubeconfig" "${TESTER_IMAGE}" "/src/scripts/$(basename "${0}")" "${@}"
fi

cd "${root_dir}"

kube_version="${1}"
config_flavour="${2}"
name="${3}"

shift 3

export KUBECONFIG="${KUBECONFIG:-"/github/workspace/kubeconfig"}"

config_path="/etc/kind/${kube_version}/${config_flavour}-cluster.yaml"

if ! [ -e "${config_path}" ] ; then
  echo "no config file exists for given version and flavour combination (${config_path}), but the following configs are available:"
  ls /etc/kind/*/*-cluster.yaml
  exit 2
fi

kind create cluster --config="${config_path}" --name="${name}" "${@}"

kubectl config use-context "kind-${name}"

chmod o+r "${KUBECONFIG}"
