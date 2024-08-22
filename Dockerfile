# syntax=docker/dockerfile:1.10@sha256:865e5dd094beca432e8c0a1d5e1c465db5f998dca4e439981029b3b81fb39ed5

# Copyright Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

# FINAL_CONTAINER specifies the source for the output
# cilium-cli-ci (default) is based on ubuntu with cloud CLIs
# cilium-cli is from scratch only including cilium binaries
ARG FINAL_CONTAINER="cilium-cli-ci"

FROM docker.io/library/golang:1.23.1-alpine3.19@sha256:e0ea2a119ae0939a6d449ea18b2b1ba30b44986ec48dbb88f3a93371b4bf8750 AS builder
WORKDIR /go/src/github.com/cilium/cilium-cli
RUN apk add --no-cache git make ca-certificates
COPY . .
RUN make

FROM scratch AS cilium-cli
COPY --from=builder /go/src/github.com/cilium/cilium-cli/cilium /usr/local/bin/cilium

FROM ubuntu:24.04@sha256:56a8952801afd93876eea675cae9ab861bf8d2e6a4f978e4b0237ce94e1c3b49 AS cilium-cli-ci
COPY --from=builder /go/src/github.com/cilium/cilium-cli/cilium /usr/local/bin/cilium

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

# Select the layer to provide the final container image from
FROM ${FINAL_CONTAINER}
LABEL maintainer="maintainer@cilium.io"
WORKDIR /root/app
ENTRYPOINT []
