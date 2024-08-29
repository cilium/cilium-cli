# syntax=docker/dockerfile:1.12@sha256:db1ff77fb637a5955317c7a3a62540196396d565f3dd5742e76dddbb6d75c4c5

# Copyright Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

FROM --platform=${BUILDPLATFORM} docker.io/library/golang:1.23.4-alpine3.19@sha256:5f3336882ad15d10ac1b59fbaba7cb84c35d4623774198b36ae60edeba45fd84 AS base
RUN apk add --no-cache --update ca-certificates git make 
WORKDIR /go/src/github.com/cilium/cilium-cli
COPY go.* .
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .

# xx is a helper for cross-compilation
FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.5.0@sha256:0c6a569797744e45955f39d4f7538ac344bfb7ebf0a54006a0a4297b153ccf0f AS xx

FROM --platform=${BUILDPLATFORM} base AS builder
ARG TARGETPLATFORM
ARG TARGETARCH
COPY --link --from=xx / /
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    xx-go --wrap && \
    make && \
    xx-verify --static /go/src/github.com/cilium/cilium-cli/cilium

# cilium-cli is from scratch only including cilium binaries
FROM --platform=${BUILDPLATFORM} scratch AS cilium-cli
ENTRYPOINT [""]
USER 1000:1000
LABEL maintainer="maintainer@cilium.io"
WORKDIR /root/app
COPY --link --from=builder --chown=root:root --chmod=755 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --link --from=builder --chown=1000:1000 --chmod=755 /go/src/github.com/cilium/cilium-cli/cilium /usr/local/bin/cilium

# cilium-cli-ci is based on ubuntu with cloud CLIs
FROM ubuntu:24.04@sha256:80dd3c3b9c6cecb9f1667e9290b3bc61b78c2678c02cbdae5f0fea92cc6734ab AS cilium-cli-ci
ENTRYPOINT []
LABEL maintainer="maintainer@cilium.io"
WORKDIR /root/app
COPY --link --from=builder /go/src/github.com/cilium/cilium-cli/cilium /usr/local/bin/cilium

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
