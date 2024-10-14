# syntax=docker/dockerfile:1.10@sha256:865e5dd094beca432e8c0a1d5e1c465db5f998dca4e439981029b3b81fb39ed5

# Copyright Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

# FINAL_CONTAINER specifies the source for the output
# cilium-cli-ci (default) is based on ubuntu with cloud CLIs
# cilium-cli is from scratch only including cilium binaries
ARG FINAL_CONTAINER="cilium-cli-ci"

FROM docker.io/library/golang:1.23.2-alpine3.19@sha256:f6392ffebb028fed5ffe743ddb9716e38402c978779edd66474bb5d05f5e65e4 AS builder
WORKDIR /go/src/github.com/cilium/cilium-cli
RUN apk add --no-cache curl git make ca-certificates
COPY . .
RUN make

FROM scratch AS cilium-cli
COPY --from=builder /go/src/github.com/cilium/cilium-cli/cilium /usr/local/bin/cilium
ENTRYPOINT ["cilium"]

FROM ubuntu:24.04@sha256:ab64a8382e935382638764d8719362bb50ee418d944c1f3d26e0c99fae49a345 AS cilium-cli-ci
COPY --from=builder /go/src/github.com/cilium/cilium-cli/cilium /usr/local/bin/cilium
ENTRYPOINT []

# Install cloud CLIs. Based on these instructions:
# - https://cloud.google.com/sdk/docs/install#deb
# - https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html
# - https://learn.microsoft.com/en-us/cli/azure/install-azure-cli-linux?pivots=apt#install-azure-cli
RUN apt-get update -y \
  && apt-get install -y curl gnupg unzip \
  && curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | gpg --dearmor -o /usr/share/keyrings/cloud.google.gpg \
  && curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - \
  && echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list \
  && apt-get update -y \
  && apt-get install -y google-cloud-cli google-cloud-sdk-gke-gcloud-auth-plugin kubectl \
  && curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" \
  && unzip awscliv2.zip \
  && ./aws/install \
  && rm -r ./aws awscliv2.zip \
  && curl -sL https://aka.ms/InstallAzureCLIDeb | bash

FROM ${FINAL_CONTAINER} 
LABEL maintainer="maintainer@cilium.io"
WORKDIR /root/app
