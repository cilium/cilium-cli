#!/bin/sh

set -x
set -e

# Install Cilium
cilium install \
  --cluster-name "${CLUSTER_NAME}" \
  --config monitor-aggregation=none
