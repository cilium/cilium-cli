# syntax=docker/dockerfile:1.11@sha256:10c699f1b6c8bdc8f6b4ce8974855dd8542f1768c26eb240237b8f1c9c6c9976

# Copyright Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

FROM docker.io/library/golang:1.23.3-alpine3.19@sha256:f72297ec1cf35152ecfe7a4d692825fc608fea8f3d3fa8f986fda70184082823 AS builder
WORKDIR /go/src/github.com/cilium/cilium-cli
RUN apk add --no-cache curl git make ca-certificates
COPY . .
RUN make

# cilium-cli is from scratch only including cilium binaries
FROM scratch AS cilium-cli
ENTRYPOINT ["cilium"]
LABEL maintainer="maintainer@cilium.io"
WORKDIR /root/app
COPY --from=builder --chown=root:root --chmod=755 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/cilium/cilium-cli/cilium /usr/local/bin/cilium

# cilium-cli-ci is based on ubuntu with cloud CLIs
FROM ubuntu:24.04@sha256:99c35190e22d294cdace2783ac55effc69d32896daaa265f0bbedbcde4fbe3e5 AS cilium-cli-ci
ENTRYPOINT []
LABEL maintainer="maintainer@cilium.io"
WORKDIR /root/app
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
