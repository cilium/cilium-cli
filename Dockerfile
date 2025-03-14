# Copyright Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

FROM --platform=${BUILDPLATFORM} docker.io/library/golang:1.24.1-alpine3.21@sha256:43c094ad24b6ac0546c62193baeb3e6e49ce14d3250845d166c77c25f64b0386 AS base
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
FROM ubuntu:24.04@sha256:72297848456d5d37d1262630108ab308d3e9ec7ed1c3286a32fe09856619a782 AS cilium-cli-ci
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
