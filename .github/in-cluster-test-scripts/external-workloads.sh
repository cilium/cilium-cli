#!/bin/sh

set -x
set -e

# Create a FIFO through which we'll pipe the sysdump before the job terminates.
# NOTE: This MUST be the first step in the script.
mkfifo /tmp/cilium-sysdump-out

# Run connectivity test
cilium connectivity test --all-flows

# Retrieve Cilium status
cilium status
cilium clustermesh status
cilium clustermesh vm status

# Grab a sysdump
cilium sysdump --output-filename=cilium-sysdump-out

# Wait for the sysdump to be read.
# NOTE: This MUST be the last step in the script.
cat cilium-sysdump-out >> /tmp/cilium-sysdump-out
