#!/bin/sh

set -x
set -e

# Run connectivity test
cilium connectivity test --all-flows

# Retrieve Cilium status
cilium status
cilium clustermesh status
cilium clustermesh vm status

# Grab a sysdump and move it to the persistent volume.
cilium sysdump --output-filename cilium-sysdump-out
mv cilium-sysdump-out.zip /output/cilium-sysdump-out.zip
