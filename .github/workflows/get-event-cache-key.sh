#!/bin/bash

# Copyright 2017-2020 Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o pipefail
set -o nounset

if [ -z "${GITHUB_ACTIONS+x}" ] ; then
  exit 1
fi

# use git hash-object purelly out of convenient, as e.g. with sha256sum output needs to be parsed
event_hash="$(git hash-object "${GITHUB_EVENT_PATH}")"
printf "::set-output name=key::%s\n" "${event_hash}"
