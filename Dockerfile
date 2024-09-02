# syntax=docker/dockerfile:1.9@sha256:fe40cf4e92cd0c467be2cfc30657a680ae2398318afd50b0c80585784c604f28

# Copyright Authors of Cilium
# SPDX-License-Identifier: Apache-2.0
FROM --platform=${BUILDPLATFORM} golang:1.23.0-alpine3.20@sha256:d0b31558e6b3e4cc59f6011d79905835108c919143ebecc58f35965bf79948f4 AS base
RUN apk add --no-cache --update ca-certificates git make 
WORKDIR /go/src/github.com/cilium/cilium-cli
COPY go.* .
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .

# xx is a helper for cross-compilation
# when bumping to a new version analyze the new version for security issues
# then use crane to lookup the digest of that version so we are immutable
# crane digest tonistiigi/xx:1.5.0
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

FROM --platform=${BUILDPLATFORM} cgr.dev/chainguard/wolfi-base:latest@sha256:72c8bfed3266b2780243b144dc5151150015baf5a739edbbde53d154574f1607
LABEL maintainer="maintainer@cilium.io"
ENTRYPOINT [""]
CMD ["bash"]
ARG cilium_uid=1000
ARG cilium_gid=1000
ARG cilium_home=/home/cilium
RUN apk add --update --no-cache bash busybox kubectl && \
    addgroup -g ${cilium_gid} cilium && \
    adduser -D -h ${cilium_home} -u ${cilium_uid} -G cilium cilium
WORKDIR ${cilium_home}
COPY --link --from=builder --chown=${cilium_uid}:${cilium_gid} --chmod=755 /go/src/github.com/cilium/cilium-cli/cilium /usr/local/bin/cilium
COPY --link --from=builder --chown=root:root --chmod=755 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
USER ${cilium_uid}:${cilium_gid}
