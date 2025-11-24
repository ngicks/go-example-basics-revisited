# syntax=docker/dockerfile:1

ARG TAG_GOVER="1.25.0"
ARG TAG_DISTRO="bookworm"

ARG HTTP_PROXY
ARG HTTPS_PROXY

ARG SSL_CERT_FILE="/etc/ssl/certs/ca-certificates.crt"
ARG NODE_EXTRA_CA_CERTS="/etc/ssl/certs/ca-certificates.crt"
# ARG SSH_AUTH_SOCK="/run/buildkit/S.gpg-agent.ssh"

ARG GOPATH="/go"
ARG GOCACHE="/root/.cache/go-build"
ARG GOPRIVATE=""
ARG CGO_ENABLED="0"
ARG MAIN_PKG_PATH="."

FROM --platform=${BUILDPLATFORM} docker.io/library/golang:${TAG_GOVER}-${TAG_DISTRO} AS builder

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    --mount=type=secret,id=certs,target=/etc/ssl/certs/ca-certificates.crt \
<<EOF
    rm -f /etc/apt/apt.conf.d/docker-clean
    echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' > /etc/apt/apt.conf.d/keep-cache
    apt-get update
    apt-get install -yqq --no-install-recommends git-lfs openssh-client
EOF

RUN <<EOF
    git config --global url."ssh://git@github.com".insteadOf https://github.com
    mkdir -p -m 0700 ~/.ssh 
    ssh-keyscan github.com >> ~/.ssh/known_hosts
EOF

WORKDIR /app/src
COPY go.mod go.sum ./
RUN --mount=type=secret,id=certs,target=/etc/ssl/certs/ca-certificates.crt \
    --mount=type=ssh,dst=/run/buildkit/S.gpg-agent.ssh \
    --mount=type=secret,id=goenv,target=/root/.config/go/env \
    --mount=type=cache,target=/go \
    --mount=type=cache,target=/root/.cache/go-build \
<<EOF
    go mod download
EOF

# https://docs.docker.com/build/building/context/#dockerignore-files
# COPY . .

RUN --mount=type=secret,id=certs,target=/etc/ssl/certs/ca-certificates.crt\
    --mount=type=ssh \
    --mount=type=secret,id=goenv,target=/root/.config/go/env \
    --mount=type=cache,target=/go \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=bind=target=. \
<<EOF
    # go generate ./...
    go build -o ../bin ${MAIN_PKG_PATH}
EOF

WORKDIR /app

FROM gcr.io/distroless/static-debian12@sha256:6ceafbc2a9c566d66448fb1d5381dede2b29200d1916e03f5238a1c437e7d9ea

COPY --from=builder /app/bin /app/bin

ENTRYPOINT [ "/app/bin" ]
