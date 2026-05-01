# Copyright Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

ARG GOLANG_IMAGE=docker.io/library/golang:1.26.2@sha256:5f3787b7f902c07c7ec4f3aa91a301a3eda8133aa32661a3b3a3a86ab3a68a36
# BUILDPLATFORM is an automatic platform ARG enabled by Docker BuildKit.
# Represents the plataform where the build is happening, do not mix with
# TARGETARCH
FROM --platform=${BUILDPLATFORM} ${GOLANG_IMAGE} AS builder
WORKDIR /go/src/github.com/cilium/cilium-cli

# TARGETOS is an automatic platform ARG enabled by Docker BuildKit.
ARG TARGETOS
# TARGETARCH is an automatic platform ARG enabled by Docker BuildKit.
ARG TARGETARCH
RUN --mount=type=bind,readwrite,target=/go/src/github.com/cilium/cilium-cli \
    --mount=type=cache,target=/root/.cache \
    --mount=type=cache,target=/go/pkg \
    make GOARCH=${TARGETARCH} DESTDIR=/out/${TARGETOS}/${TARGETARCH} install

FROM gcr.io/distroless/static:latest@sha256:47b2d72ff90843eb8a768b5c2f89b40741843b639d065b9b937b07cd59b479c6 AS release
LABEL maintainer="maintainer@cilium.io"
ENTRYPOINT []
WORKDIR /root/app
# TARGETOS is an automatic platform ARG enabled by Docker BuildKit.
ARG TARGETOS
# TARGETARCH is an automatic platform ARG enabled by Docker BuildKit.
ARG TARGETARCH
COPY --from=builder /out/${TARGETOS}/${TARGETARCH}/usr/local/bin/cilium /usr/local/bin/cilium
