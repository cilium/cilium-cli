#!/bin/sh

set -x
set -e

# Create a FIFO through which we'll pipe the sysdump before the job terminates.
# NOTE: This MUST be the first step in the script.
mkfifo /tmp/cilium-sysdump-out

# Enable Relay
cilium hubble enable

# Wait for Cilium status to be ready
cilium status --wait

# Port forward Relay
cilium hubble port-forward&
sleep 10s

# Run connectivity test
cilium connectivity test --all-flows

# Retrieve Cilium  status
cilium status

# Grab a sysdump
cilium sysdump --output-filename=cilium-sysdump-out

# Wait for the sysdump to be read.
# NOTE: This MUST be the last step in the script.
cat cilium-sysdump-out >> /tmp/cilium-sysdump-out
